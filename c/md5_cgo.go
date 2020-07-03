// +build cgo

package c

// #cgo CFLAGS: -I.

/*
#include "md5.h"
#include "md5.hxx"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/hex"
	"unsafe"
)

func MD5FileByC(filename string) (string, error) {
	pName := C.CString(filename)
	defer C.free(unsafe.Pointer(pName))
	var md5bytes [16]byte
	r := C.MD5File(pName, (*C.uchar)(&md5bytes[0]))
	if int(r) != 0 {
		return "", ErrOpenFileError
	}
	return hex.EncodeToString(md5bytes[:]), nil
}

func MD5FileByCEx(filename string) (string, int64, error) {
	pName := C.CString(filename)
	defer C.free(unsafe.Pointer(pName))
	var md5bytes [16]byte
	r := C.MD5FileExt(pName, (*C.uchar)(&md5bytes[0]))
	if int(r) < 0 {
		return "", 0, ErrOpenFileError
	}
	return hex.EncodeToString(md5bytes[:]), int64(r), nil
}
