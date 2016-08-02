package main

import (
	"fmt"

	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/influxdb"
	"github.com/arsgo/lib4go/sysinfo"
)

func main() {
	fmt.Println(sysinfo.GetMemory())
	fmt.Println(sysinfo.GetCPU())
	fmt.Println(sysinfo.GetDisk())
	client := cluster.NewZKClient()

	err := influxdb.SaveMapsToInfluxDB("db", "influxdb001", sysinfo.GetMemory(), client)
	fmt.Println(err)

	err = influxdb.SaveMapsToInfluxDB("db", "influxdb002_cpu", sysinfo.GetCPU(), client)
	fmt.Println(err)

	err = influxdb.SaveMapsToInfluxDB("db", "influxdb003_disk", sysinfo.GetDisk(), client)
	fmt.Println(err)

}
