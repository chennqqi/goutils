package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Unzipsafe(t *testing.T) {
	testdir := "test"
	os.Mkdir(testdir, 0644)
	err := UnzipSafe("data/upload_media.zip", testdir, 0)
	if err != nil {
		t.Error(err)
		return
	}
	filepath.Walk(testdir, func(filename string, fi os.FileInfo, err error) error {
		t.Log(filename, fi)
		return nil
	})
}

func Test_ScanZip(t *testing.T) {
	testdir := "test"
	os.Mkdir(testdir, 0644)
	err := ScanZip("data/upload_media.zip", testdir, 0, func(filename string) error {
		f, e := os.Open(filename)
		if os.IsNotExist(e) {
			t.Error(filename, e)
			return e
		}
		f.Close()
		t.Log(filename)
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
}
