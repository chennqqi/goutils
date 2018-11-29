package net

import (
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	_nsCache = NewNsCache(30 * time.Second)
)

func HttpTransportWithCache(maxIdle, cacheTo time.Duration) {
	http.DefaultClient.Transport = &http.Transport{
		MaxIdleConnsPerHost: 64,
		Dial: func(network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ip, _ := _nsCache.FetchOneString(address[:separator])
			return net.Dial("tcp", ip+address[separator:])
		},
	}
}
