package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// SysInfo saves the basic system information
type SysInfo struct {
	Hostname  string  `json:hostname`
	Platform  string  `json:platform`
	Uptime    uint64  `json:uptime`
	CPU       string  `json:cpu`
	CPU_mhz   float64 `json:cpu_mhz`
	CPU_c     int32   `json:cpu_c`
	CPU_usage float64 `json:cpu_usage`
	RAM       uint64  `json:ram`
	RAM_p     float64 `json:ram_p`
	RAM_free  uint64  `json:ram_free`
	Disk      uint64  `json:disk`
	Disk_free uint64  `json:disk_free`
	Disk_used uint64  `json:disk_used`
	Disk_p    float64 `json:disk_p`
}

func main() {
	for {
		//receiving data from system
		percent, _ := cpu.Percent(time.Second, false)
		hostStat, _ := host.Info()
		cpuStat, _ := cpu.Info()
		vmStat, _ := mem.VirtualMemory()
		diskStat, _ := disk.Usage("\\") // If you're in Unix change this "\\" for "/"

		m := &SysInfo{Hostname: hostStat.Hostname,
			Platform:  hostStat.Platform,          //	Information about platform
			Uptime:    hostStat.Uptime / 60,       //	Uptime in minutes
			CPU:       cpuStat[0].ModelName,       //	Cpu model
			CPU_mhz:   cpuStat[0].Mhz,             // 	Max cpu clock speed
			CPU_c:     cpuStat[0].Cores,           // 	Cpu cores
			CPU_usage: percent[0],                 //	Cpu usage
			RAM:       vmStat.Total / 1024 / 1024, // 	Total RAM
			RAM_p:     vmStat.UsedPercent,         // 	RAM usage %
			RAM_free:  vmStat.Free / 1024 / 1024,  //	Free RAM
			Disk:      diskStat.Total / 1000000,   //	Total disk space
			Disk_free: diskStat.Free / 1000000,    //	Free disk space
			Disk_used: diskStat.Used / 1000000,    //	Disk used space
			Disk_p:    diskStat.UsedPercent,       //	Disk usage %
		}
		data, _ := json.Marshal(m)
		http.Post("http://127.0.0.1:5000/data", "application/json", bytes.NewBuffer(data))
		fmt.Println(string(data))
		time.Sleep(3 * time.Second)
	}
}
