package hosts

import (
	"fmt"
	"net"
)

type Device struct {
	Name   string // 设备名称
	IpAddr net.IP
	HwAddr net.HardwareAddr

	DevName string // pcap 设备名称
	DevDesc string // pcap 设备描述
}

func (d *Device) String() string {
	return fmt.Sprintf("%s - %s - %s - %s", d.DevDesc, d.IpAddr, d.HwAddr, d.DevName)
}
