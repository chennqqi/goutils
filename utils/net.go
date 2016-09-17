package utils

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	gnet "github.com/shirou/gopsutil/net"
)

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
