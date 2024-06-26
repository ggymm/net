package hosts

import (
	"bufio"
	"io"
	"net"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"net/pkg/log"
)

const (
	dnsFile  = "dns.txt"
	tempFile = "dns-temp.txt"
)

func GetDns() {
	url := "https://public-dns.info/nameservers.txt"
	_, err := resty.New().R().SetOutput(dnsFile).Get(url)
	if err != nil {
		log.Error().
			Str("url", url).
			Err(errors.WithStack(err)).Msg("fetch dns error")
		return
	}

	var (
		f1 *os.File
		f2 *os.File
	)
	f1, err = os.OpenFile(dnsFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", dnsFile).
			Err(errors.WithStack(err)).Msg("read dns file error")
		return
	}
	f2, err = os.OpenFile(tempFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", tempFile).
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
	err = os.Rename(tempFile, dnsFile)
	if err != nil {
		log.Error().
			Str("file1", dnsFile).
			Str("file2", tempFile).
			Err(errors.WithStack(err)).Msg("rename dns file error")
		return
	}
}

func ReadDns() []string {
	ips := make([]string, 0)
	dnsF, err := os.OpenFile(dnsFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error().
			Str("file", dnsFile).
			Err(errors.WithStack(err)).Msg("read dns file error")
		return ips
	}
	buf := bufio.NewReader(dnsF)
	for {
		l, _, err1 := buf.ReadLine()
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			continue
		}
		ips = append(ips, string(l))
	}
	_ = dnsF.Close()
	return ips
}
