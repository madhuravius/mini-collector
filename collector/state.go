package collector

import (
	"time"
)

func MakeNoContainerState(time time.Time) State {
	return State{
		Time:                time,
		AccumulatedCpuUsage: MaxUint64,
	}
}
