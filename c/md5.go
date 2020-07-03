package c

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

var ErrOpenFileError = errors.New("Open File Error")

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

func MD5FileByGoEx(filename string) (string, int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := md5.New()
	size, _ := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), size, nil
}
