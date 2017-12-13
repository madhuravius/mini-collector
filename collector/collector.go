package collector

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"os"
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

func (c *collector) getCgroupPoint(lastState State) (cgroupPoint, State, error) {
	for _, subsys := range c.subsystems {
		cgPath := fmt.Sprintf("%s/%s/docker/%s", c.cgroupPath, subsys.Name(), c.dockerName)

		err := subsys.GetStats(cgPath, &c.statsBuffer)

		if err != nil {
			if os.IsNotExist(err) {
				return cgroupPoint{Running: false}, MakeNoContainerState(), nil
			}
			return cgroupPoint{}, State{}, fmt.Errorf("%s.GetStats failed: %v", subsys.Name(), err)
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

	state := State{
		Time:                pollTime,
		AccumulatedCpuUsage: accumulatedCpuUsage,
	}

	return cgroupPoint{
		MilliCpuUsage: milliCpuUsage,
		MemoryTotalMb: virtualMemory / MbInBytes,
		MemoryRssMb:   (baseRssMemory + mappedFileMemory) / MbInBytes,
		MemoryLimitMb: (limitMemory) / MbInBytes,
		Running:       true,
	}, state, nil

}

func (c *collector) getDiskPoint() (diskPoint, error) {
	if c.mountPath == "" {
		return diskPoint{}, nil
	}

	err := syscall.Statfs(c.mountPath, &c.fsBuffer)
	if err != nil {
		return diskPoint{}, fmt.Errorf("Statfs failed: %v", err)
	}

	blockSize := uint64(c.fsBuffer.Bsize)

	diskUsageBytes := (c.fsBuffer.Blocks - c.fsBuffer.Bavail) * blockSize
	diskLimitBytes := c.fsBuffer.Blocks * blockSize

	return diskPoint{
		DiskUsageMb: diskUsageBytes / MbInBytes,
		DiskLimitMb: diskLimitBytes / MbInBytes,
	}, nil
}

func (c *collector) GetPoint(lastState State) (Point, State, error) {
	cgroupPoint, thisState, err := c.getCgroupPoint(lastState)
	if err != nil {
		return Point{}, State{}, fmt.Errorf("getCgroupPoint failed: %v", err)
	}

	diskPoint, err := c.getDiskPoint()
	if err != nil {
		return Point{}, State{}, fmt.Errorf("getDiskPoint failed: %v", err)
	}

	return Point{cgroupPoint, diskPoint}, thisState, nil
}
