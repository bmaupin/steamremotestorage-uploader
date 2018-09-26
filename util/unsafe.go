package util

import "unsafe"

func GoStringToCString(s string) uintptr {
	return uintptr(unsafe.Pointer(&(([]byte(s))[0])))
}
