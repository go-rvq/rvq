// Package backup makes states of an application — its database, its files,
// its binary, its configuration —: a snapshot of what a client syncs into a
// directory with the application's structure, where the application starts.
//
// The sources (Source) register what they back up through the state's FS
// (State.FS): files hashed then (Add, AddTree), files whose sha256 they know
// (AddHashed), content made then (AddGenerated: a database's dump). A state
// is BuildState's: in its folder (Options.Dir/<id>), the manifest (one line
// an entry: sha256, mode, mtime, size, type, path — WriteManifest) and the
// snapshot of its files, hard links of them where the file system allows
// it, so it stays as it was while the files change.
package backup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

// Options of a Builder.
type Options struct {
	// Root is the application's directory: the paths of the manifest are
	// of it.
	Root string
	// Dir is where the states are (a folder each).
	Dir string
	// Cache keeps the sums of the files between states; nil: a FileCache at
	// Dir/.hash-cache.gz, saved after each state.
	Cache HashCache
}

// SourceFunc registers what a source backs up through st.FS.
type SourceFunc func(ctx context.Context, st *State) error

type source struct {
	name string
	fn   SourceFunc
}

// Builder makes the states of the sources registered.
type Builder struct {
	root, dir string
	cache     HashCache
	mu        sync.RWMutex
	sources   []source
}

// New is a Builder of o.
func New(o Options) (*Builder, error) {
	if o.Root == "" || o.Dir == "" {
		return nil, errors.New("backup: Root and Dir are needed")
	}
	root, err := filepath.Abs(o.Root)
	if err != nil {
		return nil, err
	}
	dir, err := filepath.Abs(o.Dir)
	if err != nil {
		return nil, err
	}
	b := &Builder{root: root, dir: dir, cache: o.Cache}
	if b.cache == nil {
		b.cache = OpenFileCache(filepath.Join(dir, ".hash-cache.gz"))
	}
	return b, nil
}

// Root is the application's directory.
func (b *Builder) Root() string { return b.root }

// Dir is where the states are.
func (b *Builder) Dir() string { return b.dir }

var nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// Source registers fn as the source name ("db", "site-files"): the states
// have it in this order. A name registered twice, or not of [a-z0-9._-],
// panics.
func (b *Builder) Source(name string, fn SourceFunc) *Builder {
	if !nameRe.MatchString(name) {
		panic(fmt.Sprintf("backup: invalid source name %q", name))
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.sources {
		if s.name == name {
			panic(fmt.Sprintf("backup: source %q registered twice", name))
		}
	}
	b.sources = append(b.sources, source{name, fn})
	return b
}

// Sources are the names of the sources, in their order.
func (b *Builder) Sources() (out []string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.sources {
		out = append(out, s.name)
	}
	return
}

// State is a state being made: what is asked (ID, Sources, Progress) and,
// made, what it is (Result). The sources register their files through FS.
type State struct {
	// ID names it, its folder: [A-Za-z0-9][A-Za-z0-9._-]*.
	ID string
	// Created is when it is of (zero: now).
	Created time.Time
	// Sources are those it has; nil: all (a state without some is partial).
	Sources []string
	// Progress is told as the files are registered (nil: not told).
	Progress func(Progress)
	// FS is where the sources register their files (BuildState sets it).
	FS *FS
	// Result is what it is, once made.
	Result Result
}

// Progress is a state being made: the source registering and what it has.
type Progress struct {
	Source string
	Files  int
	Bytes  int64
}

// Result is a state made.
type Result struct {
	Files, Dirs, Links int
	Bytes              int64
	Partial            bool
	Sources            []SourceResult
	// ManifestSize and ManifestSum are of the manifest's file.
	ManifestSize int64
	ManifestSum  Sum
}

// SourceResult is what a source registered.
type SourceResult struct {
	Name     string
	Files    int
	Bytes    int64
	Duration time.Duration
}

var idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// The files of a state's folder.
const (
	manifestFile = "manifest.gz"
	originsFile  = "origins"
	filesDir     = "files"
)

// BuildState makes st: each source registers its files (in st.FS), then
// the manifest is written. It is made in a folder beside and renamed to
// Dir/<id> whole: a state there is complete. A source's error, or ctx's,
// fails it, and what it made is removed.
func (b *Builder) BuildState(ctx context.Context, st *State) (err error) {
	if !idRe.MatchString(st.ID) {
		return fmt.Errorf("backup: invalid state id %q", st.ID)
	}
	all := b.Sources()
	sel := st.Sources
	if sel == nil {
		sel = all
	}
	for _, s := range sel {
		if !slices.Contains(all, s) {
			return fmt.Errorf("backup: no source %q", s)
		}
	}
	if st.Created.IsZero() {
		st.Created = time.Now()
	}
	final := filepath.Join(b.dir, st.ID)
	if _, err := os.Lstat(final); err == nil {
		return fmt.Errorf("backup: state %q exists", st.ID)
	}
	if err := os.MkdirAll(b.dir, 0o700); err != nil {
		return err
	}
	work, err := os.MkdirTemp(b.dir, ".building-"+st.ID+"-")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(work)
		}
	}()
	fsys := &FS{b: b, st: st, ctx: ctx, snap: filepath.Join(work, filesDir), entries: map[string]Entry{}, origins: map[string]string{}}
	st.FS = fsys
	st.Result = Result{Partial: len(sel) < len(all)}
	for _, s := range b.sourcesOf(sel) {
		if err = ctx.Err(); err != nil {
			return err
		}
		fsys.source = SourceResult{Name: s.name}
		start := time.Now()
		if err = s.fn(ctx, st); err != nil {
			return fmt.Errorf("backup: source %s: %w", s.name, err)
		}
		fsys.source.Duration = time.Since(start)
		st.Result.Sources = append(st.Result.Sources, fsys.source)
	}
	if err = fsys.finish(work, sel); err != nil {
		return err
	}
	if err = os.Rename(work, final); err != nil {
		return err
	}
	if c, ok := b.cache.(interface{ Save() error }); ok {
		_ = c.Save() // a cache not saved costs time, not the state
	}
	return nil
}

