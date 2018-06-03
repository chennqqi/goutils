package utils

import (
	"path/filepath"
)

func CleanFileName(parent, name string) string {
	nname := filepath.Clean(name)
	if nname == name {
		return name
	}
	if filepath.HasPrefix(nname, parent) {
		return nname
	}
	tname, err := TempFileName(parent, "clean_")
	if err == nil {
		return tname
	}
	//无论如何要返回一个结果，大不了就冲突吧
	return filepath.Join(parent, filepath.Base(name))
}
