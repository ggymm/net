package hosts

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"net/pkg/log"
)

type Scanner struct {
	gwIp  net.IP           // 网关 ip 地址
	gwMac net.HardwareAddr // 网关 mac 地址

	handle  *pcap.Handle
	options gopacket.SerializeOptions
}

func NewScanner() *Scanner {
	return &Scanner{
		options: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
	}
}

func (s *Scanner) Init(d Device) {
	s.netInfo(d)
}

func (s *Scanner) netInfo(d Device) {
	var (
		err    error
		srcIp  = d.IpAddr
		srcMac = d.HwAddr
	)

	// 获取网关 mac 地址
	s.gwIp, err = ParseGateway(srcIp)
	if err != nil {
		log.Error().
			Str("1.dev", d.Name).
			Str("2.srcIp", srcIp.String()).
			Str("3.srcMac", srcMac.String()).
			Err(errors.WithStack(err)).Msg("parse gateway ip error")
		return
	}

	eth := &layers.Ethernet{
		SrcMAC:       srcMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		Operation:         layers.ARPRequest,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(s.gwIp),
		SourceHwAddress:   []byte(srcMac),
		SourceProtAddress: []byte(srcIp),
	}

	s.handle, err = pcap.OpenLive(d.DevName, 100, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Str("1.dev", d.Name).
			Str("2.name", d.DevName).
			Err(errors.WithStack(err)).Msg("open device error")
		return
	}
	buf := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buf, s.options, eth, arp)
	if err != nil {
		log.Error().
			Str("1.dev", d.Name).
			Str("2.name", d.DevName).
			Str("3.gwIp", s.gwIp.String()).
			Str("4.srcIp", srcIp.String()).
			Str("5.srcMac", srcMac.String()).
			Err(errors.WithStack(err)).Msg("serialize arp packet error")
	}
	err = s.handle.WritePacketData(buf.Bytes())
	if err != nil {
		log.Error().
			Str("1.dev", d.Name).
			Str("2.name", d.DevName).
			Str("3.gwIp", s.gwIp.String()).
			Str("4.srcIp", srcIp.String()).
			Str("5.srcMac", srcMac.String()).
			Err(errors.WithStack(err)).Msg("write packet data error")
	}
	defer s.handle.Close()

	for {
		var data []byte
		data, _, err = s.handle.ReadPacketData()
		if err != nil {
			if errors.Is(err, pcap.NextErrorTimeoutExpired) {
				continue
			}
			log.Error().
				Str("1.dev", d.Name).
				Str("2.name", d.DevName).
				Str("3.gwIp", s.gwIp.String()).
				Str("4.srcIp", srcIp.String()).
				Str("5.srcMac", srcMac.String()).
				Err(errors.WithStack(err)).Msg("read packet data error")
			return
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if packet == nil {
			continue
		}
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp = arpLayer.(*layers.ARP)
			if arp.Operation == layers.ARPReply && net.IP(arp.SourceProtAddress).Equal(s.gwIp) {
				s.gwMac = arp.SourceHwAddress
				break
			}
		}
	}
	return
}
