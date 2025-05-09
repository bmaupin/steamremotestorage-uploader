package util

/*
#include <stdlib.h>
*/
import "C"
import (
    "../steam"
    "unsafe"
)

func GoStringToCString(s string) uintptr {
	return uintptr(unsafe.Pointer(&(([]byte(s))[0])))
}

func GoStringArrayToSteamStringArray(strings []string) (*steam.SteamParamStringArray_t, func()) {
	// Allocate an array of C strings (char**)
	cStringArray := C.malloc(C.size_t(len(strings)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	if cStringArray == nil {
		panic("failed to allocate memory for string array")
	}

	// Convert Go strings to C strings
	stringPtrs := (*[1 << 30]*C.char)(cStringArray)[:len(strings):len(strings)]
	for i, s := range strings {
		stringPtrs[i] = C.CString(s)
	}

	// Convert C string array to SteamParamStringArray_t
	steamStringArray := steam.NewSteamParamStringArray_t()
	steamStringArray.SetM_ppStrings((*string)(unsafe.Pointer(cStringArray)))
	steamStringArray.SetM_nNumStrings(len(strings))

	// Define cleanup function to free memory
	cleanup := func() {
		// Free each C string
		for _, ptr := range stringPtrs {
			C.free(unsafe.Pointer(ptr))
		}
		// Free the array itself
		C.free(cStringArray)
	}

	return &steamStringArray, cleanup
}
