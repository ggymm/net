package hosts

import (
	"net"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"net/pkg/log"
)

var (
	//device = "eth0"
	device = "WLAN"
	github = []string{
		"https://github.com/",
	}
)

func sendPack(ip net.IP, mac net.HardwareAddr) {
	var (
		snapLen = int32(1024)
		timeout = 10 * time.Second
	)
	handle, err := pcap.OpenLive(device, snapLen, false, timeout)
	if err != nil {
		log.Error().Err(errors.WithStack(err)).Msg("open device error")
		return
	}
	defer handle.Close()

	// 以太网帧头部
	eth := &layers.Ethernet{
		SrcMAC: mac,
	}
	println(eth)
}

func ReadIps(domain, nameserver string) {
}
