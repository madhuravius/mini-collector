package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
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

type wrappedSubsystem struct {
	subsystem subsystem
	optional  bool
}

type collector struct {
	cgroupPath  string
	mountPath   string
	dockerName  string
	statsBuffer cgroups.Stats
	fsBuffer    syscall.Statfs_t
	subsystems  []wrappedSubsystem
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

	subsystems := []wrappedSubsystem{
		wrappedSubsystem{
			subsystem: &fs.CpuGroup{},
			optional:  false,
		},
		wrappedSubsystem{
			subsystem: &fs.MemoryGroup{},
			optional:  false,
		},
		wrappedSubsystem{
			subsystem: &fs.CpuacctGroup{},
			optional:  false,
		},
		wrappedSubsystem{
			subsystem: &fs.BlkioGroup{},
			optional:  false,
		},
		wrappedSubsystem{
			subsystem: &fs.PidsGroup{},
			optional:  true,
		},
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

	for _, wrapper := range c.subsystems {
		cgPath := fmt.Sprintf("%s/%s/docker/%s", c.cgroupPath, wrapper.subsystem.Name(), c.dockerName)

		// Warning - this is done because opencontainers/runc no longer returns an error.
		// This was reported from this issue: https://github.com/opencontainers/runc/issues/1789
		// This is likely a "spooky error", as this must have returned an error before because we have
		// tests that would indicate this via running false with no error. There are no other
		// locations the CgroupPoint{} could have been returned from the codebase but here.
		var cgroupExistsErr error
		if wrapper.subsystem.Name() == "cpu" {
			_, cgroupExistsErr = cgroups.OpenFile(cgPath, "cpu.stat", os.O_RDONLY)
		}

		err := wrapper.subsystem.GetStats(cgPath, &c.statsBuffer)

		if err != nil || cgroupExistsErr != nil {
			if wrapper.optional {
				continue
			}

			if os.IsNotExist(err) || os.IsNotExist(cgroupExistsErr) {
				return CgroupPoint{Running: false}, MakeNoContainerState(pollTime), nil
			}

			errorForDebugPrint := err
			if errorForDebugPrint == nil {
				errorForDebugPrint = cgroupExistsErr
			}
			return CgroupPoint{}, State{}, fmt.Errorf("%s.GetStats failed: %v", wrapper.subsystem.Name(), errorForDebugPrint)
		}
	}

	cpuQuota, cpuPeriod, err := c.getMaxCpu()
	if err != nil {
		return CgroupPoint{}, State{}, fmt.Errorf("Failed to get max CPU: %v", err)
	}

	ioStats := computeIoStats(&c.statsBuffer)

	thisState := State{
		Time:                pollTime,
		AccumulatedCpuUsage: c.statsBuffer.CpuStats.CpuUsage.TotalUsage,
		IoStats:             ioStats,
	}

	milliCpuUsage := computeMilliCpuUsage(thisState, lastState)
	milliCpuLimit := computeMilliCpuLimit(thisState, lastState, cpuQuota, cpuPeriod)

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
		MilliCpuLimit: milliCpuLimit,
		MemoryTotalMb: virtualMemory / MbInBytes,
		MemoryRssMb:   (baseRssMemory + mappedFileMemory) / MbInBytes,
		MemoryLimitMb: (limitMemory) / MbInBytes,
		DiskReadKbps:  readKbps,
		DiskWriteKbps: writeKbps,
		DiskReadIops:  readIops,
		DiskWriteIops: writeIops,
		PidsCurrent:   c.statsBuffer.PidsStats.Current,
		PidsLimit:     c.statsBuffer.PidsStats.Limit,
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

func (c *collector) getMaxCpu() (cpuQuotaUs int64, cpuPeriodUs int64, err error) {
	// The Quota will be negative if no limit is set.
	// The CPU quota and period are available in the same place on the filesystem as the other cgroup
	// info. opencontainers/runc just doesn't make it easy to access through their interface.
	cgPath := fmt.Sprintf("%s/%s/docker/%s", c.cgroupPath, "cpu", c.dockerName)

	f, err := os.Open(filepath.Join(cgPath, "cpu.cfs_quota_us"))
	if err == nil {
		defer f.Close()
		cpuQuotaUs, err = readInt(f)
		if err != nil {
			return
		}
	} else {
		// If the file doesn't exist then assume there is no limit.
		// In this case set the quota to a negative value.
		if !os.IsNotExist(err) {
			return
		}
		cpuQuotaUs = -1
	}

	f, err = os.Open(filepath.Join(cgPath, "cpu.cfs_period_us"))
	if err != nil {
		return
	}
	defer f.Close()
	cpuPeriodUs, err = readInt(f)
	return
}
