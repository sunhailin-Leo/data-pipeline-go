package utils

import (
	"unsafe"
)

// StringToByte string to byte
func StringToByte(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ByteToString byte to string
func ByteToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
