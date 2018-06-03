package utils

import (
	"io/ioutil"
	"os"
)

func TempFileName(dir, prefix string) (string, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	return f.Name(), nil
}
