package utils

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"

	gnet "github.com/shirou/gopsutil/net"
)

// 获取文件大小的接口
type Size interface {
	Size() int64
}

// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

func GetExternalIP() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	txt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(txt), nil
}

func GetInternalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipv4str := ipnet.IP.To4()
			if ipv4str != nil {
				return ipv4str.String(), nil
			}
		}
	}
	return "", nil
}

func GetRequestIP(r *http.Request) string {
	if r.Header.Get("X-Real-IP") != "" {
		return r.Header.Get("X-Real-IP")
	}
	if r.Header.Get("X-Forwarded-For") != "" {
		return r.Header.Get("X-Forwarded-For")
	}
	if r.Header.Get("Proxy-Client-IP") != "" {
		return r.Header.Get("Proxy-Client-IP")
	}
	if r.Header.Get("WL-Proxy-Client-IP") != "" {
		return r.Header.Get("WL-Proxy-Client-IP")
	}
	return r.RemoteAddr
}

func GetInternalIPByDevName(dev string) ([]string, error) {
	addrs, err := gnet.Interfaces()
	if err != nil {
		return []string{}, err
	}
	for _, a := range addrs {
		if a.Name == dev {
			var retIPs []string
			for _, addr := range a.Addrs {
				retIPs = append(retIPs, addr.String())
			}
			return retIPs, nil
		}
	}
	return []string{}, errors.New("not found dev or ip addr")
}

func GetUploadFileSize(upfile multipart.File) (int64, error) {
	if statInterface, ok := upfile.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		return fileInfo.Size(), nil
	}
	if sizeInterface, ok := upfile.(Size); ok {
		fsize := sizeInterface.Size()
		return fsize, nil
	}
	return 0, errors.New("not found stat and size interface")
}
