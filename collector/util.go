package collector

import (
	"io"
	"os"
	"strconv"
	"strings"

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

func extractMilliseconds(thisState State, lastState State) int64 {
	return thisState.Time.Sub(lastState.Time).Nanoseconds() / 1000000
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

func computeMilliCpuLimit(thisState State, lastState State, cpuQuota int64, cpuPeriod int64) uint64 {
	ms := extractMilliseconds(thisState, lastState)

	if ms <= 0 {
		return 0
	}
	if cpuQuota <= 0 {
		cpuQuota = 0
	}

	return uint64(cpuQuota * ms / cpuPeriod)
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

func readInt(f *os.File) (int64, error) {
	data, err := readAll(f)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}

	return int64(result), nil
}

// Unfortunately io.ReadAll is not available in go 1.9, so we've copied it:
// https://cs.opensource.google/go/go/+/refs/tags/go1.19.5:src/io/io.go;l=650-670

// readAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because readAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func readAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
