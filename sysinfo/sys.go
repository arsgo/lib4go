package sysinfo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/colinyl/lib4go/utility"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type stat struct {
	Data interface{} `json:"data"`
	IP   string      `json:"ip"`
	timestamp int64       `json:"timestamp"`
}

func GetMemory() string {
	v, _ := mem.VirtualMemory()
	data, _ := json.Marshal(&stat{v, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}

func GetCPU() string {
	v, _ := cpu.Times(true)
	data, _ := json.Marshal(&stat{v, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}

func GetDisk() string {
	var stats []*disk.UsageStat
	if runtime.GOOS == "windows" {
		v, _ := disk.Partitions(true)
		for _, p := range v {
			s, _ := disk.Usage(p.Device)
			stats = append(stats, s)
		}
	} else {
		s, _ := disk.Usage("/")
		stats = append(stats, s)
	}
	data, _ := json.Marshal(&stat{stats, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}

func GetMemoryMap() []map[string]interface{} {
	v, _ := mem.VirtualMemory()
	data := make(map[string]interface{})
	buffer, _ := json.Marshal(&v)
	json.Unmarshal(buffer, &data)
	data["ip"] = utility.GetLocalIP("192")
	data["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	var result []map[string]interface{}
	result = append(result, data)
	return result
}

func GetCPUMap() []map[string]interface{} {
	v, _ := cpu.Times(true)
	buffer, _ := json.Marshal(&v)
	var data []map[string]interface{}
	json.Unmarshal(buffer, &data)
	var ip = utility.GetLocalIP("192")
	var timestamp = fmt.Sprintf("%d", time.Now().Unix())

	for _, v := range data {
		v["ip"] = ip
		v["timestamp"] = timestamp
	}
	return data
}
func GetDiskMap() []map[string]interface{} {
	var stats []*disk.UsageStat
	if runtime.GOOS == "windows" {
		v, _ := disk.Partitions(true)
		for _, p := range v {
			s, _ := disk.Usage(p.Device)
			stats = append(stats, s)
		}
	} else {
		s, _ := disk.Usage("/")
		stats = append(stats, s)
	}

	buffer, _ := json.Marshal(&stats)
	var data []map[string]interface{}
	json.Unmarshal(buffer, &data)

	var ip = utility.GetLocalIP("192")
	var timestamp = fmt.Sprintf("%d", time.Now().Unix())

	for _, v := range data {
		v["ip"] = ip
		v["timestamp"] = timestamp
	}

	return data

}
