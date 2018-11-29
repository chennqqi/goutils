package net

// Package dnscache caches DNS lookups
// origin code https://github.com/viki-org/dnscache/blob/master/dnscache.go
// MIT license

import (
	"net"
	"time"
)

const (
	CACHE_IP_MAX = 64
)

type NSCache struct {
	cache map[string][]net.IP
}

func NewNsCache(refreshRate time.Duration) *NSCache {
	nc := &NSCache{
		cache: make(map[string][]net.IP, CACHE_IP_MAX),
	}
	if refreshRate > 0 {
		go nc.autoRefresh(refreshRate)
	}
	return nc
}

func (r *NSCache) Fetch(address string) ([]net.IP, error) {
	cache := r.cache
	ips, exists := cache[address]
	if exists {
		return ips, nil
	}

	return r.Lookup(address)
}

func (r *NSCache) FetchOne(address string) (net.IP, error) {
	ips, err := r.Fetch(address)
	if err != nil || len(ips) == 0 {
		return nil, err
	}
	return ips[0], nil
}

func (r *NSCache) FetchOneString(address string) (string, error) {
	ip, err := r.FetchOne(address)
	if err != nil || ip == nil {
		return "", err
	}
	return ip.String(), nil
}

func (r *NSCache) Refresh() {
	i := 0

	cache := r.cache
	addresses := make([]string, len(cache))
	for key, _ := range r.cache {
		addresses[i] = key
		i++
	}

	for _, address := range addresses {
		r.Lookup(address)
		time.Sleep(time.Second * 2)
	}
}

func (r *NSCache) Lookup(address string) ([]net.IP, error) {
	ips, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}
	oldCache := r.cache
	newCache := make(map[string][]net.IP, CACHE_IP_MAX)

	delete(oldCache, address)
	for k, v := range oldCache {
		newCache[k] = v
	}
	newCache[address] = ips
	r.cache = newCache
	return ips, nil
}

func (r *NSCache) autoRefresh(rate time.Duration) {
	for {
		time.Sleep(rate)
		r.Refresh()
	}
}
