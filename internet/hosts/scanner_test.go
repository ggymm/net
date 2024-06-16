package hosts

import (
	"net"
	"testing"

	"github.com/google/gopacket/pcap"

	"net/pkg/app"
	"net/pkg/log"
)

func Test_NetInfo(t *testing.T) {
	app.Init()
	log.Init("hosts")

	dev := Device{
		Name:   "eth0",
		IpAddr: net.IP{192, 168, 1, 27},
	}
	mac, err := net.ParseMAC("54:05:db:83:7f:a5")
	if err != nil {
		t.Fatal(err)
	}
	dev.HwAddr = mac
	dev.DevName = `\Device\NPF_{81A86FFA-2C4F-4E6B-AD4E-29036647FB75}`
	dev.DevDesc = "Realtek PCIe GbE Family Controller"

	s := NewScanner()
	s.Init(dev)

	t.Logf("gwIp: %s, gwMac: %s", s.gwIp, s.gwMac)
}

func Test_FindDevice(t *testing.T) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		t.Fatal(err)
	}
	for i, dev := range devs {
		t.Log(i, dev)
		for i2, address := range dev.Addresses {
			t.Log(i2, address)
		}
	}
}
