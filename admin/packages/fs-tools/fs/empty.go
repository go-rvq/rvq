package fs

import (
	"errors"
	"io/fs"

	"github.com/hack-pad/hackpadfs"
)

var ErrEmptyFS = errors.New("Empty FS")

type emptyFS struct {
}

func (emptyFS) Open(string) (fs.File, error) {
	return nil, ErrEmptyFS
}

func EmptyFS() hackpadfs.FS {
	return NewROfs(&emptyFS{})
}
