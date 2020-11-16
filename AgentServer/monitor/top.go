package monitor

import (
	"fmt"
	psCpu "github.com/shirou/gopsutil/cpu"
	psDisk "github.com/shirou/gopsutil/disk"
	psHost "github.com/shirou/gopsutil/host"
	psLoad "github.com/shirou/gopsutil/load"
	psMem "github.com/shirou/gopsutil/mem"
	"net"
	"regexp"
	"strings"
	"tafagent/common"
	"time"
)

type CPUDetail struct {
	Cores 	int		`json:"cores"`
	Percent float64	`json:"percent"`
}

type CpuInfo struct {
	ErrInfo      common.ErrorInfo   `json:"err_info"`
	CpuDetail    []CPUDetail `json:"cpu_detail"`
	CurrentLoads [3]float64  `json:"current_loads"`
}

type MemInfo struct {
	ErrInfo     common.ErrorInfo `json:"err_info"`
	Total       int64     `json:"total"`
	Available   int64     `json:"available"`
	Used        int64     `json:"used"`
	UsedPercent float64   `json:"used_percent"`
}

type HostInfo struct {
	ErrInfo         common.ErrorInfo `json:"err_info"`
	Hostname        string    `json:"hostname"`
	ProcessNumber   uint64    `json:"process_number"`
	OS              string    `json:"os"`
	Platform        string    `json:"platform"`
	PlatformVersion string    `json:"platform_version"`
	KernelArch      string    `json:"kernel_arch"`
	KernelVersion   string    `json:"kernel_version"`
	BootTime        string    `json:"boot_time"`
}

type diskDetail struct {
	Device 		string	`json:"device"`
	MountPoint 	string	`json:"mount_point"`
	FsType 		string	`json:"fs_type"`
	UsedPercent float64	`json:"used_percent"`
	FreeSize 	int64	`json:"free_size"`
}

type DiskInfo struct {
	ErrInfo    common.ErrorInfo    `json:"err_info"`
	DiskDetail []diskDetail `json:"disk_detail"`
}

type NetInfo struct {
	ErrInfo       common.ErrorInfo `json:"err_info"`
	EthInterfaces []string  `json:"eth_interfaces"`
	AvailPorts    []int     `json:"avail_ports"`
}

// ----------------------------------------

type TopData struct {
	Cpu 	*CpuInfo
	Mem 	*MemInfo
	Host 	*HostInfo
	Disk 	*DiskInfo
	Net 	*NetInfo
}

func (td * TopData) NewCpuInfo() *CpuInfo {
	td.Cpu = &CpuInfo{CpuDetail: make([]CPUDetail, 0, 1)}
	return td.Cpu
}

func (td * TopData) NewMemInfo() *MemInfo {
	td.Mem = &MemInfo{}
	return td.Mem
}

func (td * TopData) NewHostInfo() *HostInfo {
	td.Host = &HostInfo{}
	return td.Host
}

func (td * TopData) NewDiskInfo() *DiskInfo {
	td.Disk = &DiskInfo{DiskDetail: make([]diskDetail, 0, 1)}
	return td.Disk
}

func (td * TopData) NewNetInfo() *NetInfo {
	td.Net = &NetInfo{EthInterfaces: make([]string, 0, 1), AvailPorts: make([]int, 0, 1)}
	return td.Net
}


