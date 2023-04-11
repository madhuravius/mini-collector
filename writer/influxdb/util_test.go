package influxdb

import (
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/batch"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEntryToFieldsIsValid(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &api.PublishRequest{
			MilliCpuUsage: 123,
			MemoryTotalMb: 100,
			MemoryRssMb:   50,
			MemoryLimitMb: 200,
			DiskUsageMb:   -1,
			DiskLimitMb:   -1,
			DiskReadIops:  0,
			Running:       true,
		},
	}

	field := entryToFields(&entry)

	assert.Equal(t, int64(123), field["milli_cpu_usage"])
	assert.Equal(t, int64(100), field["memory_total_mb"])
	assert.Equal(t, int64(50), field["memory_rss_mb"])
	assert.Equal(t, int64(200), field["memory_limit_mb"])
	assert.Equal(t, int64(0), field["disk_read_iops"])

	assert.Equal(t, true, field["running"])

	var ok bool

	_, ok = field["disk_usage_mb"]
	assert.False(t, ok)

	_, ok = field["disk_limit_mb"]
	assert.False(t, ok)
}

func TestBuildBatchPointsIsValid(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := time.Unix(10, 0)

	entries := []*batch.Entry{
		{
			Time:           t0,
			Tags:           map[string]string{"foo": "bar"},
			PublishRequest: &api.PublishRequest{MilliCpuUsage: 123},
		},
		{
			Time:           t1,
			Tags:           map[string]string{"qux": "baz"},
			PublishRequest: &api.PublishRequest{MilliCpuUsage: 456},
		},
	}

	bp := buildBatchPoints("myDb", entries)

	assert.Equal(t, "myDb", bp.Database())

	points := bp.Points()
	if assert.Equal(t, 2, len(points)) {
		p0 := points[0]
		assert.Equal(t, t0, p0.Time())

		p0Fields, err := p0.Fields()
		if assert.Nil(t, err) {
			assert.Equal(t, int64(123), p0Fields["milli_cpu_usage"])
		}
		assert.Equal(t, "bar", p0.Tags()["foo"])
		assert.Equal(t, "", p0.Tags()["qux"])

		p1 := points[1]
		assert.Equal(t, t1, p1.Time())
		p1Fields, err := p1.Fields()
		if assert.Nil(t, err) {
			assert.Equal(t, int64(456), p1Fields["milli_cpu_usage"])
		}
		assert.Equal(t, "baz", p1.Tags()["qux"])
		assert.Equal(t, "", p1.Tags()["foo"])
	}
}
