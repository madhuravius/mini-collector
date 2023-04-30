package collector

import (
	"time"
)

type CgroupPoint struct {
	MilliCpuUsage uint64

	MemoryTotalMb uint64
	MemoryRssMb   uint64
	MemoryLimitMb uint64

	DiskReadKbps  uint64
	DiskWriteKbps uint64
	DiskReadIops  uint64
	DiskWriteIops uint64

	PidsCurrent uint64
	PidsLimit   uint64

	Running bool
}

type DiskPoint struct {
	DiskUsageMb int64
	DiskLimitMb int64
}

type Point struct {
	CgroupPoint
	DiskPoint
}

type State struct {
	Time                time.Time
	AccumulatedCpuUsage uint64
	IoStats             IoStats
}

type Collector interface {
	GetPoint(lastState State) (Point, State, error)
}

type IoStats struct {
	ReadBytes  uint64
	WriteBytes uint64
	ReadOps    uint64
	WriteOps   uint64
}
