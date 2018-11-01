package net

import (
	"encoding/json"
	"fmt"
	"net"

	//"syscall"
	//"golang.org/x/sys/windows"
	"time"

	"github.com/StackExchange/wmi"
)

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

func GetDefaultGateway() (net.IP, error) {
	var items []Win32_IP4RouteTable
	q := wmi.CreateQuery(&items, `WHERE Destination="0.0.0.0"`)
	err := wmi.Query(q, &items)
	if err != nil {
		return net.ParseIP("0.0.0.0"), err
	}

	for _, it := range items {
		return net.ParseIP(it.NextHop), nil
	}
	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func GetGatewayByNic(nicName string) (net.IP, error) {
	ifi, err := net.InterfaceByName(nicName)
	if err != nil {
		return net.ParseIP("0.0.0.0"), err
	}

	var items []Win32_IP4RouteTable
	q := wmi.CreateQuery(&items, fmt.Sprintf(`WHERE InterfaceIndex=%d`, ifi.Index))
	err = wmi.Query(q, &items)
	if err != nil {
		return net.ParseIP("0.0.0.0"), err
	}

	for _, it := range items {
		return net.ParseIP(it.NextHop), nil
	}
	return net.ParseIP("0.0.0.0"), ErrNotFound
}

func ListGateway() ([]*RouteItem, error) {
	var routes []*RouteItem
	var items []Win32_IP4RouteTable
	q := wmi.CreateQuery(&items, `WHERE Destination!="255.255.255.255"`)
	err := wmi.Query(q, &items)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		route := new(RouteItem)
		route.Gateway = net.ParseIP(it.NextHop)
		route.Dst = net.ParseIP(it.Destination)
		ifi, err := net.InterfaceByIndex(int(it.InterfaceIndex))
		if err == nil {
			route.Iface = ifi.Name
			route.MTU = ifi.MTU
		}
		route.Flags = uint32(ifi.Flags)
		route.Mask = netStringToIPv4Mask(it.Mask)
		route.Metric = int(it.Metric1)
	}

	return routes, nil
}
