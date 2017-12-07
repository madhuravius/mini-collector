package collector

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"syscall"
	"time"
)

type subsystem interface {
	Name() string
	GetStats(path string, stats *cgroups.Stats) error
}

type collector struct {
	cgroupPath  string
	mountPath   string
	dockerName  string
	statsBuffer cgroups.Stats
	fsBuffer    syscall.Statfs_t
	subsystems  []subsystem
}

func NewCollector(cgroupPath string, dockerName string, mountPath string) Collector {
	statsBuffer := *cgroups.NewStats()

	subsystems := []subsystem{
		&fs.CpuGroup{},
		&fs.MemoryGroup{},
		&fs.CpuacctGroup{},
	}

	return &collector{
		cgroupPath:  cgroupPath,
		mountPath:   mountPath,
		dockerName:  dockerName,
		statsBuffer: statsBuffer,
		fsBuffer:    syscall.Statfs_t{},
		subsystems:  subsystems,
	}
}

func (c *collector) GetPoint(lastState State) (Point, State) {
	for _, subsys := range c.subsystems {
		cgPath := fmt.Sprintf("%s/%s/docker/%s", c.cgroupPath, subsys.Name(), c.dockerName)

		err := subsys.GetStats(cgPath, &c.statsBuffer)
		if err != nil {
			// TODO: Logging
			fmt.Printf("%s.GetStats Error: %+v\n", subsys.Name(), err)
			return MakeNoContainerPoint(), MakeNoContainerState()
		}
	}

	pollTime := time.Now()

	accumulatedCpuUsage := c.statsBuffer.CpuStats.CpuUsage.TotalUsage

	var milliCpuUsage uint64
	if accumulatedCpuUsage > lastState.AccumulatedCpuUsage {
		elapsedCpu := float64(accumulatedCpuUsage - lastState.AccumulatedCpuUsage)
		elapsedTime := float64(pollTime.Sub(lastState.Time).Nanoseconds())
		milliCpuUsage = uint64(1000 * elapsedCpu / elapsedTime)
	} else {
		milliCpuUsage = 0
	}

	baseRssMemory := c.statsBuffer.MemoryStats.Stats["rss"]
	mappedFileMemory := c.statsBuffer.MemoryStats.Stats["mapped_file"]
	virtualMemory := c.statsBuffer.MemoryStats.Usage.Usage
	limitMemory := c.statsBuffer.MemoryStats.Usage.Limit

	var (
		diskUsageBytes uint64 = 0
		diskLimitBytes uint64 = 0
	)

	if c.mountPath != "" {
		err := syscall.Statfs(c.mountPath, &c.fsBuffer)
		if err != nil {
			fmt.Printf("Statfs Error: %+v\n", err)
			return MakeNoContainerPoint(), MakeNoContainerState()
		}
		blockSize := uint64(c.fsBuffer.Bsize)
		diskLimitBytes = c.fsBuffer.Blocks * blockSize
		diskUsageBytes = (c.fsBuffer.Blocks - c.fsBuffer.Bavail) * blockSize
	}

	point := Point{
		MilliCpuUsage: milliCpuUsage,
		MemoryTotalMb: virtualMemory / MbInBytes,
		MemoryRssMb:   (baseRssMemory + mappedFileMemory) / MbInBytes,
		MemoryLimitMb: (limitMemory) / MbInBytes,
		DiskUsageMb:   diskUsageBytes / MbInBytes,
		DiskLimitMb:   diskLimitBytes / MbInBytes,
		Running:       true,
	}

	state := State{
		Time:                pollTime,
		AccumulatedCpuUsage: accumulatedCpuUsage,
	}

	return point, state
}
