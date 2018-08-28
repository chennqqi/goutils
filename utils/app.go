package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
ApplicationName: return `this program`'s basename, if run in windows,
returned string auto removed ".exe"*/
func ApplicationName() string {
	baseName := filepath.Base(os.Args[0])
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(baseName))
		if ext == ".exe" {
			return baseName[:len(baseName)-4]
		}
	}

	return baseName
}
