package fs

import (
	"errors"

	"github.com/hack-pad/hackpadfs"
	"github.com/hack-pad/hackpadfs/mem"
)

// ErrEmptyFS: what an empty FS was asked has nothing.
var ErrEmptyFS = errors.New("empty FS")

// EmptyFS is a file system of no file: its root, an empty folder, only
// read.
func EmptyFS() hackpadfs.FS {
	m, err := mem.NewFS()
	if err != nil {
		panic(err)
	}
	return NewROfs(m)
}
