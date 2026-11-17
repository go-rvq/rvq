package backup

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/go-rvq/rvq/x/packages/backup"
	"golang.org/x/net/webdav"
)

// stateFS is the files of a state as a file system, for the explorer of
// its page (Files): its tree is the manifest's — the folders a path has but
// the manifest does not, made —, each file opened by the core (its
// snapshot, or where it was). Read only: every change is os.ErrPermission.
type stateFS struct {
	b    *Builder
	id   string
	ents map[string]*stateEntry
}

type stateEntry struct {
	backup.Entry
	kids []string
}

// stateFSOf is the file system of the state id, made of its manifest once
// (kept with it, until the state is removed).
func (b *Builder) stateFSOf(id string) (*stateFS, error) {
	if v, ok := b.stateFSs.Load(id); ok {
		return v.(*stateFS), nil
	}
	m, err := b.manifest(id)
	if err != nil {
		return nil, err
	}
	s := &stateFS{b: b, id: id, ents: map[string]*stateEntry{"": {Entry: backup.Entry{Type: backup.TypeDir, Mode: 0o755, MTime: m.Header.Created}}}}
	var add func(p string) *stateEntry
	add = func(p string) *stateEntry {
		if e, ok := s.ents[p]; ok {
			return e
		}
		e := &stateEntry{Entry: backup.Entry{Path: p, Type: backup.TypeDir, Mode: 0o755, MTime: m.Header.Created}}
		s.ents[p] = e
		parent := path.Dir(p)
		if parent == "." {
			parent = ""
		}
		pe := add(parent)
		pe.kids = append(pe.kids, path.Base(p))
		return e
	}
	for _, e := range m.Entries {
		add(e.Path).Entry = e
	}
	for _, e := range s.ents {
		sort.Strings(e.kids)
	}
	b.stateFSs.Store(id, s)
	return s, nil
}

func (s *stateFS) entry(name string) (*stateEntry, error) {
	e, ok := s.ents[strings.Trim(path.Clean("/"+name), "/")]
	if !ok {
		return nil, os.ErrNotExist
	}
	return e, nil
}

func (s *stateFS) Mkdir(context.Context, string, os.FileMode) error { return os.ErrPermission }
func (s *stateFS) RemoveAll(context.Context, string) error          { return os.ErrPermission }
func (s *stateFS) Rename(context.Context, string, string) error     { return os.ErrPermission }

func (s *stateFS) Stat(_ context.Context, name string) (os.FileInfo, error) {
	e, err := s.entry(name)
	if err != nil {
		return nil, err
	}
	return stateInfo{e}, nil
}

func (s *stateFS) OpenFile(_ context.Context, name string, flag int, _ os.FileMode) (webdav.File, error) {
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND) != 0 {
		return nil, os.ErrPermission
	}
	e, err := s.entry(name)
	if err != nil {
		return nil, err
	}
	switch e.Type {
	case backup.TypeDir:
		return &stateDir{s: s, e: e}, nil
	case backup.TypeFile:
		f, err := s.b.core.OpenFile(s.id, e.Entry)
		if err != nil {
			return nil, err
		}
		return &stateFile{File: f, e: e}, nil
	}
	// a link: listed, not opened
	return nil, os.ErrNotExist
}

// stateInfo is an entry as a fs.FileInfo.
type stateInfo struct{ e *stateEntry }

func (i stateInfo) Name() string {
	if i.e.Path == "" {
		return "/"
	}
	return path.Base(i.e.Path)
}
func (i stateInfo) Size() int64        { return i.e.Size }
func (i stateInfo) ModTime() time.Time { return i.e.MTime }
func (i stateInfo) IsDir() bool        { return i.e.Type == backup.TypeDir }
func (i stateInfo) Sys() any           { return nil }
func (i stateInfo) Mode() fs.FileMode {
	switch i.e.Type {
	case backup.TypeDir:
		return i.e.Mode | fs.ModeDir
	case backup.TypeLink:
		return i.e.Mode | fs.ModeSymlink
	}
	return i.e.Mode
}

// stateDir is a folder opened: its entries.
type stateDir struct {
	s   *stateFS
	e   *stateEntry
	pos int
}

func (d *stateDir) Close() error                   { return nil }
func (d *stateDir) Read([]byte) (int, error)       { return 0, io.EOF }
func (d *stateDir) Seek(int64, int) (int64, error) { return 0, os.ErrInvalid }
func (d *stateDir) Write([]byte) (int, error)      { return 0, os.ErrPermission }
func (d *stateDir) Stat() (os.FileInfo, error)     { return stateInfo{d.e}, nil }
func (d *stateDir) Readdir(count int) ([]fs.FileInfo, error) {
	var out []fs.FileInfo
	for d.pos < len(d.e.kids) && (count <= 0 || len(out) < count) {
		out = append(out, stateInfo{d.s.ents[strings.TrimPrefix(d.e.Path+"/"+d.e.kids[d.pos], "/")]})
		d.pos++
	}
	if count > 0 && len(out) == 0 {
		return nil, io.EOF
	}
	return out, nil
}

// stateFile is a file opened: its content, read only.
type stateFile struct {
	*os.File
	e *stateEntry
}

func (f *stateFile) Write([]byte) (int, error)          { return 0, os.ErrPermission }
func (f *stateFile) Stat() (os.FileInfo, error)         { return stateInfo{f.e}, nil }
func (f *stateFile) Readdir(int) ([]fs.FileInfo, error) { return nil, os.ErrInvalid }
