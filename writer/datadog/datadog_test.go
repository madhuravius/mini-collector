package datadog

import (
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/batch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatBatchAddsHostIfPresent(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: api.PublishRequest{},
		Tags: map[string]string{
			"host": "my-host",
		},
	}

	payload := formatBatch(batch.Batch{
		Entries: []batch.Entry{entry},
	})

	assert.True(t, len(payload.Series) > 0)

	for _, series := range payload.Series {
		assert.Equal(t, "my-host", series.Host)
	}
}

func TestFormatBatchSkipsHostIfNotPresent(t *testing.T) {
	entry := batch.Entry{
		PublishRequest: api.PublishRequest{},
		Tags:           map[string]string{},
	}

	payload := formatBatch(batch.Batch{
		Entries: []batch.Entry{entry},
	})

	assert.True(t, len(payload.Series) > 0)

	for _, series := range payload.Series {
		assert.Equal(t, "", series.Host)
	}
}
