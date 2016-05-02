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
	"sync"
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

type DataMap struct {
	data map[string]string
	lk   sync.Mutex
}

func NewDataMap() *DataMap {
	return &DataMap{data: make(map[string]string)}
}
func NewDataMaps(d map[string]interface{}) *DataMap {
	current := make(map[string]string)
	for k, v := range d {
		current[fmt.Sprintf("@%s", k)] = fmt.Sprint(v)
	}
	return &DataMap{data: current}
}

//Add 添加变量
func (d *DataMap) Set(k string, v string) {
	d.lk.Lock()
	defer d.lk.Unlock()
	d.data[fmt.Sprintf("@%s", k)] = v
}
func (d *DataMap) Get(k string) string {
	d.lk.Lock()
	defer d.lk.Unlock()
	return d.data[fmt.Sprintf("@%s", k)]
}

//Merge merge new map from current
func (d *DataMap) Merge(n DataMap) *DataMap {
	d.lk.Lock()
	defer d.lk.Unlock()
	nmap := NewDataMap()
	MergeStringMap(d.data, nmap.data)
	MergeStringMap(n.data, nmap.data)
	return nmap
}

//Copy Copy the current map to another
func (d *DataMap) Copy() *DataMap {
	d.lk.Lock()
	defer d.lk.Unlock()
	nmap := NewDataMap()
	for k, v := range d.data {
		nmap.data[k] = v
	}
	return nmap
}

//Translate 翻译带有@变量的字符串
func (d *DataMap) Translate(format string) string {
	d.lk.Lock()
	defer d.lk.Unlock()
	brackets, _ := regexp.Compile(`\{@\w+\}`)
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		return d.data[s[1:len(s)-1]]
	})
	word, _ := regexp.Compile(`@\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		return d.data[s]
	})
	return result
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
func MergeMaps(source map[string]interface{}, targets []map[string]interface{}) []map[string]interface{} {
	for k, v := range source {
		for _, target := range targets {
			target[k] = v
		}
	}
	return targets
}
func MergeMap(source map[string]interface{}, target map[string]interface{}) map[string]interface{} {
	for k, v := range source {
		target[k] = v
	}
	return target
}
func MergeStringMap(source map[string]string, target map[string]string) map[string]string {
	for k, v := range source {
		target[k] = v
	}
	return target
}
