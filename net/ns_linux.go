package net

import (
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"sort"
)

var (
	ifcfgNsExp  = regexp.MustCompile(`(?m)^\s*DNS(\d)\s*=\s*(.*)$`)
	resolvNsExp = regexp.MustCompile(`(?m)^\s*nameserver\s*(.*)$`)
)

type nsRecord struct {
	Index int
	Value string
}

type nsRecords []nsRecord

func (t nsRecords) Less(i, j int) bool { return t[i].Index < t[j].Index }
func (t nsRecords) Len() int           { return len(t) }
func (t nsRecords) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func GetRhNsByNic(name string) (NSServers, error) {
	var ns NSServers
	fname := fmt.Sprintf("/etc/sysconfig/network-scripts/ifcfg-%v", name)
	txt, err := ioutil.ReadFile(fname)
	if err != nil {
		return ns, err
	}
	var recordsList nsRecords

	match := ifcfgNsExp.FindAllSubmatchIndex(txt, 2)
	for i := 0; i < len(match); i++ {
		if len(match[i]) >= 6 {
			var n nsRecord
			id := string(txt[match[i][2]:match[i][3]])
			n.Value = string(txt[match[i][4]:match[i][5]])
			fmt.Sscanf(id, "%d", &n.Index)
			recordsList = append(recordsList, n)
		}
	}
	sort.Sort(recordsList)
	ns = make(NSServer, len(recordsList))
	for i := 0; i < len(recordsList); i++ {
		ns[i] = recordsList[i].Value
	}
	return ns, ErrNotFound
}

func GetLocalNS() (NSServers, error) {
	var ns NSServers
	var rerr error

	//order
	//1. /etc/hosts
	//2. /etc/sysconfig/network-scripts/ifcfg-eth0
	//3. /etc/resolv.conf
	//nic list
	ns.NicNS = make(map[string]NSServer)

	ifaces, err := net.Interfaces()
	if err == nil {
		for i := 0; i < len(ifaces); i++ {
			iface := &ifaces[i]
			//only support centos/redhat;
			//TODO: add debian
			nicns, err := GetRhNsByNic(iface.Name)
			if err == nil {
				ns.NicNS[iface.Name] = nicns
			}
		}
	} else {
		rerr = err
	}

	txt, err := ioutil.ReadFile("/etc/resolv.conf")

	//TODO:
	if err == nil {
		match := resolvNsExp.FindAllSubmatchIndex(txt, 2)
		for i := 0; i < len(match); i++ {
			ns.NSServer = append(ns.NSServer, string(txt[match[i][2]:match[i][3]]))
		}
	} else {
		rerr = err
	}
	return ns, rerr
}
