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
	dev string

	loIp net.IP           // 本地 ip 地址
	loHw net.HardwareAddr // 本地 mac 地址
	gwIp net.IP           // 网关 ip 地址
	gwHw net.HardwareAddr // 网关 mac 地址

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

func (s *Scanner) Init(d Device) error {
	s.dev = d.Name
	s.loIp = d.IpAddr
	s.loHw = d.HwAddr
	return s.netInfo(d)
}

func (s *Scanner) netInfo(d Device) (err error) {
	// 获取网关 mac 地址
	s.gwIp, err = ParseGateway(s.loIp)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Err(errors.WithStack(err)).Msg("parse gateway ip error")
		return
	}

	en := &layers.Ethernet{
		SrcMAC:       s.loHw,
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
		SourceHwAddress:   []byte(s.loHw),
		SourceProtAddress: []byte(s.loIp),
	}

	s.handle, err = pcap.OpenLive(s.dev, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Error().
			Err(errors.WithStack(err)).Msg("open device error")
		return
	}
	err = s.sendPacket(en, arp)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Str("4.gwIp", s.gwIp.String()).
			Err(errors.WithStack(err)).Msg("write packet data error")
	}
	//defer s.handle.Close()

	for {
		var data []byte
		data, _, err = s.handle.ReadPacketData()
		if err != nil {
			if errors.Is(err, pcap.NextErrorTimeoutExpired) {
				continue
			}
			log.Error().
				Str("1.dev", s.dev).
				Str("2.loIp", s.loIp.String()).
				Str("3.loHw", s.loHw.String()).
				Str("4.gwIp", s.gwIp.String()).
				Err(errors.WithStack(err)).Msg("read packet data error")
			return
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if packet == nil {
			continue
		}
		if layer := packet.Layer(layers.LayerTypeARP); layer != nil {
			arp = layer.(*layers.ARP)
			if arp.Operation == layers.ARPReply && net.IP(arp.SourceProtAddress).Equal(s.gwIp) {
				s.gwHw = arp.SourceHwAddress
				break
			}
		}
	}
	return
}

func (s *Scanner) sendPacket(l ...gopacket.SerializableLayer) error {
	buf := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buf, s.options, l...)
	if err != nil {
		log.Error().
			Str("1.dev", s.dev).
			Str("2.loIp", s.loIp.String()).
			Str("3.loHw", s.loHw.String()).
			Str("4.gwIp", s.gwIp.String()).
			Err(errors.WithStack(err)).Msg("serialize packet error")
		return err
	}
	return s.handle.WritePacketData(buf.Bytes())
}
