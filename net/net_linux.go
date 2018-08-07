//utils of network
package net

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
)

var (
	gateWayExp  = regexp.MustCompile(`(?m)^(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s*$`)
	ErrNotFound = errors.New("NOT FOUND")
)

func netHexToIPAddr(s string) net.IP {
	var v uint32
	fmt.Scanf(s, "%x", &v)
	ipstr := fmt.Sprintf("%d.%d.%d.%d",
		v&0xFF, (v>>8)&0xFF, (v>>16)&0xFF, (v>>24)&0xFF)
	return net.ParseIP(ipstr)
}

func netHexToIPMask(s string) net.IPMask {
	var v uint32
	fmt.Scanf(s, "%x", &v)
	return net.IPv4Mask(byte(v&0xFF), byte((v>>8)&0xFF),
		byte((v>>16)&0xFF), byte((v>>24)&0xFF))
}

func GetDefaultGateWay() (net.IP, error) {
	txt, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		return net.ParseIP("0.0.0.0"), err
	}
	matched := gateWayExp.FindAllSubmatchIndex(txt, -1)
	if len(matched) > 0 {
		//skip title line
		for i := 1; i < len(matched); i++ {
			item := &RouteItem{}
			var voff int //value offset
			voff++
			item.Iface = string(txt[matched[i][voff*2]:matched[i][voff*2+1]])

			//dst hex string
			voff++
			dst := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Dst = netHexToIPAddr(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = netHexToIPAddr(gw)

			if item.Dst.String() == "0.0.0.0" {
				return item.Gateway, nil
			}
		}
	}

	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func GetGateWayByNic(nicName string) (net.IP, error) {
	txt, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		return net.ParseIP("0.0.0.0"), err
	}
	matched := gateWayExp.FindAllSubmatchIndex(txt, -1)
	if len(matched) > 0 {
		//skip title line
		for i := 1; i < len(matched); i++ {
			item := &RouteItem{}
			var voff int //value offset
			voff++
			item.Iface = string(txt[matched[i][voff*2]:matched[i][voff*2+1]])

			//dst hex string
			voff++
			dst := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Dst = netHexToIPAddr(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = netHexToIPAddr(gw)
			if item.Iface == nicName {
				return item.Gateway, nil
			}
		}
	}

	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func ListGateWay() ([]*RouteItem, error) {
	txt, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		return nil, err
	}
	var ret []*RouteItem
	matched := gateWayExp.FindAllSubmatchIndex(txt, -1)

	//Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
	if len(matched) > 0 {
		//skip title line

		for i := 1; i < len(matched); i++ {
			item := new(RouteItem)
			var voff int //value offset
			voff++
			item.Iface = string(txt[matched[i][voff*2]:matched[i][voff*2+1]])

			//dst hex string
			voff++
			dst := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Dst = netHexToIPAddr(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = netHexToIPAddr(gw)

			voff++
			flags := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(flags, "%d", &item.Flags)

			voff++
			refCnt := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(refCnt, "%d", &item.RefCnt)

			voff++
			use := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(use, "%d", &item.Use)

			voff++
			metric := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(metric, "%d", &item.Metric)

			voff++
			mask := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Mask = netHexToIPMask(mask)

			voff++
			mtu := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(mtu, "%d", &item.MTU)

			voff++
			window := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(window, "%d", &item.Window)

			voff++
			irtt := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			fmt.Sscan(irtt, "%d", &item.IRTT)

			ret = append(ret, item)
		}
		return ret, nil
	}

	return ret, ErrNotFound
}
