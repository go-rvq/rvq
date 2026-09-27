package fs

import (
	"os"

	"syscall"
)

func IsWriteMode(mode int) bool {
	return (mode&os.O_CREATE != 0 || mode&os.O_APPEND != 0 || mode&os.O_TRUNC != 0 || mode&os.O_WRONLY != 0 || mode&syscall.O_RDWR != 0)
}
