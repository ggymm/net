package hosts

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"net"
)

var (
	github = []string{
		"https://github.com/",
	}
)

func sendPacket(s *Scanner, domain, nameserver string) {
	en := &layers.Ethernet{
		SrcMAC:       s.loHw,
		DstMAC:       s.gwHw,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    s.loIp.To4(),
		DstIP:    net.ParseIP(nameserver).To4(),
		Protocol: layers.IPProtocolTCP,
	}
	udp := &layers.UDP{
		SrcPort: 12345,
		DstPort: 53,
	}
	_ = udp.SetNetworkLayerForChecksum(ip)
	dns := &layers.DNS{
		ID:           1,
		QR:           false,
		OpCode:       layers.DNSOpCodeQuery,
		AA:           false,
		TC:           false,
		RD:           true,
		RA:           false,
		Z:            2,
		ResponseCode: layers.DNSResponseCodeNoErr,
		QDCount:      1,
		ANCount:      0,
		NSCount:      0,
		ARCount:      0,
		Questions: []layers.DNSQuestion{
			{
				Name:  []byte(domain),
				Type:  layers.DNSTypeA,
				Class: layers.DNSClassIN,
			},
		},
	}
	err := s.sendPacket(en, ip, udp, dns)
	if err != nil {
		panic(err)
	}
}

func ReadIps(domain, nameserver string) {
	d := Device{
		IpAddr: net.IP{192, 168, 1, 27},
	}
	mac, _ := net.ParseMAC("54:05:db:83:7f:a5")
	d.HwAddr = mac
	d.Name = `\Device\NPF_{81A86FFA-2C4F-4E6B-AD4E-29036647FB75}`
	d.Desc = "Realtek PCIe GbE Family Controller"

	s := NewScanner()
	err := s.Init(d)
	if err != nil {
		panic(err)
	}

	go func() {
		var handle *pcap.Handle
		handle, err = pcap.OpenLive(d.Name, 65535, true, pcap.BlockForever)
		if err != nil {
			panic(err)
		}
		err = handle.SetBPFFilter("udp and port 53")
		if err != nil {
			panic(err)
		}
		for {
			var data []byte
			data, _, err = s.handle.ReadPacketData()
			if err != nil {
				if errors.Is(err, pcap.NextErrorTimeoutExpired) {
					continue
				}
				panic(err)
			}

			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
			if packet == nil {
				continue
			}
			if layer := packet.Layer(layers.LayerTypeDNS); layer != nil {
				dns := layer.(*layers.DNS)
				for _, answer := range dns.Answers {
					if answer.Type == layers.DNSTypeA {
						fmt.Println("ip", answer.IP.String())
					}
				}
			}
		}
	}()
	sendPacket(s, domain, nameserver)
}
