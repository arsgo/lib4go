package utility

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

//GetGuid 生成Guid字串
func GetGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func Md5(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

var localIP string

type DataMap map[string]string

func NewDataMap() DataMap {
	return make(map[string]string)
}

//Add 添加变量
func (d DataMap) Set(k string, v string) {
	d[fmt.Sprintf("@%s", k)] = v
}

//Merge merge new map from current
func (d DataMap) Merge(n DataMap) DataMap {
	nmap := NewDataMap()
	for k, v := range d {
		nmap[k] = v
	}
	for k, v := range n {
		nmap[k] = v
	}
	return nmap
}

//Copy Copy the current map to another
func (d DataMap) Copy() DataMap {
	nmap := NewDataMap()
	for k, v := range d {
		nmap[k] = v
	}
	return nmap
}

//Translate 翻译带有@变量的字符串
func (d DataMap) Translate(format string) string {

	brackets, _ := regexp.Compile(`\{@\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		return d[s[1:len(s)-1]]
	})
	word, _ := regexp.Compile(`@\w+`)
	if word == nil {
		return format
	}
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		return d[s]
	})
	return result

}

//GetLocalIP 获取本机IP地址
func GetLocalIP(mask string) string {
	if localIP == "" {
		localIP = getLocalIP(mask)
	}
	return localIP
}

//------------------------------------------内部函数-----------------------------------
func getLocalIP(mask string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), mask) {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
