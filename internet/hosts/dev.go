package hosts

import (
	"fmt"
	"net"
)

type Device struct {
	Name string // pcap 设备名称
	Desc string // pcap 设备描述

	IpAddr net.IP
	HwAddr net.HardwareAddr
}

func (d *Device) String() string {
	return fmt.Sprintf("%s - %s - %s - %s", d.Desc, d.IpAddr, d.HwAddr, d.Name)
}
