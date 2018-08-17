package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestbuildCheckSuffix(t *testing.T) {
	checkSuffix := buildCheckSuffix(".go")
	var ok bool
	ok = checkSuffix("a.go")
	assert.True(t, ok)

	ok = checkSuffix("a.goxx")
	assert.False(t, ok)

	ok = checkSuffix(".go")
	assert.True(t, ok)

	ok = checkSuffix("go")
	assert.False(t, ok)

	checkSuffix = buildCheckSuffix("*.go")
	ok = checkSuffix("a.go")
	assert.False(t, ok)

	checkSuffix = buildCheckSuffix("")
	ok = checkSuffix("a.go")
	assert.True(t, ok)
	ok = checkSuffix("a.python")
	assert.True(t, ok)
	ok = checkSuffix("python")
	assert.True(t, ok)
	ok = checkSuffix("")
	assert.True(t, ok)
}

func TestDoListDir(t *testing.T) {
	var count int
	err := DoListDir(".", "", func(filename string) error {
		count++
		assert.FileExists(t, filename)
		assert.NotContains(t, filename, "sub.txt")
		return nil
	})
	assert.Nil(t, err)
	assert.True(t, count > 0)
}

func TestDoListDirEx(t *testing.T) {
	var count int
	err := DoListDirEx(".", "", func(fullname, filename string) error {
		count++
		assert.FileExists(t, fullname)
		assert.NotContains(t, fullname, "sub.txt")
		assert.NotContains(t, filename, `/`)
		assert.NotContains(t, filename, `\`)
		return nil
	})
	assert.Nil(t, err)
	assert.True(t, count > 0)
}

func TestListDir(t *testing.T) {
	ch := make(chan string)
	stopChan := make(chan struct{})
	var err error
	go func() {
		err = ListDir(".", "", ch)
	}()

	var count int
	for filename := range ch {
		count++
		assert.FileExists(t, filename)
		assert.NotContains(t, filename, "sub.txt")
	}
	close(stopChan)

	assert.Nil(t, err)
	assert.True(t, count > 0)
}

func TestDoWalkDir(t *testing.T) {
	//进入子目录
	var count int
	var hasSubdir bool
	var hasDir bool
	err := DoWalkDir(".", "", func(filename string, isdir bool) error {
		count++
		if strings.Contains(filename, "sub.txt") {
			hasSubdir = true
		}
		if isdir {
			hasDir = true
			assert.DirExists(t, filename)
		} else {
			assert.FileExists(t, filename)
		}
		return nil
	})

	assert.Nil(t, err)
	assert.True(t, count > 0)
	assert.True(t, hasSubdir)
	assert.True(t, hasDir)
}

func TestWalkDir(t *testing.T) {
	ch := make(chan string)
	var err error

	go func() {
		err = WalkDir(".", "", ch)
	}()

	var count int
	var hasSubdir bool

	for filename := range ch {
		count++
		if strings.Contains(filename, "sub.txt") {
			hasSubdir = true
		}

		assert.FileExists(t, filename)
	}
	assert.Nil(t, err)
	assert.True(t, count > 0)
	assert.True(t, hasSubdir)
}

func TestPathExists(t *testing.T) {
	exist := PathExists("file_test.go")
	assert.True(t, exist)
}

func TestPathExists2(t *testing.T) {
	exist, err := PathExists2("file_test.go")
	assert.Nil(t, err)
	assert.True(t, exist)
}
