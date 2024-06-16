package hosts

import (
	"testing"

	"github.com/google/gopacket/pcap"
)

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
