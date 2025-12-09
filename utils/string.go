package utils

import "unsafe"

func UnsafeStringToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))[:]
}
