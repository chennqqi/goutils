package httputil

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func ReadJson(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	return d.Decode(v)
}

func DownloadFile(url string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(filename)
		return err
	}
	//TODO: uncompress
	buf := bufio.NewWriter(file)
	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	buf.Flush()
	return err
}
