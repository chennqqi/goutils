package net

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	serverList []string
)

type Resolver struct {
	serverList []string
	mutex      sync.Mutex
}

func NewResolver(list []string) *net.Resolver {
	var r net.Resolver
	r.PreferGo = true
	r.Dial = dialer
	return &r
}

func UpdateResolvList(list []string) {
	serverList = list
}

func OverrideDefaultResolver(list []string) {
	var r net.Resolver
	r.PreferGo = true
	r.Dial = dialer
	net.DefaultResolver = &r
}

func dialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	list := serverList

	count := len(list)
	id := rand.Intn(count)
	net, err := d.DialContext(ctx, "udp", list[id])
	if err == nil {
		return net, err
	}
	//skip last
	for i := 0; i+1 < count; i++ {
		id++
		id = id % count
		net, err = d.DialContext(ctx, "udp", list[id])
		if err == nil {
			return net, err
		}
	}
	return net, err
}
