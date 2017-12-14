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
	clock       clock
}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func NewCollector(cgroupPath string, dockerName string, mountPath string) Collector {
	return newCollector(cgroupPath, dockerName, mountPath)
}

func newCollector(cgroupPath string, dockerName string, mountPath string) *collector {
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
		clock:       &realClock{},
	}
}

func (c *collector) getCgroupPoint(lastState State) (CgroupPoint, State, error) {
	pollTime := c.clock.Now()

	for _, subsys := range c.subsystems {
		cgPath := fmt.Sprintf("%s/%s/docker/%s", c.cgroupPath, subsys.Name(), c.dockerName)

		err := subsys.GetStats(cgPath, &c.statsBuffer)

		if err != nil {
			if os.IsNotExist(err) {
				return CgroupPoint{Running: false}, MakeNoContainerState(pollTime), nil
			}
			return CgroupPoint{}, State{}, fmt.Errorf("%s.GetStats failed: %v", subsys.Name(), err)
		}
	}

	accumulatedCpuUsage := c.statsBuffer.CpuStats.CpuUsage.TotalUsage

	var milliCpuUsage uint64
	if accumulatedCpuUsage > lastState.AccumulatedCpuUsage {
		elapsedCpu := accumulatedCpuUsage - lastState.AccumulatedCpuUsage
		elapsedTime := pollTime.Sub(lastState.Time).Nanoseconds()
		if elapsedTime > 0 {
			milliCpuUsage = uint64(1000 * float64(elapsedCpu) / float64(elapsedTime))
		} else {
			milliCpuUsage = 0
		}
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

	return CgroupPoint{
		MilliCpuUsage: milliCpuUsage,
		MemoryTotalMb: virtualMemory / MbInBytes,
		MemoryRssMb:   (baseRssMemory + mappedFileMemory) / MbInBytes,
		MemoryLimitMb: (limitMemory) / MbInBytes,
		Running:       true,
	}, state, nil

}

func (c *collector) getDiskPoint() (DiskPoint, error) {
	if c.mountPath == "" {
		return DiskPoint{}, nil
	}

	err := syscall.Statfs(c.mountPath, &c.fsBuffer)
	if err != nil {
		return DiskPoint{}, fmt.Errorf("Statfs failed: %v", err)
	}

	blockSize := uint64(c.fsBuffer.Bsize)

	diskUsageBytes := (c.fsBuffer.Blocks - c.fsBuffer.Bavail) * blockSize
	diskLimitBytes := c.fsBuffer.Blocks * blockSize

	return DiskPoint{
		DiskUsageMb: diskUsageBytes / MbInBytes,
		DiskLimitMb: diskLimitBytes / MbInBytes,
	}, nil
}

func (c *collector) GetPoint(lastState State) (Point, State, error) {
	CgroupPoint, thisState, err := c.getCgroupPoint(lastState)
	if err != nil {
		return Point{}, State{}, fmt.Errorf("getCgroupPoint failed: %v", err)
	}

	DiskPoint, err := c.getDiskPoint()
	if err != nil {
		return Point{}, State{}, fmt.Errorf("getDiskPoint failed: %v", err)
	}

	return Point{CgroupPoint, DiskPoint}, thisState, nil
}
