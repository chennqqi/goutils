package net

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	gnet "github.com/shirou/gopsutil/net"
	"github.com/tomasen/realip"
)

var (
	ErrNotFound = errors.New("NOT FOUND")
)

// 获取文件大小的接口
type Size interface {
	Size() int64
}

// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

// 返回公网出口IP
func GetExternalIP() (string, error) {
	resp, err := http.Get("http://ipaddr.site")
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

//将整数IP转换为Go
func Inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

//将IP转换为整数值
func Inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//返回是否是公网IP
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// 返回主机HOST IP
func GetHostIP() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return "", err
	}

	return addr.String(), nil
}

// 返回主机内部IP（顺序第一个）
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

//返回HTTP请求IP,优先级(X-Real-IP>X-Forwarded-For>Proxy-Client-IP>WL-Proxy-Client-IP")
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
	return realip.FromRequest(r)
}

//根据网卡设备名称返回内部IP
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

//返回HTTP upload文件的大小
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

//在主机拥有多个IP地址时，返回能够连通对端网络的自身主机IP
func GetLocalConnectIP(proto string, addr string) (string, error) {
	if proto == "" {
		proto = "tcp"
	}
	conn, err := net.Dial(proto, addr)
	if err != nil {
		return "", nil
	}
	defer conn.Close()
	laddr := conn.LocalAddr().String()
	var rip string
	v := strings.Split(laddr, ":")
	if len(v) > 1 {
		rip = v[0]
	}
	return rip, nil
}

type RouteItem struct {
	Iface   string
	Dst     net.IP
	Gateway net.IP

	Flags  uint32
	RefCnt int
	Use    int
	Metric int
	Mask   net.IPMask
	MTU    int
	Window int
	IRTT   int
}

func NetHexToIPv4(s string) net.IP {
	var v uint32
	fmt.Sscanf(s, "%x", &v)
	ipstr := fmt.Sprintf("%d.%d.%d.%d",
		v&0xFF, (v>>8)&0xFF, (v>>16)&0xFF, (v>>24)&0xFF)
	return net.ParseIP(ipstr)
}

func netHexToIPMask(s string) net.IPMask {
	var v uint32
	fmt.Sscanf(s, "%x", &v)
	return net.IPv4Mask(byte(v&0xFF), byte((v>>8)&0xFF),
		byte((v>>16)&0xFF), byte((v>>24)&0xFF))
}

func netStringToIPv4Mask(s string) net.IPMask {
	var a, b, c, d uint32
	fmt.Sscanf(s, "%d.%d.%d.%d", &a, &b, &c, &d)
	return net.IPv4Mask(byte(a&0xFF), byte(b&0xFF),
		byte(c&0xFF), byte(d&0xFF))
}