func (td * TopData) getCpuInfo() *CpuInfo {
	cpuInfo := td.NewCpuInfo()

	cpus, err := psCpu.Info()
	if err != nil {
		cpuInfo.ErrInfo.ErrCode=-1
		cpuInfo.ErrInfo.ErrMsg = err.Error()
		return cpuInfo
	}

	var percents []float64
	percents, err = psCpu.Percent(time.Second, false)
	if err != nil {
		cpuInfo.ErrInfo.ErrCode=-1
		cpuInfo.ErrInfo.ErrMsg = err.Error()
		return cpuInfo
	}
	if len(cpus) != len(percents) {
		cpuInfo.ErrInfo.ErrCode=-1
		cpuInfo.ErrInfo.ErrMsg = fmt.Errorf("len(cpu):%d != len(percent):%d error.", len(cpus), len(percents)).Error()
		return cpuInfo
	}

	for i := 0; i < len(cpus); i++ {
		cpuInfo.CpuDetail = append(cpuInfo.CpuDetail, CPUDetail{int(cpus[i].Cores), percents[i]})
	}

	var loads *psLoad.AvgStat
	loads, err = psLoad.Avg()
	if err != nil {
		cpuInfo.ErrInfo.ErrCode=-1
		cpuInfo.ErrInfo.ErrMsg = err.Error()
		return cpuInfo
	}
	cpuInfo.CurrentLoads[0] = loads.Load1
	cpuInfo.CurrentLoads[1] = loads.Load5
	cpuInfo.CurrentLoads[2] = loads.Load15

	return cpuInfo
}


func (td * TopData) getMemInfo() *MemInfo {
	memInfo := td.NewMemInfo()

	mem, err := psMem.VirtualMemory()
	if err != nil {
		memInfo.ErrInfo.ErrCode = -1;
		memInfo.ErrInfo.ErrMsg 	= err.Error()
		return memInfo
	}
	memInfo.Total 		= int64(mem.Total)
	memInfo.Available 	= int64(mem.Available)
	memInfo.Used 		= int64(mem.Used)
	memInfo.UsedPercent = mem.UsedPercent

	return memInfo
}

func (td * TopData) getHostInfo() *HostInfo {
	hostInfo := td.NewHostInfo()

	hInfo, err := psHost.Info()
	if err != nil {
		hostInfo.ErrInfo.ErrCode = -1
		hostInfo.ErrInfo.ErrMsg 	= err.Error()
		return hostInfo
	}
	hostInfo.Hostname 		= hInfo.Hostname
	hostInfo.ProcessNumber 	= hInfo.Procs
	hostInfo.OS 			= hInfo.OS
	hostInfo.Platform 		= hInfo.Platform
	hostInfo.PlatformVersion = hInfo.PlatformVersion
	hostInfo.KernelArch 	= hInfo.KernelArch
	hostInfo.KernelVersion 	= hInfo.KernelVersion
	hostInfo.BootTime 		= time.Unix(int64(hInfo.BootTime), 0).Format("2006-1-2 15:04:05")

	return hostInfo
}

func (td * TopData) getDiskInfo() *DiskInfo {
	diskInfo := td.NewDiskInfo()

	parts, err := psDisk.Partitions(true)
	if err != nil {
		diskInfo.ErrInfo.ErrCode 	= -1
		diskInfo.ErrInfo.ErrMsg 	= err.Error()
		return diskInfo
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

	return diskInfo
}

func (td * TopData) getNetInfo() *NetInfo {
	netInfo := td.NewNetInfo()

	interfaces, err := net.Interfaces()
	if err != nil {
		netInfo.ErrInfo.ErrCode 	= -1
		netInfo.ErrInfo.ErrMsg 	= err.Error()
		return netInfo
	}

	availEth := make([]string, 0, 1)

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
				availEth = append(availEth, v.String())
			}
		}
	}

	pd := PortData{}
	for _, v := range availEth {
		port := 0
		if len(netInfo.AvailPorts) > 0 {
			port = netInfo.AvailPorts[0]
		}

		ip := (strings.Split(v, "/"))[0]
		ans := pd.GetAvailPort(ip, port)
		if ans.ErrInfo.ErrCode == -1 {
			if port == 0 {
				continue
			}
			ans= pd.GetAvailPort(ip, 0)
			if ans.ErrInfo.ErrCode == -1 {
				continue
			}
		}
		netInfo.EthInterfaces = append(netInfo.EthInterfaces, v)
		netInfo.AvailPorts = append(netInfo.AvailPorts, ans.Port)
	}

	return netInfo
}


