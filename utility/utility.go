package utility

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	"strings"

	"github.com/colinyl/lib4go/security/md5"
)

//GetGuid 生成Guid字串
func GetGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}

func GetLocalIPAddress(masks ...string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	var ipLst []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipLst = append(ipLst, ipnet.IP.String())
		}
	}
	if len(masks) == 0 && len(ipLst) > 0 {
		return ipLst[0]
	}
	for _, ip := range ipLst {
		for _, m := range masks {
			if strings.HasPrefix(ip, m) {
				return ip
			}
		}
	}
	return "127.0.0.1"
}
