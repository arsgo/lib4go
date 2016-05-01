package main

import (
	"fmt"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/influxdb"
	"github.com/colinyl/ars/sysinfo"
)

func main() {
	fmt.Println(sys.GetMemory())
	fmt.Println(sys.GetCPU())
	fmt.Println(sys.GetDisk())
	client := cluster.NewZKClient()

	err := influxdb.SaveMapsToInfluxDB("db", "influxdb001", sys.GetMemoryMap(), client)	
	fmt.Println(err)
	
	err = influxdb.SaveMapsToInfluxDB("db", "influxdb002_cpu", sys.GetCPUMap(), client)	
	fmt.Println(err)
	
	err = influxdb.SaveMapsToInfluxDB("db", "influxdb003_disk", sys.GetDiskMap(), client)	
	fmt.Println(err)

}
