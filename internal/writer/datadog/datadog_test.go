package datadog

import (
	"github.com/aptible/mini-collector/internal/aggregator/batch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mapFromPayload(p datadogPayload) map[string]int64 {
	out := map[string]int64{}

	for _, series := range p.Series {
		metric := series.Metric
		val := series.Points[0][1]
		out[metric] = val.(int64)
	}

	return out
}

func TestFormatBatchAddsHostIfPresent(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &protobufs.PublishRequest{},
		Tags: map[string]string{
			"host_name": "my-host",
		},
	}

	payload := formatBatch(batch.Batch{
		Entries: []*batch.Entry{&entry},
	})

	assert.True(t, len(payload.Series) > 0)

	for _, series := range payload.Series {
		assert.Equal(t, "my-host", series.Host)
	}
}

func TestFormatBatchSkipsHostIfNotPresent(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &protobufs.PublishRequest{},
		Tags:           map[string]string{},
	}

	payload := formatBatch(batch.Batch{
		Entries: []*batch.Entry{&entry},
	})

	assert.True(t, len(payload.Series) > 0)

	for _, series := range payload.Series {
		assert.Equal(t, "", series.Host)
	}
}

func TestFormatBatchIsValid(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &protobufs.PublishRequest{
			MilliCpuUsage: 123,
			MemoryTotalMb: 100,
			MemoryRssMb:   50,
			MemoryLimitMb: 200,
			DiskUsageMb:   -1,
			DiskLimitMb:   -1,
			DiskReadIops:  0,
		},
		Tags: map[string]string{},
	}

	m := mapFromPayload(formatBatch(batch.Batch{
		Entries: []*batch.Entry{&entry},
	}))

	var (
		v  int64
		ok bool
	)

	v, ok = m["enclave.milli_cpu_usage"]
	assert.True(t, ok)
	assert.Equal(t, int64(123), v)

	v, ok = m["enclave.disk_usage_mb"]
	assert.False(t, ok)
}

func TestFormatBatchRunning(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &protobufs.PublishRequest{
			Running: true,
		},
		Tags: map[string]string{},
	}

	m := mapFromPayload(formatBatch(batch.Batch{
		Entries: []*batch.Entry{&entry},
	}))

	v, ok := m["enclave.running"]
	if assert.True(t, ok) {
		assert.Equal(t, int64(1), v)
	}
}

func TestFormatBatchNotRunning(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: &protobufs.PublishRequest{
			Running: false,
		},
		Tags: map[string]string{},
	}

	m := mapFromPayload(formatBatch(batch.Batch{
		Entries: []*batch.Entry{&entry},
	}))

	v, ok := m["enclave.running"]
	if assert.True(t, ok) {
		assert.Equal(t, int64(0), v)
	}
}
