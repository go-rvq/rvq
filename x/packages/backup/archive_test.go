package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// untar reads the tar.gz of data: its headers by name, the files' contents.
func untar(t *testing.T, data []byte) (map[string]*tar.Header, map[string]string) {
	t.Helper()
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	hds, files := map[string]*tar.Header{}, map[string]string{}
	for {
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		hds[hd.Name] = hd
		if hd.Typeflag == tar.TypeReg {
			b, _ := io.ReadAll(tr)
			files[hd.Name] = string(b)
		}
	}
	return hds, files
}

func TestWriteArchive(t *testing.T) {
	root := app(t)
	b, _ := New(Options{Root: root, Dir: filepath.Join(root, "data/backup"), Cache: &countingCache{m: map[CacheKey]Sum{}}})
	b.Source("s", func(ctx context.Context, st *State) error {
		if err := st.FS.AddTree("public"); err != nil {
			return err
		}
		return st.FS.AddTree("data/system")
	})
	if err := b.BuildState(context.Background(), &State{ID: "a"}); err != nil {
		t.Fatal(err)
	}
	m, _ := b.Manifest("a")

	// the whole state: every entry of the manifest, its files' sums, modes
	// and links kept
	var buf bytes.Buffer
	skipped, err := b.WriteArchive(&buf, "a", ArchiveOptions{})
	if err != nil || len(skipped) > 0 {
		t.Fatal(err, skipped)
	}
	hds, files := untar(t, buf.Bytes())
	for _, e := range m.Entries {
		name := e.Path
		if e.Type == TypeDir {
			name += "/"
		}
		hd, ok := hds[name]
		if !ok {
			t.Errorf("%s not in the archive", name)
			continue
		}
		if uint32(hd.Mode) != unixMode(e.Mode) || hd.ModTime.Unix() != e.MTime.Unix() {
			t.Errorf("%s: mode %o time %v", name, hd.Mode, hd.ModTime)
		}
		switch e.Type {
		case TypeFile:
			if Sum(sha256.Sum256([]byte(files[name]))) != e.Sum {
				t.Errorf("%s: sum", name)
			}
		case TypeLink:
			if hd.Typeflag != tar.TypeSymlink || hd.Linkname != e.Target {
				t.Errorf("%s: link %q", name, hd.Linkname)
			}
		}
	}
	if files["data/system/u/1.jpg"] != "jpeg" || hds["data/system/u/1.jpg"].Mode != 0o640 {
		t.Error("1.jpg")
	}

	// some paths, under a base: the names of the base
	buf.Reset()
	if _, err := b.WriteArchive(&buf, "a", ArchiveOptions{Base: "public", Paths: []string{"/public/index.html", "public/empty"}}); err != nil {
		t.Fatal(err)
	}
	hds, _ = untar(t, buf.Bytes())
	if len(hds) != 2 || hds["index.html"] == nil || hds["empty/"] == nil {
		t.Errorf("paths: %v", keys(hds))
	}

	// a file written in place after the state: left out, and said
	f, _ := os.OpenFile(filepath.Join(root, "public/index.html"), os.O_APPEND|os.O_WRONLY, 0)
	f.WriteString("!")
	f.Close()
	buf.Reset()
	skipped, err = b.WriteArchive(&buf, "a", ArchiveOptions{})
	if err != nil || len(skipped) != 1 || skipped[0] != "public/index.html" {
		t.Fatal(err, skipped)
	}
	hds, files = untar(t, buf.Bytes())
	if hds["public/index.html"] != nil || !bytes.Contains([]byte(files[SkippedFile]), []byte("public/index.html\tchanged")) {
		t.Errorf("skipped: %q", files[SkippedFile])
	}
}

func keys[V any](m map[string]V) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}