func (b *Builder) sourcesOf(names []string) (out []source) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.sources {
		if slices.Contains(names, s.name) {
			out = append(out, s)
		}
	}
	return
}

// Manifest is the manifest of the state id.
func (b *Builder) Manifest(id string) (*Manifest, error) {
	if !idRe.MatchString(id) {
		return nil, fmt.Errorf("backup: invalid state id %q", id)
	}
	f, err := os.Open(filepath.Join(b.dir, id, manifestFile))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadManifest(f)
}

// ManifestPath is the file of the manifest of the state id.
func (b *Builder) ManifestPath(id string) string {
	return filepath.Join(b.dir, id, manifestFile)
}

// OpenFile opens the file e of the state id: its snapshot, else where it
// was (a file of another file system), which must still be as e says —
// size and time (ErrChanged otherwise; its sum, the ETag, is the client's
// check).
func (b *Builder) OpenFile(id string, e Entry) (*os.File, error) {
	if !idRe.MatchString(id) || e.Type != TypeFile || CheckPath(e.Path) != nil {
		return nil, fs.ErrNotExist
	}
	dir := filepath.Join(b.dir, id)
	name := filepath.Join(dir, filesDir, filepath.FromSlash(e.Path))
	if fi, err := os.Lstat(name); err != nil || !fi.Mode().IsRegular() {
		origins, err := readOrigins(filepath.Join(dir, originsFile))
		if err != nil {
			return nil, err
		}
		if name = origins[e.Path]; name == "" {
			return nil, fs.ErrNotExist
		}
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err == nil && (!fi.Mode().IsRegular() || fi.Size() != e.Size || !fi.ModTime().Equal(e.MTime)) {
		err = fmt.Errorf("%s: %w", e.Path, ErrChanged)
	}
	if err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

// Remove removes the state id: its manifest and its snapshot.
func (b *Builder) Remove(id string) error {
	if !idRe.MatchString(id) {
		return fmt.Errorf("backup: invalid state id %q", id)
	}
	return os.RemoveAll(filepath.Join(b.dir, id))
}

// FS is where the sources register the files of a state. Paths are of the
// application's directory: relative ("public/index.html") or absolute
// within it; a source of files out of it gives the path they have in the
// state (AddFrom, AddTreeFrom).
type FS struct {
	b       *Builder
	st      *State
	ctx     context.Context
	snap    string
	mu      sync.Mutex
	entries map[string]Entry
	origins map[string]string
	source  SourceResult
}

// Root is the application's directory.
func (f *FS) Root() string { return f.b.root }

// rel is p as a path of the manifest; p is relative to the root or
// absolute within it.
func (f *FS) rel(p string) (string, error) {
	if filepath.IsAbs(p) {
		r, err := filepath.Rel(f.b.root, p)
		if err != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("backup: %s is out of %s", p, f.b.root)
		}
		p = r
	}
	p = path.Clean(filepath.ToSlash(p))
	return p, CheckPath(p)
}

// Add registers the file, folder or link at p (of the root), hashed now
// (or by the cache).
func (f *FS) Add(p string) error {
	rel, err := f.rel(p)
	if err != nil {
		return err
	}
	return f.AddFrom(rel, filepath.Join(f.b.root, filepath.FromSlash(rel)))
}

// AddFrom registers src (a file, a folder, a link — not followed) at p of
// the state.
func (f *FS) AddFrom(p, src string) error {
	rel, err := f.rel(p)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return f.addInfo(rel, src, fi, nil)
}

// AddHashed registers the file at p (of the root) whose sha256 is known (a
// file the application keeps by its hash): it is not read. fi nil: it is
// stated.
func (f *FS) AddHashed(p string, sum Sum, fi fs.FileInfo) error {
	rel, err := f.rel(p)
	if err != nil {
		return err
	}
	src := filepath.Join(f.b.root, filepath.FromSlash(rel))
	if fi == nil {
		if fi, err = os.Lstat(src); err != nil {
			return err
		}
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("backup: %s is not a file", rel)
	}
	return f.addInfo(rel, src, fi, &sum)
}

// AddGenerated registers at p a file write makes now (a database's dump):
// kept in the state's snapshot, its sum taken as it is written.
func (f *FS) AddGenerated(p string, mode fs.FileMode, write func(w io.Writer) error) error {
	rel, err := f.rel(p)
	if err != nil {
		return err
	}
	dst := filepath.Join(f.snap, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".gen-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	h := newHasher()
	if err := write(io.MultiWriter(tmp, h)); err != nil {
		tmp.Close()
		return fmt.Errorf("backup: %s: %w", rel, err)
	}
	if err := errors.Join(tmp.Chmod(mode.Perm()), tmp.Close()); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		return err
	}
	fi, err := os.Stat(dst)
	if err != nil {
		return err
	}
	return f.put(Entry{Path: rel, Type: TypeFile, Mode: mode & (fs.ModePerm | fs.ModeSetuid | fs.ModeSetgid | fs.ModeSticky),
		MTime: fi.ModTime(), Size: fi.Size(), Sum: h.sum()}, "")
}

// AddTree registers p (of the root) and all under it: folders, files,
// links (not followed). exclude are patterns of path.Match of the paths
// under p ("cache/*"); one with no "/" is of a name too (".git", "*.tmp"):
// an excluded folder is left whole.
func (f *FS) AddTree(p string, exclude ...string) error {
	rel, err := f.relTree(p)
	if err != nil {
		return err
	}
	return f.AddTreeFrom(rel, filepath.Join(f.b.root, filepath.FromSlash(rel)), exclude...)
}

// relTree is rel, "." (the root itself) allowed.
func (f *FS) relTree(p string) (string, error) {
	if q := filepath.Clean(p); q == "." || q == f.b.root {
		return ".", nil
	}
	return f.rel(p)
}

// AddTreeFrom registers src and all under it at p of the state ("." the
// root).
func (f *FS) AddTreeFrom(p, src string, exclude ...string) error {
	rel := path.Clean(filepath.ToSlash(p))
	if rel != "." {
		if err := CheckPath(rel); err != nil {
			return err
		}
	}
	for _, pat := range exclude {
		if _, err := path.Match(pat, ""); err != nil {
			return fmt.Errorf("backup: exclude %q: %w", pat, err)
		}
	}
	return filepath.WalkDir(src, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err := f.ctx.Err(); err != nil {
			return err
		}
		sub, _ := filepath.Rel(src, name)
		sub = filepath.ToSlash(sub)
		if sub != "." && excluded(sub, exclude) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		at := path.Join(rel, sub)
		if at == "." || (rel == "." && (at == ReservedDir || strings.HasPrefix(at, ReservedDir+"/"))) {
			return nil // the root itself; the client's
		}
		if name == f.b.dir || strings.HasPrefix(name, f.b.dir+string(filepath.Separator)) {
			return filepath.SkipDir // the states: never in a state
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		return f.addInfo(at, name, fi, nil)
	})
}

func excluded(sub string, patterns []string) bool {
	base := path.Base(sub)
	for _, p := range patterns {
		if ok, _ := path.Match(p, sub); ok {
			return true
		}
		if !strings.Contains(p, "/") {
			if ok, _ := path.Match(p, base); ok {
				return true
			}
		}
	}
	return false
}

// addInfo registers src (fi, of Lstat) at rel; sum, a file's known.
func (f *FS) addInfo(rel, src string, fi fs.FileInfo, sum *Sum) error {
	e := Entry{Path: rel, Mode: fi.Mode() & (fs.ModePerm | fs.ModeSetuid | fs.ModeSetgid | fs.ModeSticky), MTime: fi.ModTime()}
	switch {
	case fi.Mode().IsDir():
		e.Type = TypeDir
		return f.put(e, "")
	case fi.Mode()&fs.ModeSymlink != 0:
		t, err := os.Readlink(src)
		if err != nil {
			return err
		}
		e.Type, e.Target = TypeLink, t
		return f.put(e, "")
	case !fi.Mode().IsRegular():
		return nil // a socket, a device, a pipe: not backed up
	}
	e.Type, e.Size = TypeFile, fi.Size()
	// the snapshot first: what is hashed is the inode it keeps
	read, origin := src, src
	dst := filepath.Join(f.snap, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err == nil && os.Link(src, dst) == nil {
		read, origin = dst, ""
		if lfi, err := os.Lstat(dst); err == nil {
			fi = lfi
			e.Size, e.MTime = fi.Size(), fi.ModTime()
		}
	}
	if sum != nil {
		e.Sum = *sum
	} else {
		k := keyOf(src, fi)
		s, ok := f.b.cache.Lookup(k)
		if !ok {
			var err error
			if s, err = hashFile(read, fi); err != nil {
				return err
			}
			f.b.cache.Store(k, s)
		}
		e.Sum = s
	}
	return f.put(e, origin)
}

// ErrDuplicate: a path registered twice, differently.
var ErrDuplicate = errors.New("registered twice, differently")

func (f *FS) put(e Entry, origin string) error {
	f.mu.Lock()
	if old, ok := f.entries[e.Path]; ok {
		f.mu.Unlock()
		if old.Type == e.Type && old.Sum == e.Sum && old.Target == e.Target {
			return nil
		}
		return fmt.Errorf("backup: %s: %w", e.Path, ErrDuplicate)
	}
	f.entries[e.Path] = e
	if origin != "" {
		f.origins[e.Path] = origin
	}
	r := &f.st.Result
	switch e.Type {
	case TypeFile:
		r.Files++
		r.Bytes += e.Size
		f.source.Files++
		f.source.Bytes += e.Size
	case TypeDir:
		r.Dirs++
	case TypeLink:
		r.Links++
	}
	p := Progress{Source: f.source.Name, Files: f.source.Files, Bytes: f.source.Bytes}
	f.mu.Unlock()
	if f.st.Progress != nil && e.Type == TypeFile {
		f.st.Progress(p)
	}
	return nil
}

// finish writes the manifest and the origins of the state into work.
func (f *FS) finish(work string, sources []string) error {
	m := &Manifest{Header: Header{Version: ManifestVersion, ID: f.st.ID, Created: f.st.Created, Sources: sources, Partial: f.st.Result.Partial}}
	for _, e := range f.entries {
		m.Entries = append(m.Entries, e)
	}
	name := filepath.Join(work, manifestFile)
	mf, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	h := newHasher()
	if err := errors.Join(WriteManifest(io.MultiWriter(mf, h), m), mf.Close()); err != nil {
		return err
	}
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	f.st.Result.ManifestSize, f.st.Result.ManifestSum = fi.Size(), h.sum()
	return writeOrigins(filepath.Join(work, originsFile), f.origins)
}
