package editorjs

import (
	"unsafe"
)

// StrToBytes is the bytes of *s, without a copy: they must not be changed.
func StrToBytes(s *string) []byte {
	return unsafe.Slice(unsafe.StringData(*s), len(*s))
}
