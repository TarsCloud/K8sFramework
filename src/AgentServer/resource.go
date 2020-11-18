package main

import (
	"fmt"
	psCpu "github.com/shirou/gopsutil/cpu"
	psDisk "github.com/shirou/gopsutil/disk"
	psHost "github.com/shirou/gopsutil/host"
	psLoad "github.com/shirou/gopsutil/load"
	psMem "github.com/shirou/gopsutil/mem"
	"net"
	"regexp"
	"time"
)

type CPUDetail struct {
	Cores 	int		`json:"cores"`
	Percent float64	`json:"percent"`
}

type CpuInfo struct {
	CpuDetail    []CPUDetail    `json:"cpu_detail"`
	CurrentLoads [3]float64     `json:"current_loads"`
}

type MemInfo struct {
	Total       int64          `json:"total"`
	Available   int64          `json:"available"`
	Used        int64          `json:"used"`
	UsedPercent float64        `json:"used_percent"`
}

type HostInfo struct {
	Hostname        string         `json:"hostname"`
	ProcessNumber   uint64         `json:"process_number"`
	OS              string         `json:"os"`
	Platform        string         `json:"platform"`
	PlatformVersion string         `json:"platform_version"`
	KernelArch      string         `json:"kernel_arch"`
	KernelVersion   string         `json:"kernel_version"`
	BootTime        string         `json:"boot_time"`
}

type diskDetail struct {
	Device 		string	`json:"device"`
	MountPoint 	string	`json:"mount_point"`
	FsType 		string	`json:"fs_type"`
	UsedPercent float64	`json:"used_percent"`
	FreeSize 	int64	`json:"free_size"`
}

type DiskInfo struct {
	DiskDetail []diskDetail   `json:"disk_detail"`
}

type NetInfo struct {
	EthInterfaces []string       `json:"eth_interfaces"`
}

type PortInfo struct {
	Port    	int   `json:"port"`
	Available   bool  `json:"available"`
}

// ----------------------------------------

func GetCpuInfo() (interface{}, error) {
	cpuInfo := &CpuInfo{CpuDetail: make([]CPUDetail, 0, 1)}

	cpus, err := psCpu.Info()
	if err != nil {
		return nil, err
	}

	var percents []float64
	percents, err = psCpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	if len(cpus) != len(percents) {
		return nil, fmt.Errorf("len(cpu):%d != len(percent):%d error. ", len(cpus), len(percents))
	}

	for i := 0; i < len(cpus); i++ {
		cpuInfo.CpuDetail = append(cpuInfo.CpuDetail, CPUDetail{int(cpus[i].Cores), percents[i]})
	}

	var loads *psLoad.AvgStat
	loads, err = psLoad.Avg()
	if err != nil {
		return nil, err
	}
	cpuInfo.CurrentLoads[0] = loads.Load1
	cpuInfo.CurrentLoads[1] = loads.Load5
	cpuInfo.CurrentLoads[2] = loads.Load15

	return cpuInfo, nil
}


func GetMemInfo() (interface{}, error) {
	memInfo := &MemInfo{}

	mem, err := psMem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memInfo.Total 		= int64(mem.Total)
	memInfo.Available 	= int64(mem.Available)
	memInfo.Used 		= int64(mem.Used)
	memInfo.UsedPercent = mem.UsedPercent

	return memInfo, nil
}

func GetHostInfo() (interface{}, error) {
	hostInfo := &HostInfo{}

	hInfo, err := psHost.Info()
	if err != nil {
		return nil, err
	}
	hostInfo.Hostname 		= hInfo.Hostname
	hostInfo.ProcessNumber 	= hInfo.Procs
	hostInfo.OS 			= hInfo.OS
	hostInfo.Platform 		= hInfo.Platform
	hostInfo.PlatformVersion = hInfo.PlatformVersion
	hostInfo.KernelArch 	= hInfo.KernelArch
	hostInfo.KernelVersion 	= hInfo.KernelVersion
	hostInfo.BootTime 		= time.Unix(int64(hInfo.BootTime), 0).Format("2006-1-2 15:04:05")

	return hostInfo, nil
}

func GetDiskInfo() (interface{}, error) {
	diskInfo := &DiskInfo{DiskDetail: make([]diskDetail, 0, 1)}

	parts, err := psDisk.Partitions(true)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		dt := diskDetail{Device: part.Device, MountPoint: part.Mountpoint, FsType: part.Fstype}

		info, err := psDisk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}
		dt.UsedPercent 	= info.UsedPercent
		dt.FreeSize 	= int64(info.Free)

		diskInfo.DiskDetail = append(diskInfo.DiskDetail, dt)
	}

	return diskInfo, nil
}

func GetNetInfo() (interface{}, error) {
	netInfo := &NetInfo{EthInterfaces: make([]string, 0, 1)}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			continue
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			matched, err := regexp.MatchString(
				"^(?:(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}(?:[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\/([1-9]|[1-2]\\d|3[0-2])$",
				v.String())
			if err == nil && matched {
				netInfo.EthInterfaces = append(netInfo.EthInterfaces, v.String())
			}
		}
	}

	return netInfo, nil
}

func GetPortInfo(host string, port int) (interface{}, error) {
	if port == 0 {
		return getAvailPort(host, port)
	} else {
		return checkPortAvailable(host, port)
	}
}

func getAvailPort(host string, port int) (interface{}, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	_ = l.Close()

	return &PortInfo{Port: l.Addr().(*net.TCPAddr).Port, Available: true}, nil
}

func checkPortAvailable(host string, port int) (interface{}, error) {
	available := true

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
	if err == nil && conn != nil {
		_ = conn.Close()
		available = false
	}

	return &PortInfo{Port: port, Available: available}, nil
}

