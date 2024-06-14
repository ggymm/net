package hosts

import (
	"net/pkg/app"
	"net/pkg/log"
	"os"
	"testing"
)

func Test_GetDns(t *testing.T) {
	app.Init()
	log.Init("hosts")

	GetDns()

	f, err := os.ReadFile("dns.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(f))
}

func Test_ReadDns(t *testing.T) {
	app.Init()
	log.Init("hosts")

	ips := ReadDns()
	for i, ip := range ips {
		t.Log(i, ip)
	}
}
