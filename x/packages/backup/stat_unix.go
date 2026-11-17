//go:build unix

package backup

import (
	"io/fs"
	"syscall"
)

// devIno are the device and the inode of fi.
func devIno(fi fs.FileInfo) (uint64, uint64) {
	if st, ok := fi.Sys().(*syscall.Stat_t); ok {
		return uint64(st.Dev), uint64(st.Ino)
	}
	return 0, 0
}
