//go:build !unix

package backup

import "io/fs"

// devIno: none on this system (the size and the time tell alone).
func devIno(fs.FileInfo) (uint64, uint64) { return 0, 0 }
