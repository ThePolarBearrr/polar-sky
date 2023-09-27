package worker

import (
	"github.com/c9s/goprocinfo/linux"
	"polar-sky/log"
	"time"
)

type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.MemAvailable
}

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.MemTotal
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}

func (s *Stats) MemUsedPercent() uint64 {
	return s.MemUsedKb() / s.MemStats.MemTotal
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.All
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

func (s *Stats) CpuUsage() float64 {
	idle := s.CpuStats.Idle + s.CpuStats.IOWait
	nonIdle := s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.System + s.CpuStats.IRQ + s.CpuStats.SoftIRQ + s.CpuStats.Steal
	total := idle + nonIdle
	if total == 0 {
		return 0.0
	}
	return (float64(total) - float64(idle)) / float64(total)
}

func GetMemInfo() *linux.MemInfo {
	memInfo, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Logger.Errorf("Failed get meminfo from /proc/meminfo, error: %v\n", err)
		return &linux.MemInfo{}
	}
	return memInfo
}

func GetDiskInfo() *linux.Disk {
	diskStats, err := linux.ReadDisk("/")
	if err != nil {
		log.Logger.Errorf("Failed get disk stats from /")
		return &linux.Disk{}
	}
	return diskStats
}

func GetCpuStats() *linux.CPUStat {
	cpuStats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Logger.Errorf("Failed get stat from /proc/stat, error: %v\n", err)
		return &linux.CPUStat{}
	}
	return &cpuStats.CPUStatAll
}

func GetLoadAvg() *linux.LoadAvg {
	loadAvg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Logger.Errorf("Failed get load avg from /proc/loadavg, error: %v\n", err)
		return &linux.LoadAvg{}
	}
	return loadAvg
}

func GetStats() *Stats {
	return &Stats{
		MemStats:  GetMemInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats:  GetCpuStats(),
		LoadStats: GetLoadAvg(),
	}
}

func (w *Worker) CollectStats() {
	for {
		w.Logger.Infof("Collect stats")
		w.Stats = GetStats()
		time.Sleep(time.Second * 15)
	}
}
