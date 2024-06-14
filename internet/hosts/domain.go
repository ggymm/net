package hosts

import (
	"net"

	"github.com/pkg/errors"

	"net/pkg/log"
)

var (
	device = "以太网"
	github = []string{
		"https://github.com/",
	}
)

func netInfo() (ip net.IP, mac net.HardwareAddr) {
	ifs, err := net.Interfaces()
	if err != nil {
		log.Error().Err(errors.WithStack(err)).Msg("query interface error")
		return
	}
	for _, i := range ifs {
		if i.Name == device {
			// 获取 ip 地址
			addrs, err1 := i.Addrs()
			if err1 != nil {
				log.Error().Err(errors.WithStack(err1)).Msg("query interface address error")
				return
			}
			for _, addr := range addrs {
				ipn, ok := addr.(*net.IPNet)
				if ok && !ipn.IP.IsLoopback() {
					ip = ipn.IP.To4()
					break
				}
			}

			// 获取 mac 地址
			mac = i.HardwareAddr
			break
		}
	}
	return
}

func sendPack() {

}

func ReadIps(domain, nameserver string) {

}
