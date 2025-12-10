package fs

import (
	"context"
	"io/fs"
	"os"
	"path"

	"github.com/hack-pad/hackpadfs"
	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = (*WebDavFS)(nil)

type WebDavFile struct {
	f hackpadfs.File
}

func (w *WebDavFile) Close() error {
	return w.f.Close()
}

func (w *WebDavFile) Read(p []byte) (n int, err error) {
	return w.f.Read(p)
}

func (w *WebDavFile) Seek(offset int64, whence int) (int64, error) {
	return hackpadfs.SeekFile(w.f, offset, whence)
}

func (w *WebDavFile) Stat() (fs.FileInfo, error) {
	return w.f.Stat()
}

func (w *WebDavFile) Write(p []byte) (n int, err error) {
	return hackpadfs.WriteFile(w.f, p)
}

func (w *WebDavFile) Readdir(count int) (fi []fs.FileInfo, err error) {
	var fe []hackpadfs.DirEntry
	if fe, err = hackpadfs.ReadDirFile(w.f, count); err != nil {
		return nil, err
	}
	fi = make([]fs.FileInfo, len(fe))
	for i, d := range fe {
		if fi[i], err = d.Info(); err != nil {
			return nil, err
		}
	}
	return
}

type WebDavFS struct {
	fs hackpadfs.FS
}

func NewWebDavFS(fs hackpadfs.FS) *WebDavFS {
	return &WebDavFS{fs: fs}
}

func (w *WebDavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return hackpadfs.MkdirAll(w.fs, path.Join(".", name), perm)
}

func (w *WebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (_ webdav.File, err error) {
	var f hackpadfs.File
	if f, err = hackpadfs.OpenFile(w.fs, path.Join(".", name), flag, perm); err != nil {
		return
	}
	return &WebDavFile{f}, nil
}

func (w *WebDavFS) RemoveAll(ctx context.Context, name string) error {
	return hackpadfs.RemoveAll(w.fs, path.Join(".", name))
}

func (w *WebDavFS) Rename(ctx context.Context, oldName, newName string) error {
	return hackpadfs.Rename(w.fs, path.Join(".", oldName), path.Join(".", newName))
}

func (w *WebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return hackpadfs.Stat(w.fs, path.Join(".", name))
}
