package fs

import (
	"errors"
	"os"
	"strings"

	"io/ioutil"
	"path/filepath"
)

func IsSymLink(m os.FileMode) bool {
	return (m & os.ModeSymlink) == os.ModeSymlink
}

func buildCheckSuffix(suffix string) func(string) bool {
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	if suffix == "" {
		return func(filename string) bool {
			return true
		}
	} else {
		return func(filename string) bool {
			return strings.HasSuffix(strings.ToUpper(filename), suffix)
		}
	}
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func DoListDir(dirPth string, suffix string, f func(fileName string) error) error {
	checkSuffix := buildCheckSuffix(suffix)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil { //忽略错误
		return nil
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		newFile := dirPth + PthSep + fi.Name()
		if !checkSuffix(newFile) {
			continue
		}
		if f(newFile) != nil {
			return errors.New("user quit")
		}
	}
	return nil
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func DoListDirEx(dirPth string, suffix string, f func(fullpath string, fileName string) error) error {
	checkSuffix := buildCheckSuffix(suffix)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil { //忽略错误
		return nil
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		newFile := dirPth + PthSep + fi.Name()
		if !checkSuffix(newFile) {
			continue
		}
		if f(newFile, fi.Name()) != nil {
			return errors.New("user quit")
		}
	}
	return nil
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string, ch chan<- string) error {
	checkSuffix := buildCheckSuffix(suffix)
	defer close(ch)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil { //忽略错误
		return nil
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		newFile := dirPth + PthSep + fi.Name()
		if !checkSuffix(newFile) {
			continue
		}
		ch <- newFile
	}
	return nil
}

func DoWalkDir(dirPth, suffix string, f func(fileName string, isdir bool) error) error {
	checkSuffix := buildCheckSuffix(suffix)

	err := filepath.Walk(dirPth,
		func(filename string, fi os.FileInfo, err error) error { //遍历目录
			if err != nil { //忽略错误
				// return err
				return nil
			}
			if fi.IsDir() { // 忽略目录
				f(filename, true)
				return nil
			}
			if !checkSuffix(filename) {
				return nil
			}
			f(filename, false)
			return nil
		})
	return err
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string, ch chan<- string) error {
	checkSuffix := buildCheckSuffix(suffix)
	defer close(ch)

	err := filepath.Walk(dirPth,
		func(filename string, fi os.FileInfo, err error) error { //遍历目录
			if err != nil { //忽略错误
				// return err
				return nil
			}
			if fi.IsDir() { // 忽略目录
				return nil
			}
			if !checkSuffix(filename) {
				return nil
			}

			if fi.Mode().IsRegular() {
				ch <- filename
			}
			return nil
		})

	return err
}

func PathExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func PathExists2(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}

	if e, ok := err.(*os.PathError); ok && e.Error() == os.ErrNotExist.Error() {
		return false, nil
	}
	return false, err
}
