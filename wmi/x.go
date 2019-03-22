package wmi

type NicConfigListBrief struct {
	DefaultIPGateway string
	Description      string
	DhcpEnable       bool
	DNSDomain        string
	Index            uint
	IPAddress        string
}
