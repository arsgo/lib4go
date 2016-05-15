package sysinfo

import (
	"encoding/json"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func GetMemory() []map[string]interface{} {
	v, _ := mem.VirtualMemory()
	data := make(map[string]interface{})
	buffer, _ := json.Marshal(&v)
	json.Unmarshal(buffer, &data)
	var result []map[string]interface{}
	result = append(result, data)
	return result
}

func GetCPU() []map[string]interface{} {
	v, _ := cpu.Times(true)
	buffer, _ := json.Marshal(&v)
	var data []map[string]interface{}
	json.Unmarshal(buffer, &data)
	return data
}
func GetDisk() []map[string]interface{} {
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
	return data
}
