package net

import (
	"context"
	"errors"
	"math/rand"
	stdnet "net"
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

func NewResolver(list []string) *stdnet.Resolver {
	var r stdnet.Resolver
	r.PreferGo = true
	r.Dial = dialer
	return &r
}

func UpdateResolveList(list []string) {
	serverList = list
}

func OverrideDefaultResolver(list []string) {
	var r stdnet.Resolver
	r.PreferGo = true
	r.Dial = dialer
	UpdateResolveList(list)
	stdnet.DefaultResolver = &r
}

func dialer(ctx context.Context, network, address string) (stdnet.Conn, error) {
	d := stdnet.Dialer{}
	list := serverList
	if len(list)==1{
		return d.DialContext(ctx, "udp", list[0])
	} else if len(list) == 0 {
		return nil, errors.New("EMPTY DNS LIST")
	}

	count := len(list)
	id := rand.Intn(count)
	conn, err := d.DialContext(ctx, "udp", list[id])
	if err == nil {
		return conn, err
	}
	//skip last
	for i := 0; i+1 < count; i++ {
		id++
		id = id % count
		conn, err = d.DialContext(ctx, "udp", list[id])
		if err == nil {
			return conn, err
		}
	}
	return conn, err
}
