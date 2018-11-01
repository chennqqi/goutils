//utils of network
package net

import (
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
)

var (
	gateWayExp = regexp.MustCompile(`(?m)^(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s+(.*?)\s*$`)
)

func GetDefaultGateway() (net.IP, error) {
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
			item.Dst = NetHexToIPv4(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = NetHexToIPv4(gw)

			if item.Dst.String() == "0.0.0.0" {
				return item.Gateway, nil
			}
		}
	}

	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func GetGatewayByNic(nicName string) (net.IP, error) {
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
			item.Dst = NetHexToIPv4(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = NetHexToIPv4(gw)
			if item.Iface == nicName {
				return item.Gateway, nil
			}
		}
	}

	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func ListGateway() ([]*RouteItem, error) {
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
			item.Dst = NetHexToIPv4(dst)

			//gateway hex string
			voff++
			gw := string(txt[matched[i][voff*2]:matched[i][voff*2+1]])
			item.Gateway = NetHexToIPv4(gw)

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
