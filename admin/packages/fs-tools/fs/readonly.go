package fs

import (
	"errors"
	"time"

	"github.com/hack-pad/hackpadfs"
)

var ErrReadOnlyFS = errors.New("Read-only filesystem")

type ROfs struct {
	fs hackpadfs.FS
}

// Open implements hackpadfs.FS
func (fs *ROfs) Open(name string) (hackpadfs.File, error) {
	return fs.fs.Open(name)
}

// OpenFile implements hackpadfs.OpenFileFS
func (fs *ROfs) OpenFile(name string, flag int, perm hackpadfs.FileMode) (hackpadfs.File, error) {
	if IsWriteMode(int(perm)) {
		return nil, ErrReadOnlyFS
	}
	return hackpadfs.OpenFile(fs.fs, name, flag, perm)
}

// Mkdir implements hackpadfs.MkdirFS
func (fs *ROfs) Mkdir(name string, perm hackpadfs.FileMode) error {
	return ErrReadOnlyFS
}

// MkdirAll implements hackpadfs.MkdirAllFS
func (fs *ROfs) MkdirAll(path string, perm hackpadfs.FileMode) error {
	return ErrReadOnlyFS
}

// Remove implements hackpadfs.RemoveFS
func (fs *ROfs) Remove(name string) error {
	return ErrReadOnlyFS
}

// Rename implements hackpadfs.RenameFS
func (fs *ROfs) Rename(oldname, newname string) error {
	return ErrReadOnlyFS
}

// Stat implements hackpadfs.StatFS
func (fs *ROfs) Stat(name string) (hackpadfs.FileInfo, error) {
	return hackpadfs.Stat(fs.fs, name)
}

// Chmod implements hackpadfs.ChmodFS
func (fs *ROfs) Chmod(name string, mode hackpadfs.FileMode) error {
	return ErrReadOnlyFS
}

// Chtimes implements hackpadfs.ChtimesFS
func (fs *ROfs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return ErrReadOnlyFS
}

func NewROfs(fs hackpadfs.FS) *ROfs {
	return &ROfs{fs}
}
