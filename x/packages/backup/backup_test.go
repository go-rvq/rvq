package backup

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func write(t *testing.T, name, content string, mode fs.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(name, mode); err != nil {
		t.Fatal(err)
	}
}

func sumOf(s string) Sum { return Sum(sha256.Sum256([]byte(s))) }

// countingCache is a HashCache counting what it was asked.
type countingCache struct {
	m            map[CacheKey]Sum
	hits, misses atomic.Int32
}

func (c *countingCache) Lookup(k CacheKey) (Sum, bool) {
	s, ok := c.m[k]
	if ok {
		c.hits.Add(1)
	} else {
		c.misses.Add(1)
	}
	return s, ok
}
func (c *countingCache) Store(k CacheKey, s Sum) { c.m[k] = s }

// app is an application's directory: its site, its uploads, its binary.
func app(t *testing.T) string {
	root := t.TempDir()
	write(t, filepath.Join(root, "public/index.html"), "<h1>hi</h1>", 0o644)
	write(t, filepath.Join(root, "public/.transpiled/x.gad"), "x", 0o644)
	write(t, filepath.Join(root, "public/a.tmp"), "tmp", 0o644)
	write(t, filepath.Join(root, "data/system/u/1.jpg"), "jpeg", 0o640)
	write(t, filepath.Join(root, "app"), "binary", 0o755)
	if err := os.MkdirAll(filepath.Join(root, "public/empty"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("index.html", filepath.Join(root, "public/home.html")); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestBuildState(t *testing.T) {
	root := app(t)
	cache := &countingCache{m: map[CacheKey]Sum{}}
	b, err := New(Options{Root: root, Dir: filepath.Join(root, "data/backup"), Cache: cache})
	if err != nil {
		t.Fatal(err)
	}
	b.Source("db", func(ctx context.Context, st *State) error {
		return st.FS.AddGenerated("data/db.pg", 0o600, func(w io.Writer) error {
			_, err := io.WriteString(w, "dump")
			return err
		})
	})
	b.Source("site-files", func(ctx context.Context, st *State) error {
		return st.FS.AddTree("public", ".transpiled", "*.tmp")
	})
	b.Source("uploads", func(ctx context.Context, st *State) error {
		// the application keeps their sums: not read
		return st.FS.AddHashed(filepath.Join(root, "data/system/u/1.jpg"), sumOf("jpeg"), nil)
	})
	b.Source("app", func(ctx context.Context, st *State) error { return st.FS.Add("app") })

	var progress []Progress
	st := &State{ID: "s1", Progress: func(p Progress) { progress = append(progress, p) }}
	if err := b.BuildState(context.Background(), st); err != nil {
		t.Fatal(err)
	}
	r := st.Result
	if r.Files != 4 || r.Dirs != 2 || r.Links != 1 || r.Partial || len(r.Sources) != 4 || r.Bytes != int64(len("dump<h1>hi</h1>jpegbinary")) {
		t.Errorf("result %+v", r)
	}
	if len(progress) != 4 || progress[0].Source != "db" {
		t.Errorf("progress %+v", progress)
	}

	m, err := b.Manifest("s1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Header.ID != "s1" || strings.Join(m.Header.Sources, ",") != "db,site-files,uploads,app" {
		t.Errorf("header %+v", m.Header)
	}
	var paths []string
	for _, e := range m.Entries {
		paths = append(paths, e.Path)
	}
	if got := strings.Join(paths, " "); got != "app data/db.pg data/system/u/1.jpg public public/empty public/home.html public/index.html" {
		t.Errorf("paths %s", got)
	}
	check := func(p string, typ Type, mode fs.FileMode, sum Sum, target string) {
		t.Helper()
		e, ok := m.Find(p)
		if !ok || e.Type != typ || e.Mode != mode || e.Sum != sum || e.Target != target {
			t.Errorf("%s: %+v", p, e)
		}
	}
	check("app", TypeFile, 0o755, sumOf("binary"), "")
	check("data/db.pg", TypeFile, 0o600, sumOf("dump"), "")
	check("data/system/u/1.jpg", TypeFile, 0o640, sumOf("jpeg"), "")
	check("public/empty", TypeDir, 0o750, Sum{}, "")
	check("public/home.html", TypeLink, 0o777, Sum{}, "index.html")
	check("public/index.html", TypeFile, 0o644, sumOf("<h1>hi</h1>"), "")

	// the snapshot: a hard link of the file, the same inode
	e, _ := m.Find("public/index.html")
	f, err := b.OpenFile("s1", e)
	if err != nil {
		t.Fatal(err)
	}
	fi, _ := f.Stat()
	orig, _ := os.Stat(filepath.Join(root, "public/index.html"))
	f.Close()
	if !os.SameFile(fi, orig) {
		t.Error("not the same file")
	}
	// the generated one is the snapshot's alone
	e, _ = m.Find("data/db.pg")
	if f, err = b.OpenFile("s1", e); err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(f)
	f.Close()
	if string(got) != "dump" {
		t.Errorf("dump %q", got)
	}
	if _, err := os.Stat(filepath.Join(root, "data/db.pg")); !errors.Is(err, fs.ErrNotExist) {
		t.Error("the dump written in the application's directory")
	}

	// a file replaced (written beside and renamed): the state keeps it
	write(t, filepath.Join(root, "public/new.html"), "<h1>new</h1>", 0o644)
	if err := os.Rename(filepath.Join(root, "public/new.html"), filepath.Join(root, "public/index.html")); err != nil {
		t.Fatal(err)
	}
	e, _ = m.Find("public/index.html")
	f, err = b.OpenFile("s1", e)
	if err != nil {
		t.Fatal(err)
	}
	got, _ = io.ReadAll(f)
	f.Close()
	if string(got) != "<h1>hi</h1>" {
		t.Errorf("the snapshot changed: %q", got)
	}

	// a second state: what did not change is not hashed again
	misses := cache.misses.Load()
	if err := b.BuildState(context.Background(), &State{ID: "s2", Sources: []string{"site-files", "app"}}); err != nil {
		t.Fatal(err)
	}
	if n := cache.misses.Load() - misses; n != 1 { // index.html, replaced
		t.Errorf("hashed again: %d", n)
	}
	m2, err := b.Manifest("s2")
	if err != nil {
		t.Fatal(err)
	}
	if !m2.Header.Partial {
		t.Error("s2 not partial")
	}
	if e, _ := m2.Find("public/index.html"); e.Sum != sumOf("<h1>new</h1>") {
		t.Error("s2 has the old index")
	}
	// the states are never in a state
	if _, ok := m2.Find("data/backup"); ok {
		t.Error("the states in a state")
	}

	if err := b.BuildState(context.Background(), &State{ID: "s1"}); err == nil {
		t.Error("a state made twice")
	}
	if err := b.Remove("s1"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Manifest("s1"); err == nil {
		t.Error("s1 not removed")
	}
}

func TestBuildStateFails(t *testing.T) {
	root := app(t)
	dir := filepath.Join(t.TempDir(), "states")
	b, _ := New(Options{Root: root, Dir: dir, Cache: &countingCache{m: map[CacheKey]Sum{}}})
	b.Source("site-files", func(ctx context.Context, st *State) error { return st.FS.AddTree("public") })
	b.Source("bad", func(ctx context.Context, st *State) error { return errors.New("broken") })
	if err := b.BuildState(context.Background(), &State{ID: "x"}); err == nil || !strings.Contains(err.Error(), "bad: broken") {
		t.Fatalf("err %v", err)
	}
	if left, _ := os.ReadDir(dir); len(left) != 0 {
		t.Errorf("left %v", left)
	}
	for _, c := range []struct {
		name string
		st   *State
	}{
		{"id", &State{ID: "../x"}},
		{"source", &State{ID: "y", Sources: []string{"none"}}},
	} {
		if err := b.BuildState(context.Background(), c.st); err == nil {
			t.Errorf("%s accepted", c.name)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.BuildState(ctx, &State{ID: "z", Sources: []string{"site-files"}}); !errors.Is(err, context.Canceled) {
		t.Errorf("canceled: %v", err)
	}
}

func TestFSPaths(t *testing.T) {
	root := app(t)
	outside := t.TempDir()
	write(t, filepath.Join(outside, "f"), "out", 0o644)
	b, _ := New(Options{Root: root, Dir: filepath.Join(root, "data/backup"), Cache: &countingCache{m: map[CacheKey]Sum{}}})
	var errs []error
	b.Source("s", func(ctx context.Context, st *State) error {
		errs = append(errs,
			st.FS.Add(filepath.Join(outside, "f")),              // out of the root
			st.FS.Add("../x"),                                   // out of the root
			st.FS.Add(".backup/x"),                              // the client's
			st.FS.AddFrom("ext/f", filepath.Join(outside, "f")), // out of the root, at a path of the state
			st.FS.AddGenerated("app", 0o644, func(w io.Writer) error { _, err := io.WriteString(w, "other"); return err }),
		)
		return st.FS.Add("app")
	})
	err := b.BuildState(context.Background(), &State{ID: "p"})
	if errs[0] == nil || errs[1] == nil || errs[2] == nil || errs[3] != nil {
		t.Errorf("paths: %v", errs)
	}
	// app generated, then added from the disk with another content
	if !errors.Is(err, ErrDuplicate) {
		t.Errorf("duplicate: %v", err)
	}
}

func TestManifestRoundTrip(t *testing.T) {
	m := &Manifest{
		Header: Header{Version: ManifestVersion, ID: "s", Created: time.Unix(1700000000, 123).UTC(), Sources: []string{"a", "b"}, Partial: true},
		Entries: []Entry{
			{Path: "z", Type: TypeFile, Mode: fileMode(0o4755), MTime: time.Unix(1, 5), Size: 3, Sum: sumOf("abc")},
			{Path: "a\tb\nc", Type: TypeFile, Mode: 0o644, MTime: time.Unix(2, 0), Size: 0, Sum: sumOf("")},
			{Path: `"q`, Type: TypeDir, Mode: fileMode(0o1777), MTime: time.Unix(3, 0)},
			{Path: "l", Type: TypeLink, Mode: 0o777, MTime: time.Unix(4, 0), Target: "../a b"},
		},
	}
	var buf bytes.Buffer
	if err := WriteManifest(&buf, m); err != nil {
		t.Fatal(err)
	}
	got, err := ReadManifest(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.Header.ID != "s" || !got.Header.Created.Equal(m.Header.Created) || !got.Header.Partial || len(got.Header.Sources) != 2 {
		t.Errorf("header %+v", got.Header)
	}
	if len(got.Entries) != 4 {
		t.Fatalf("%d entries", len(got.Entries))
	}
	for i, e := range got.Entries {
		w := m.Entries[i] // sorted by WriteManifest
		if e.Path != w.Path || e.Type != w.Type || e.Mode != w.Mode || !e.MTime.Equal(w.MTime) || e.Size != w.Size || e.Sum != w.Sum || e.Target != w.Target {
			t.Errorf("%d: %+v, want %+v", i, e, w)
		}
	}
	if e, ok := got.Find("z"); !ok || e.Sum != sumOf("abc") {
		t.Error("find")
	}
}

func TestManifestRejects(t *testing.T) {
	for _, body := range []string{
		"nope",
		"#backup-manifest\t2\ts\t2020-01-01T00:00:00Z\t\t0",
		"#backup-manifest\t1\ts\t2020-01-01T00:00:00Z\t\t0\n-\t644\t0\t0\tf\ta",                     // a file with no sum
		"#backup-manifest\t1\ts\t2020-01-01T00:00:00Z\t\t0\n-\t755\t0\t0\td\t../x",                  // out
		"#backup-manifest\t1\ts\t2020-01-01T00:00:00Z\t\t0\n-\t755\t0\t0\td\tb\n-\t755\t0\t0\td\ta", // out of order
		"#backup-manifest\t1\ts\t2020-01-01T00:00:00Z\t\t0\n-\t755\t0\t0\td\t.backup",               // the client's
	} {
		var buf bytes.Buffer
		zw := gzipWriter(&buf)
		io.WriteString(zw, body+"\n")
		zw.Close()
		if _, err := ReadManifest(&buf); err == nil {
			t.Errorf("accepted %q", body)
		}
	}
}

func TestFileCache(t *testing.T) {
	name := filepath.Join(t.TempDir(), "c.gz")
	c := OpenFileCache(name)
	k1, k2 := CacheKey{Path: "/a\tb", Size: 1, MTime: 2, Dev: 3, Inode: 4}, CacheKey{Path: "/c", Size: 5}
	c.Store(k1, sumOf("1"))
	c.Store(k2, sumOf("2"))
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}
	c = OpenFileCache(name)
	if s, ok := c.Lookup(k1); !ok || s != sumOf("1") {
		t.Error("k1")
	}
	// k2 not used: dropped by the next save
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}
	c = OpenFileCache(name)
	if _, ok := c.Lookup(k2); ok {
		t.Error("k2 kept")
	}
	if fi, _ := os.Stat(name); fi.Mode().Perm() != 0o600 {
		t.Errorf("mode %v", fi.Mode())
	}
}

func TestOpenFile(t *testing.T) {
	root := app(t)
	b, _ := New(Options{Root: root, Dir: filepath.Join(root, "data/backup"), Cache: &countingCache{m: map[CacheKey]Sum{}}})
	b.Source("s", func(ctx context.Context, st *State) error { return st.FS.AddTree("public") })
	if err := b.BuildState(context.Background(), &State{ID: "o"}); err != nil {
		t.Fatal(err)
	}
	m, _ := b.Manifest("o")
	e, _ := m.Find("public/index.html")

	// written in place after the state (the same inode): not served
	f, err := os.OpenFile(filepath.Join(root, "public/index.html"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("!")
	f.Close()
	if _, err := b.OpenFile("o", e); !errors.Is(err, ErrChanged) {
		t.Errorf("changed: %v", err)
	}

	// a file of another file system (no hard link): served from where it is
	other := filepath.Join(t.TempDir(), "elsewhere.jpg")
	write(t, other, "jpeg", 0o640)
	fi, _ := os.Stat(other)
	snap := filepath.Join(b.Dir(), "o", filesDir, "data/x.jpg")
	os.Remove(snap)
	if err := writeOriginsAppend(filepath.Join(b.Dir(), "o", originsFile), "data/x.jpg", other); err != nil {
		t.Fatal(err)
	}
	x := Entry{Path: "data/x.jpg", Type: TypeFile, Size: 4, MTime: fi.ModTime(), Sum: sumOf("jpeg")}
	if f, err = b.OpenFile("o", x); err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(f)
	f.Close()
	if string(got) != "jpeg" {
		t.Errorf("origin %q", got)
	}
	// not in the state: nothing
	if _, err := b.OpenFile("o", Entry{Path: "public/none", Type: TypeFile}); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("none: %v", err)
	}
	if _, err := b.OpenFile("o", Entry{Path: "../../etc/passwd", Type: TypeFile}); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("out: %v", err)
	}
}

// The meta and the warnings of the sources; a tree whose sums its source
// knows (files named by their hash) is not read.
func TestSourceMetaWarnSums(t *testing.T) {
	root := t.TempDir()
	h := sumOf("jpeg").Hex()
	write(t, filepath.Join(root, "data/system/files", h[:2], h[2:4], h), "jpeg", 0o640)
	write(t, filepath.Join(root, "data/system/media/a.png"), "png", 0o640)
	cache := &countingCache{m: map[CacheKey]Sum{}}
	b, _ := New(Options{Root: root, Dir: filepath.Join(root, "data/backup"), Cache: cache})
	b.Source("uploads", func(ctx context.Context, st *State) error {
		st.FS.Warn("%d things", 2)
		if err := st.FS.Meta("dir", "data/system"); err != nil {
			return err
		}
		if err := st.FS.Meta("Bad Key", "x"); err == nil {
			t.Error("a bad key")
		}
		return st.FS.AddTreeWith("data/system", "", TreeOptions{Sum: func(rel string, fi fs.FileInfo) (Sum, bool) {
			if !strings.HasPrefix(rel, "files/") {
				return Sum{}, false
			}
			s, err := SumFromHex(filepath.Base(rel))
			return s, err == nil
		}})
	})
	st := &State{ID: "m"}
	if err := b.BuildState(context.Background(), st); err != nil {
		t.Fatal(err)
	}
	if cache.misses.Load() != 1 { // a.png alone
		t.Errorf("read %d", cache.misses.Load())
	}
	if w := st.Result.Sources[0].Warnings; len(w) != 1 || w[0] != "2 things" {
		t.Errorf("warnings %v", w)
	}
	m, _ := b.Manifest("m")
	if m.Header.Meta["uploads.dir"] != "data/system" {
		t.Errorf("meta %v", m.Header.Meta)
	}
	if e, ok := m.Find("data/system/files/" + h[:2] + "/" + h[2:4] + "/" + h); !ok || e.Sum != sumOf("jpeg") {
		t.Errorf("hashed %+v", e)
	}
}

func TestAppSource(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "bin/app"), "binary", 0o755)
	write(t, filepath.Join(root, "app.yml"), "env: {}", 0o600)
	write(t, filepath.Join(root, ".secure_key"), "secret", 0o600)
	outside := filepath.Join(t.TempDir(), "other-app")
	write(t, outside, "other", 0o755)
	for _, c := range []struct{ bin, at string }{{filepath.Join(root, "bin/app"), "bin/app"}, {outside, "app"}} {
		b, _ := New(Options{Root: root, Dir: filepath.Join(t.TempDir(), "states"), Cache: &countingCache{m: map[CacheKey]Sum{}}})
		b.Source("app", AppSource(AppOptions{Binary: c.bin, Config: []string{"app.yml", "app.yaml", ".secure_key"}}))
		if err := b.BuildState(context.Background(), &State{ID: "a"}); err != nil {
			t.Fatal(err)
		}
		m, _ := b.Manifest("a")
		if e, ok := m.Find(c.at); !ok || e.Mode != 0o755 {
			t.Errorf("%s: binary %+v", c.at, e)
		}
		if _, ok := m.Find(".secure_key"); ok {
			t.Error("the secure key in the state")
		}
		if m.Header.Meta["app.binary"] != c.at || m.Header.Meta["app.config"] != "app.yml" {
			t.Errorf("meta %v", m.Header.Meta)
		}
	}
}

func TestManifestMeta(t *testing.T) {
	m := &Manifest{Header: Header{ID: "s", Created: time.Unix(1, 0), Meta: map[string]string{"db.dump": "data/db.pg", "x.path": "a\tb"}}}
	var buf bytes.Buffer
	if err := WriteManifest(&buf, m); err != nil {
		t.Fatal(err)
	}
	got, err := ReadManifest(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.Header.Meta["db.dump"] != "data/db.pg" || got.Header.Meta["x.path"] != "a\tb" {
		t.Errorf("meta %v", got.Header.Meta)
	}
	m.Header.Meta = map[string]string{"Bad": "x"}
	if err := WriteManifest(&bytes.Buffer{}, m); err == nil {
		t.Error("a bad key written")
	}
	if s, err := SumFromHex(sumOf("a").Hex()); err != nil || s != sumOf("a") {
		t.Error("hex")
	}
}
