package hosts

import (
	"bufio"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/pkg/log"
	"os"
)

const (
	dnsUrl = "https://public-dns.info/nameservers.txt"
)

func GetDNS() {
	dnsF := "dns.txt"
	dnsT := "dns-temp.txt"
	_, err := resty.New().R().SetOutput(dnsF).Get(dnsUrl)
	if err != nil {
		log.Error().
			Str("url", dnsUrl).
			Err(errors.WithStack(err)).Msg("fetch dns error")
		return
	}

	var (
		f1 *os.File
		f2 *os.File
	)
	f1, err = os.OpenFile(dnsF, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", dnsF).
			Err(errors.WithStack(err)).Msg("read dns file error")
		return
	}
	f2, err = os.OpenFile(dnsT, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", dnsT).
			Err(errors.WithStack(err)).Msg("create dns file error")
	}
	buf := bufio.NewReader(f1)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		ip := net.ParseIP(string(l))
		if ip.To4() != nil {
			_, _ = f2.Write(l)
			_, _ = f2.WriteString("\n")
		}
	}
	_ = f1.Close()
	_ = f2.Close()

	// 重命名文件
	err = os.Rename(dnsT, dnsF)
	if err != nil {
		log.Error().
			Str("file", dnsF).
			Str("temp", dnsT).
			Err(errors.WithStack(err)).Msg("rename dns file error")
		return
	}
}
