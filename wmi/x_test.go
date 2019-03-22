package wmi

import (
	"encoding/json"
	"fmt"
	"syscall"
	"testing"
	"time"
	"unsafe"

	stdwmi "github.com/StackExchange/wmi"
)

var (
	iphlpapi             = syscall.NewLazyDLL("iphlpapi.dll")
	procGetNetworkParams = iphlpapi.NewProc("GetNetworkParams")
)

const MAX_HOSTNAME_LEN = 128    // arb.
const MAX_DOMAIN_NAME_LEN = 128 // arb.
const MAX_SCOPE_ID_LEN = 256    // arb.

type IpAddapterParams struct {
	HostName         [MAX_HOSTNAME_LEN + 4]byte
	DomainName       [MAX_DOMAIN_NAME_LEN + 4]byte
	CurrentDnsServer *syscall.IpAddrString
	DnsServerList    syscall.IpAddrString
	NodeType         uint32
	ScopeId          [MAX_SCOPE_ID_LEN + 4]byte
	EnableRouting    uint32
	EnableProxy      uint32
	EnableDns        uint32
}

func getNetworkParams() (*IpAddapterParams, error) {
	b := make([]byte, 2048)
	l := uint32(len(b))
	a := (*IpAddapterParams)(unsafe.Pointer(&b[0]))

	// TODO(mikio): GetAdaptersInfo returns IP_ADAPTER_INFO that
	// contains IPv4 address list only. We should use another API
	// for fetching IPv6 stuff from the kernel.
	r0, _, _ := syscall.Syscall(procGetNetworkParams.Addr(), 2, uintptr(unsafe.Pointer(a)), uintptr(unsafe.Pointer(&l)), 0)
	if r0 != 0 {
		return nil, syscall.Errno(r0)
	}
	return a, nil
}

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

type Win32_IP4RouteTable struct {
	Age            int32
	Caption        string
	Description    string
	Destination    string
	Information    string
	InstallDate    time.Time
	InterfaceIndex int32
	Mask           string
	Metric1        int32
	Metric2        int32
	Metric3        int32
	Metric4        int32
	Metric5        int32
	Name           string
	NextHop        string
	Protocol       uint32
	Status         string
	Type           uint32
}

type Win32_NetworkAdapterConfiguration struct {
	Caption                      string
	Description                  string
	SettingID                    string
	ArpAlwaysSourceRoute         bool
	ArpUseEtherSNAP              bool
	DatabasePath                 string
	DeadGWDetectEnabled          bool
	DefaultIPGateway             *string
	DefaultTOS                   uint8
	DefaultTTL                   uint8
	DHCPEnabled                  bool
	DHCPLeaseExpires             time.Time
	DHCPLeaseObtained            time.Time
	DHCPServer                   string
	DNSDomain                    string
	DNSDomainSuffixSearchOrder   *string
	DNSEnabledForWINSResolution  bool
	DNSHostName                  string
	DNSServerSearchOrder         *string
	DomainDNSRegistrationEnabled bool
	ForwardBufferMemory          uint32
	FullDNSRegistrationEnabled   bool
	GatewayCostMetric            []uint16
	IGMPLevel                    uint8
	Index                        uint32
	InterfaceIndex               uint32
	IPAddress                    []string
	IPConnectionMetric           uint32
	IPEnabled                    bool
	IPFilterSecurityEnabled      bool
	IPPortSecurityEnabled        bool
	IPSecPermitIPProtocols       []string
	IPSecPermitTCPPorts          []string
	IPSecPermitUDPPorts          []string
	IPSubnet                     []string
	IPUseZeroBroadcast           bool
	IPXAddress                   string
	IPXEnabled                   bool
	IPXFrameType                 []uint32
	IPXMediaType                 uint32
	IPXNetworkNumber             []string
	IPXVirtualNetNumber          string
	KeepAliveInterval            uint32
	KeepAliveTime                uint32
	MACAddress                   string
	MTU                          uint32
	NumForwardPackets            uint32
	PMTUBHDetectEnabled          bool
	PMTUDiscoveryEnabled         bool
	ServiceName                  string
	TcpipNetbiosOptions          uint32
	TcpMaxConnectRetransmissions uint32
	TcpMaxDataRetransmissions    uint32
	TcpNumConnections            uint32
	TcpUseRFC1122UrgentPointer   bool
	TcpWindowSize                uint16
	WINSEnableLMHostsLookup      bool
	WINSHostLookupFile           string
	WINSPrimaryServer            string
	WINSScopeID                  string
	WINSSecondaryServer          string
}

type Win32_NetworkAdapter struct {
	AdapterType                 string
	AdapterTypeID               uint16
	AutoSense                   bool
	Availability                uint16
	Caption                     string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	Description                 string
	DeviceID                    string
	ErrorCleared                bool
	ErrorDescription            string
	GUID                        string
	Index                       uint32
	InstallDate                 time.Time
	Installed                   bool
	InterfaceIndex              uint32
	LastErrorCode               uint32
	MACAddress                  string
	Manufacturer                string
	MaxNumberControlled         uint32
	MaxSpeed                    uint64
	Name                        string
	NetConnectionID             string
	NetConnectionStatus         uint16
	NetEnabled                  bool
	NetworkAddresses            []string
	PermanentAddress            string
	PhysicalAdapter             bool
	PNPDeviceID                 string
	PowerManagementCapabilities []uint16
	PowerManagementSupported    bool
	ProductName                 string
	ServiceName                 string
	Speed                       uint64
	Status                      string
	StatusInfo                  uint16
	SystemCreationClassName     string
	SystemName                  string
	TimeOfLastReset             time.Time
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

func (t Win32_NetworkAdapter) String() string {
	txt, _ := json.MarshalIndent(&t, " ", "\t")
	return string(txt)
}

func (t Win32_IP4RouteTable) String() string {
	txt, _ := json.MarshalIndent(&t, " ", "\t")
	return string(txt)
}

func (t Win32_NetworkAdapterConfiguration) String() string {
	txt, _ := json.MarshalIndent(&t, " ", "\t")
	return string(txt)
}

func TestRouteTable(t *testing.T) {
	var items []Win32_IP4RouteTable

	err := stdwmi.Query(
		`SELECT * FROM Win32_IP4RouteTable WHERE InterfaceIndex=14`, &items)
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
