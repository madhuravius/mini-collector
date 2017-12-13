package collector

import (
	"time"
)

type CgroupPoint struct {
	MilliCpuUsage uint64

	MemoryTotalMb uint64
	MemoryRssMb   uint64
	MemoryLimitMb uint64

	Running bool
}

type DiskPoint struct {
	DiskUsageMb uint64
	DiskLimitMb uint64
}

type Point struct {
	CgroupPoint
	DiskPoint
}

type State struct {
	Time                time.Time
	AccumulatedCpuUsage uint64
}

type Collector interface {
	GetPoint(lastState State) (Point, State, error)
}
