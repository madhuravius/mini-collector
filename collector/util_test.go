package collector

import (
	"testing"
)

func BenchmarkComputeMilliCpuUsage(b *testing.B) {
	lastState := State{
		Time:                t0,
		AccumulatedCpuUsage: 10000000000,
	}

	thisState := State{
		Time:                t1,
		AccumulatedCpuUsage: 20000000000,
	}

	for n := 0; n < b.N; n++ {
		computeMilliCpuUsage(thisState, lastState)
	}
}
