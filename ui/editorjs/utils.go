package editorjs

import (
	"reflect"
	"unsafe"
)

func StrToBytes(s *string) (bytes []byte) {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(s))
	// Our string is no longer referenced anywhere and could potentially be garbage collected
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	sliceHeader.Data = stringHeader.Data
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Cap = stringHeader.Len
	// runtime.KeepAlive(&t.Str)
	return
}
