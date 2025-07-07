package storage

import (
	"io"
	"os"
	"time"
)

// Storage define common API to operate storage
type Storage interface {
	Get(path string) (*os.File, error)
	GetStream(path string) (io.ReadCloser, error)
	Put(path string, reader io.Reader) (*Object, error)
	Delete(path string) error
	List(path string) ([]*Object, error)
	GetURL(path string) (string, error)
	GetEndpoint() string
}

// Object content object
type Object struct {
	Path             string
	Name             string
	LastModified     *time.Time
	StorageInterface Storage
}

// Get retrieve object's content
func (object Object) Get() (*os.File, error) {
	return object.StorageInterface.Get(object.Path)
}

// StorageLinker is an Storage interface including methods to create symbolic link.
type Symlinker interface {
	Storage

	// Symlink create symbolic link from name to target.
	// If symbolic link creation isn't supported, return ErrSymlinkNotSupported.
	Symlink(target string, name string) (err error)
}
