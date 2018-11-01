package net

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	iphlpapi             = syscall.NewLazyDLL("iphlpapi.dll")
	procGetNetworkParams = iphlpapi.NewProc("GetNetworkParams")
)

func getAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 2048)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))

	// TODO(mikio): GetAdaptersInfo returns IP_ADAPTER_INFO that
	// contains IPv4 address list only. We should use another API
	// for fetching IPv6 stuff from the kernel.

	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

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

func GetLocalNS() (NSServers, error) {
	var ns NSServers
	var rerr error

	//windows only support one couple dns
	iphelper, err := getNetworkParams()
	if err != nil {
		return ns, err
	}
	nsip := strings.Trim(fmt.Sprintf(`%s`, iphelper.DnsServerList.IpAddress.String), "\t ")
	ns.NSServer = append(ns.NSServer, nsip)
	pIPAddr := iphelper.DnsServerList.Next
	for pIPAddr != nil {
		nsip = strings.Trim(fmt.Sprintf(`%s`, pIPAddr.IpAddress.String), "\t ")
		ns.NSServer = append(ns.NSServer, nsip)
		pIPAddr = pIPAddr.Next
	}
	return ns, rerr
}
