package md5

//#cgo CFLAGS: -I .

/*
#include "md5.h"
#include <stdlib.h>
*/
import "C"

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"unsafe"
)

var ErrOpenFileError = errors.New("Open File Error")

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

func MD5FileByGo(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), nil
}
