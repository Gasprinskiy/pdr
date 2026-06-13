package a_unsafe

import "unsafe"

func StringToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
