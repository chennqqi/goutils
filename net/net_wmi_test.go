// +build windows

package net

import (
	"encoding/json"
	"fmt"
	"testing"

	stdwmi "github.com/StackExchange/wmi"
)

func TestGetNetworkParams(t *testing.T) {
	ip, err := getNetworkParams()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("HOST NAME: %s\n", ip.HostName)
	fmt.Printf("DOMAIn NAME: %s\n", ip.DomainName)
	fmt.Printf("DNS Servers:\n")
	fmt.Printf("\t%s\n", ip.DnsServerList.IpAddress.String)
	pIPAddr := ip.DnsServerList.Next
	for pIPAddr != nil {
		fmt.Printf("\t%s\n", pIPAddr.IpAddress.String)
		pIPAddr = pIPAddr.Next
	}
	fmt.Printf("\tNode Type . . . . . . . . . : ")
	switch ip.NodeType {
	case 1:
		fmt.Printf("%s\n", "Broadcast")
	case 2:
		fmt.Printf("%s\n", "Peer to peer")

	case 4:
		fmt.Printf("%s\n", "Mixed")

	case 8:
		fmt.Printf("%s\n", "Hybrid")

	default:
		fmt.Printf("\n")
	}

	fmt.Printf("\tNetBIOS Scope ID. . . . . . : %s\n", ip.ScopeId)
	if ip.EnableRouting == 1 {
		fmt.Printf("\tIP Routing Enabled. . . . . : yes\n")
	} else {
		fmt.Printf("\tIP Routing Enabled. . . . . : no\n")
	}
	if ip.EnableProxy == 1 {
		fmt.Printf("\tWINS Proxy Enabled. . . . . : %s\n", "yes")
	} else {
		fmt.Printf("\tWINS Proxy Enabled. . . . . : %s\n", "no")
	}
	if ip.EnableDns == 1 {
		fmt.Printf("\tNetBIOS Resolution Uses DNS : %s\n", "yes")

	} else {
		fmt.Printf("\tNetBIOS Resolution Uses DNS : %s\n", "no")
	}
}

type MicrosoftDNS_Domain struct {
	DnsServerName string
	ContainerName string
	Name          string
}

type MicrosoftDNS_Server struct {
	Name                      string
	Version                   uint32
	LogLevel                  uint32
	LogFilePath               string
	LogFileMaxSize            uint32
	LogIPFilterList           []string
	EventLogLevel             uint32
	RpcProtocol               int32
	NameCheckFlag             uint32
	AddressAnswerLimit        uint32
	RecursionRetry            uint32
	RecursionTimeout          uint32
	DsPollingInterval         uint32
	DsTombstoneInterval       uint32
	MaxCacheTTL               uint32
	MaxNegativeCacheTTL       uint32
	SendPort                  uint32
	XfrConnectTimeout         uint32
	BootMethod                uint32
	AllowUpdate               uint32
	UpdateOptions             uint32
	DsAvailable               bool
	DisableAutoReverseZones   bool
	AutoCacheUpdate           bool
	NoRecursion               bool
	RoundRobin                bool
	LocalNetPriority          bool
	StrictFileParsing         bool
	LooseWildcarding          bool
	BindSecondaries           bool
	WriteAuthorityNS          bool
	ForwardDelegations        uint32
	SecureResponses           bool
	DisjointNets              bool
	AutoConfigFileZones       uint32
	ScavengingInterval        uint32
	DefaultRefreshInterval    uint32
	DefaultNoRefreshInterval  uint32
	DefaultAgingState         bool
	EDnsCacheTimeout          uint32
	EnableEDnsProbes          bool
	EnableDnsSec              uint32
	ServerAddresses           []string
	ListenAddresses           []string
	Forwarders                []string
	ForwardingTimeout         uint32
	IsSlave                   bool
	EnableDirectoryPartitions bool
}

func (t MicrosoftDNS_Server) String() string {
	txt, _ := json.MarshalIndent(&t, " ", "\t")
	return string(txt)
}

func (t MicrosoftDNS_Domain) String() string {
	txt, _ := json.MarshalIndent(&t, " ", "\t")
	return string(txt)
}

func TestRouteTable(t *testing.T) {
	return
	var items []Win32_IP4RouteTable
	err := stdwmi.Query(
		`SELECT * FROM Win32_IP4RouteTable`, &items)
	if err != nil {
		t.Error(err)
	}
	for i, v := range items {
		fmt.Println(i, v)
	}
}

func TestNICCONFIG(t *testing.T) {
	return
	var items []Win32_NetworkAdapterConfiguration
	q := stdwmi.CreateQuery(&items, "WHERE Index=14")
	err := stdwmi.Query(q, &items)
	if err != nil {
		t.Error(err)
	}
	for i, v := range items {
		if v.Index == 14 {
			fmt.Println(i, v)
		}
	}
}

func TestNIC(t *testing.T) {
	return //OK
	var items []Win32_NetworkAdapter
	q := stdwmi.CreateQuery(&items, "")
	err := stdwmi.Query(q, &items)
	if err != nil {
		t.Error(err)
	}
	for i, v := range items {
		if v.Index == 14 {
			fmt.Println(i, v)
		}
	}
}

func TestDnsDomain(t *testing.T) {
	return //ERROR
	var items []MicrosoftDNS_Domain
	q := stdwmi.CreateQuery(&items, "")
	err := stdwmi.Query(q, &items)
	if err != nil {
		t.Error(err)
	}
	for i, v := range items {
		fmt.Println(i, v)
	}
}

func TestDnsServer(t *testing.T) {
	return //ERROR
	var items []MicrosoftDNS_Server
	q := stdwmi.CreateQuery(&items, "")
	err := stdwmi.Query(q, &items)
	if err != nil {
		t.Error(err)
	}
	for i, v := range items {
		fmt.Println(i, v)
	}
}
