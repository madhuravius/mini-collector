package collector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"
)

const (
	testContainerId = "cg"
	testMountPath   = "/"
)

var (
	t0 = time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 = time.Date(2017, 1, 1, 0, 0, 20, 0, time.UTC) // 20 seconds later
)

func getTestDataCgroupPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Join(path.Dir(filename), "testdata")
}

type stubClock struct {
	time time.Time
}

func (sc *stubClock) Now() time.Time {
	return sc.time
}

func TestGetCgroupPointReturnsData(t *testing.T) {
	c := newCollector(getTestDataCgroupPath(), testContainerId, testMountPath)
	c.clock = &stubClock{time: t1}

	point, _, err := c.getCgroupPoint(MakeNoContainerState(t0))

	if assert.Nil(t, err) {
		// CPU should be zero because we don't have any state history
		assert.Equal(t, uint64(0), point.MilliCpuUsage)
		assert.Equal(t, uint64(2), point.MemoryTotalMb)
		assert.Equal(t, uint64(1), point.MemoryRssMb)
		assert.Equal(t, uint64(8), point.MemoryLimitMb)
		assert.Equal(t, true, point.Running)
	}
}

func TestGetCgroupPointReturnsAccumulatedCpuUsage(t *testing.T) {
	c := newCollector(getTestDataCgroupPath(), testContainerId, testMountPath)
	c.clock = &stubClock{time: t1}

	lastState := State{
		Time:                t0,
		AccumulatedCpuUsage: 10000000000,
	}

	point, thisState, err := c.getCgroupPoint(lastState)

	if assert.Nil(t, err) {
		// We had about 10.04 seconds of CPU usage over 20 seconds of
		// runtime. This gives us 0.502 CPUs, i.e. 502 milli-cpus
		assert.Equal(t, t1, thisState.Time)
		assert.Equal(t, uint64(20042190861), thisState.AccumulatedCpuUsage)
		assert.Equal(t, uint64(502), point.MilliCpuUsage)
	}
}

func TestGetCgroupPointReturnsAccumulatedIoStats(t *testing.T) {
	c := newCollector(getTestDataCgroupPath(), testContainerId, testMountPath)
	c.clock = &stubClock{time: t1}

	lastState := State{
		Time: t0,
		IoStats: IoStats{
			ReadBytes:  16097280,
			WriteBytes: 24694784,
			ReadOps:    81,
			WriteOps:   20,
		},
	}

	point, thisState, err := c.getCgroupPoint(lastState)

	if assert.Nil(t, err) {
		assert.Equal(t, t1, thisState.Time)

		// 16097280 bytes in 20 seconds = 786 kbps
		assert.Equal(t, uint64(32194560), thisState.IoStats.ReadBytes)
		assert.Equal(t, uint64(786), point.DiskReadKbps)

		// 24694784 bytes in 20 seconds = 1205 kbps
		assert.Equal(t, uint64(49389568), thisState.IoStats.WriteBytes)
		assert.Equal(t, uint64(1205), point.DiskWriteKbps)

		// 600 ops in 20 seconds = 30 iops
		assert.Equal(t, uint64(681), thisState.IoStats.ReadOps)
		assert.Equal(t, uint64(30), point.DiskReadIops)

		// 60 ops in 20 seconds = 3 iops
		assert.Equal(t, uint64(80), thisState.IoStats.WriteOps)
		assert.Equal(t, uint64(3), point.DiskWriteIops)
	}
}

func TestGetCgroupPointReturnsZeroUsageForZeroTime(t *testing.T) {
	c := newCollector(getTestDataCgroupPath(), testContainerId, testMountPath)
	c.clock = &stubClock{time: t0}

	point, _, err := c.getCgroupPoint(State{Time: t0})
	if assert.Nil(t, err) {
		assert.Equal(t, uint64(0), point.MilliCpuUsage)
	}
}

func TestGetCgroupPointReturnsNotRunningForNoCgroup(t *testing.T) {
	c := newCollector(
		getTestDataCgroupPath(),
		fmt.Sprintf("%sfoobar", testContainerId),
		testMountPath,
	)
	c.clock = &stubClock{time: t1}

	point, state, err := c.getCgroupPoint(State{Time: t0})
	if assert.Nil(t, err) {
		assert.Equal(t, t1, state.Time)
		assert.Equal(t, MaxUint64, state.AccumulatedCpuUsage)
		assert.Equal(t, false, point.Running)
	}
}

func TestGetCgroupPointReturnsErrorForOtherError(t *testing.T) {
	// Set up a copy in a temp directory so we can play with the
	// permissions. This is the easiest way for us to test a "not
	// not-found" error.
	dir, tempDirErr := ioutil.TempDir("", "work")
	if tempDirErr != nil {
		t.Fatalf("TempDir failed: %v", tempDirErr)
	}
	defer func() {
		exec.Command("rm", "-r", dir).Run()
	}()

	cpCmd := exec.Command("cp", "-r", getTestDataCgroupPath(), dir)
	cpErr := cpCmd.Run()
	if cpErr != nil {
		t.Fatalf("cp failed: %v", cpErr)
	}

	cgPath := path.Join(dir, "testdata")

	memoryStat := path.Join(
		cgPath,
		"memory",
		"docker",
		testContainerId,
		"memory.stat",
	)

	chmodCmd := exec.Command("chmod", "000", memoryStat)
	chmodErr := chmodCmd.Run()
	if chmodErr != nil {
		t.Fatalf("chmod failed: %v", chmodErr)
	}

	c := newCollector(
		cgPath,
		testContainerId,
		testMountPath,
	)
	c.clock = &stubClock{time: t1}

	_, _, err := c.getCgroupPoint(State{Time: t0})
	assert.NotNil(t, err)
}

func TestGetDiskPointReturnsData(t *testing.T) {
	c := newCollector(getTestDataCgroupPath(), testContainerId, testMountPath)
	point, err := c.getDiskPoint()
	if assert.Nil(t, err) {
		assert.True(t, point.DiskLimitMb > point.DiskUsageMb)
	}
}

func TestGetDiskPointReturnsErrorForPathNotExistent(t *testing.T) {
	c := newCollector(
		getTestDataCgroupPath(),
		testContainerId,
		fmt.Sprintf("%s/does-not-exist", getTestDataCgroupPath()),
	)
	_, err := c.getDiskPoint()
	assert.NotNil(t, err)
}
