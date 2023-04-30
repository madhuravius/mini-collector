package collector

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
)

const (
	kbInBytes = 1024
)

type valueExtractor = func(state State) uint64
type durationExtractor = func(thisState State, lastState State) int64

func extractMicroseconds(thisState State, lastState State) int64 {
	return thisState.Time.Sub(lastState.Time).Nanoseconds() / 1000
}

func extractSeconds(thisState State, lastState State) int64 {
	return int64(thisState.Time.Sub(lastState.Time).Seconds())
}

// NOTE: This is a little less efficient than implementing all the methods by
// hand (in fact, it's twice slower), but considering this only takes a handful
// of nanoseconds, so that's worth it considering we have 5 of these and a risk
// of getting them wrong.
func computeDeltaOverTime(thisState State, lastState State, valueExtractor valueExtractor, durationExtractor durationExtractor) uint64 {
	thisValue := valueExtractor(thisState)
	lastValue := valueExtractor(lastState)

	if thisValue <= lastValue {
		return 0
	}

	valueDelta := thisValue - lastValue
	timeDelta := durationExtractor(thisState, lastState)

	if timeDelta <= 0 {
		return 0
	}

	return uint64(float64(valueDelta) / float64(timeDelta))
}

func computeMilliCpuUsage(thisState State, lastState State) uint64 {
	// NOTE: We use extractMicroseconds but the AccumulatedCpuUsage is in
	// Nanoseconds. This is how we get MilliCpus (what we want).
	return computeDeltaOverTime(
		thisState,
		lastState,
		func(state State) uint64 { return state.AccumulatedCpuUsage },
		extractMicroseconds,
	)
}

func computeReadKbps(thisState State, lastState State) uint64 {
	return computeDeltaOverTime(
		thisState,
		lastState,
		func(state State) uint64 { return state.IoStats.ReadBytes },
		extractSeconds,
	) / kbInBytes
}

func computeWriteKbps(thisState State, lastState State) uint64 {
	return computeDeltaOverTime(
		thisState,
		lastState,
		func(state State) uint64 { return state.IoStats.WriteBytes },
		extractSeconds,
	) / kbInBytes
}

func computeReadIops(thisState State, lastState State) uint64 {
	return computeDeltaOverTime(
		thisState,
		lastState,
		func(state State) uint64 { return state.IoStats.ReadOps },
		extractSeconds,
	)
}

func computeWriteIops(thisState State, lastState State) uint64 {
	return computeDeltaOverTime(
		thisState,
		lastState,
		func(state State) uint64 { return state.IoStats.WriteOps },
		extractSeconds,
	)
}

func computeIoStats(statsBuffer *cgroups.Stats) IoStats {
	ioStats := IoStats{}

	for _, e := range (*statsBuffer).BlkioStats.IoServiceBytesRecursive {
		if e.Op == "Read" {
			ioStats.ReadBytes += e.Value
		}

		if e.Op == "Write" {
			ioStats.WriteBytes += e.Value
		}
	}

	for _, e := range (*statsBuffer).BlkioStats.IoServicedRecursive {
		if e.Op == "Read" {
			ioStats.ReadOps += e.Value
		}

		if e.Op == "Write" {
			ioStats.WriteOps += e.Value
		}
	}

	return ioStats
}
