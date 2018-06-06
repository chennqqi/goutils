package utils

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DEFAULT_UNZIP_LIMIT = 64 * 1024 * 1024 * 1024
)

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ZipSize(zipR *zip.ReadCloser) uint64 {
	var totalSize uint64
	for idx := 0; idx < len(zipR.File); idx++ {
		totalSize += zipR.File[idx].UncompressedSize64
	}
	return totalSize
}

func UnzipSafe(archive, target string, sizeLimit uint64) error {
	if sizeLimit == 0 {
		sizeLimit = DEFAULT_UNZIP_LIMIT
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	if ZipSize(reader) > sizeLimit {
		return errors.New("SIZE OVER LIMIT")
	}

	for _, file := range reader.File {
		filePath := filepath.Join(target, file.Name)
		filePath = CleanFileName(target, filePath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}
		os.MkdirAll(filepath.Dir(filePath), 0750)

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func ScanZip(archive, tmpDir string, sizeLimit uint64, scanCall func(filename string) error) error {
	if sizeLimit == 0 {
		sizeLimit = DEFAULT_UNZIP_LIMIT
	}
	nTmp, err := ioutil.TempDir(tmpDir, "szipproc_")
	if err != nil {
		return err
	}
	defer os.RemoveAll(nTmp)

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	if ZipSize(reader) > sizeLimit {
		return errors.New("ZIP UNCOMPRESS OVERLIMIT")
	}

	for _, file := range reader.File {
		filePath := filepath.Join(nTmp, file.Name)
		filePath = CleanFileName(nTmp, filePath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}
		os.MkdirAll(filepath.Dir(filePath), 0750)

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer os.Remove(filePath)

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			targetFile.Close()
			return err
		}
		targetFile.Close()

		err = scanCall(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
