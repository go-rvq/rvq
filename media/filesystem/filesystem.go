package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/qor5/admin/v3/media/base"
)

var (
	_ base.Media          = (*FileSystem)(nil)
	_ base.MediaSymlinker = (*FileSystem)(nil)
)

// FileSystem defined a media library storage using file system
type FileSystem struct {
	base.Base
}

// GetFullPath return full file path from a relative file path
func (f *FileSystem) GetFullPath(url string, option *base.Option) (path string, err error) {
	if option != nil && option.Get("path") != "" {
		path = filepath.Join(option.Get("path"), url)
	} else {
		path = filepath.Join("./public", url)
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	return
}

// Store save reader's context with name
func (f *FileSystem) Store(name string, option *base.Option, reader io.Reader) (err error) {
	if fullpath, err := f.GetFullPath(name, option); err == nil {
		if dst, err := os.Create(fullpath); err == nil {
			_, err = io.Copy(dst, reader)
		}
	}
	return err
}

// Symlink create symbolic link
func (f *FileSystem) Symlink(srcName string, name string, option *base.Option) (err error) {
	var fullpathSrc, fullpath string
	if fullpathSrc, err = f.GetFullPath(srcName, option); err != nil {
		return
	}
	if fullpath, err = f.GetFullPath(name, option); err != nil {
		return
	}

	return os.Link(fullpathSrc, fullpath)
}

// Retrieve retrieve file content with url
func (f *FileSystem) Retrieve(url string) (base.FileInterface, error) {
	if fullpath, err := f.GetFullPath(url, nil); err == nil {
		return os.Open(fullpath)
	}
	return nil, os.ErrNotExist
}
