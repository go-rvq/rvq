package fs

import (
	"os"
	"path"
	"syscall"
)

// slashClean is equivalent to but slightly more efficient than
// path.Clean("/" + name).
func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}

func IsWriteMode(mode int) bool {
	return (mode&os.O_CREATE != 0 || mode&os.O_APPEND != 0 || mode&os.O_TRUNC != 0 || mode&os.O_WRONLY != 0 || mode&syscall.O_RDWR != 0)
}
