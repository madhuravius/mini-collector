package collector

import (
	"time"
)

type cgroupPoint struct {
	MilliCpuUsage uint64

	MemoryTotalMb uint64
	MemoryRssMb   uint64
	MemoryLimitMb uint64

	Running bool
}

type diskPoint struct {
	DiskUsageMb uint64
	DiskLimitMb uint64
}

type Point struct {
	cgroupPoint
	diskPoint
}

type State struct {
	Time                time.Time
	AccumulatedCpuUsage uint64
}

type Collector interface {
	GetPoint(lastState State) (Point, State, error)
}
