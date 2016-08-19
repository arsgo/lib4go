package utility

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"io"
	"net"
	"path/filepath"
	"strings"

	"github.com/arsgo/lib4go/security/md5"
)

func GetSessionID() string {
	id := GetGUID()
	return id[:8]
}

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

func Escape(input string) string {
	r := strings.Replace(input, "\\u0026", "&", -1)
	r = strings.Replace(r, "\\u003c", "<", -1)
	r = strings.Replace(r, "\\u003e", ">", -1)
	return r
}

func GetExcPath(p ...string) string {
	if len(p) == 0 {
		return ""
	}
	if strings.HasPrefix(p[0], ".") {
		path, err := getExecRoot()
		if err != nil {
			return p[0]
		}
		for i := 1; i < len(p); i++ {
			path = strings.Trim(path, p[i])
		}
		return filepath.Join(path, strings.Trim(p[0], "."))
	}
	return p[0]
}
func Clone(src interface{}) (dst interface{}, err error) {
	var buf bytes.Buffer
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		return
	}
	err = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
	return
}
