package collector

import (
	"time"
)

func MakeNoContainerState(time time.Time) State {
	return State{
		Time:                time,
		AccumulatedCpuUsage: MaxUint64,
		IoStats: IoStats{
			ReadBytes:  MaxUint64,
			WriteBytes: MaxUint64,
			ReadOps:    MaxUint64,
			WriteOps:   MaxUint64,
		},
	}
}
