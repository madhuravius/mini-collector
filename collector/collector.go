package collector

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"os"
	"syscall"
	"time"
)

var (
	noDiskPoint = DiskPoint{
		DiskUsageMb: -1,
		DiskLimitMb: -1,
	}
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
		&fs.BlkioGroup{},
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

	ioStats := computeIoStats(&c.statsBuffer)

	thisState := State{
		Time:                pollTime,
		AccumulatedCpuUsage: c.statsBuffer.CpuStats.CpuUsage.TotalUsage,
		IoStats:             ioStats,
	}

	milliCpuUsage := computeMilliCpuUsage(thisState, lastState)

	readKbps := computeReadKbps(thisState, lastState)
	writeKbps := computeWriteKbps(thisState, lastState)
	readIops := computeReadIops(thisState, lastState)
	writeIops := computeWriteIops(thisState, lastState)

	baseRssMemory := c.statsBuffer.MemoryStats.Stats["rss"]
	mappedFileMemory := c.statsBuffer.MemoryStats.Stats["mapped_file"]
	virtualMemory := c.statsBuffer.MemoryStats.Usage.Usage
	limitMemory := c.statsBuffer.MemoryStats.Usage.Limit

	return CgroupPoint{
		MilliCpuUsage: milliCpuUsage,
		MemoryTotalMb: virtualMemory / MbInBytes,
		MemoryRssMb:   (baseRssMemory + mappedFileMemory) / MbInBytes,
		MemoryLimitMb: (limitMemory) / MbInBytes,
		DiskReadKbps:  readKbps,
		DiskWriteKbps: writeKbps,
		DiskReadIops:  readIops,
		DiskWriteIops: writeIops,
		Running:       true,
	}, thisState, nil

}

func (c *collector) getDiskPoint() (DiskPoint, error) {
	// When no disk is to be scanned, we return negative values for the
	// disk point.
	if c.mountPath == "" {
		return noDiskPoint, nil
	}

	err := syscall.Statfs(c.mountPath, &c.fsBuffer)
	if err != nil {
		return DiskPoint{}, fmt.Errorf("Statfs failed: %v", err)
	}

	blockSize := uint64(c.fsBuffer.Bsize)

	diskUsageBytes := (c.fsBuffer.Blocks - c.fsBuffer.Bavail) * blockSize
	diskLimitBytes := c.fsBuffer.Blocks * blockSize

	return DiskPoint{
		DiskUsageMb: int64(diskUsageBytes / MbInBytes),
		DiskLimitMb: int64(diskLimitBytes / MbInBytes),
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
