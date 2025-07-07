package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem file system storage
type FileSystem struct {
	Base string
}

// New initialize FileSystem storage
func NewFileSystem(base string) *FileSystem {
	absbase, err := filepath.Abs(base)
	if err != nil {
		fmt.Println("FileSystem storage's directory haven't been initialized")
	}
	return &FileSystem{Base: absbase}
}

// GetFullPath get full path from absolute/relative path
func (f *FileSystem) GetFullPath(path string) string {
	fullpath := path
	if !strings.HasPrefix(path, f.Base) {
		fullpath, _ = filepath.Abs(filepath.Join(f.Base, path))
	}
	return fullpath
}

// Get receive file with given path
func (f *FileSystem) Get(path string) (*os.File, error) {
	return os.Open(f.GetFullPath(path))
}

// GetStream get file as stream
func (f *FileSystem) GetStream(path string) (io.ReadCloser, error) {
	return os.Open(f.GetFullPath(path))
}

// Put store a reader into given path
func (f *FileSystem) Put(path string, reader io.Reader) (*Object, error) {
	var (
		fullpath = f.GetFullPath(path)
		err      = os.MkdirAll(filepath.Dir(fullpath), os.ModePerm)
	)

	if err != nil {
		return nil, err
	}

	dst, err := os.Create(fullpath)

	if err == nil {
		defer dst.Close()
		if seeker, ok := reader.(io.ReadSeeker); ok {
			seeker.Seek(0, 0)
		}
		_, err = io.Copy(dst, reader)
	}

	return &Object{Path: path, Name: filepath.Base(path), StorageInterface: f}, err
}

// Delete delete file
func (f *FileSystem) Delete(path string) error {
	return os.Remove(f.GetFullPath(path))
}

// List list all objects under current path
func (f *FileSystem) List(path string) ([]*Object, error) {
	var (
		objects  []*Object
		fullpath = f.GetFullPath(path)
	)

	filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if path == fullpath {
			return nil
		}

		if err == nil && !info.IsDir() {
			modTime := info.ModTime()
			objects = append(objects, &Object{
				Path:             strings.TrimPrefix(path, f.Base),
				Name:             info.Name(),
				LastModified:     &modTime,
				StorageInterface: f,
			})
		}
		return nil
	})

	return objects, nil
}

// GetEndpoint get endpoint, FileSystem's endpoint is /
func (f *FileSystem) GetEndpoint() string {
	return "/"
}

// GetURL get public accessible URL
func (f *FileSystem) GetURL(path string) (url string, err error) {
	return path, nil
}

// Symlink create symbolic link
func (f *FileSystem) Symlink(target string, name string) (err error) {
	var (
		targetPth    = f.GetFullPath(target)
		pth          = f.GetFullPath(name)
		targetRelDir string
	)

	if targetRelDir, err = filepath.Rel(filepath.Dir(targetPth), filepath.Dir(pth)); err != nil {
		return
	}

	targetPth = filepath.Join(targetRelDir, filepath.Base(targetPth))

	return os.Symlink(targetPth, pth)
}
