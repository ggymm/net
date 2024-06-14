package hosts

import "testing"

func Test_netInfo(t *testing.T) {
	t.Log(netInfo())
}

func Test_ReadIps(t *testing.T) {
	ReadIps("github.com.", "8.8.8.8:53")
}
