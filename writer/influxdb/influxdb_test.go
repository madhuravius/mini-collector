package influxdb

import (
	"fmt"
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/batch"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type errorClient struct{}

func (c *errorClient) Write(bp client.BatchPoints) error {
	return fmt.Errorf("permanent failure")
}

func (c *errorClient) Close() error {
	return nil
}

type successClient struct {
	writes chan (client.BatchPoints)
}

func (c *successClient) Write(bp client.BatchPoints) error {
	c.writes <- bp
	return nil
}

func (c *successClient) Close() error {
	return nil
}

func TestItWrites(t *testing.T) {
	c := &successClient{writes: make(chan client.BatchPoints, 1)}

	w := &influxdbWriter{
		database: "foo",
		client:   c,
	}

	t0 := time.Unix(10, 0)

	err := w.Write(batch.Batch{
		Entries: []*batch.Entry{{Time: t0, PublishRequest: &api.PublishRequest{}}},
	})

	if assert.Nil(t, err) {
		select {
		case write := <-c.writes:
			assert.Equal(t, 1, len(write.Points()))
			point := (*write.Points()[0])
			assert.Equal(t, t0, point.Time())
		default:
			t.Fatal("did not receive write")
		}
	}
}

func TestItReturnsAnError(t *testing.T) {
	w := &influxdbWriter{
		database: "foo",
		client:   &errorClient{},
	}

	t0 := time.Unix(10, 0)

	err := w.Write(batch.Batch{
		Entries: []*batch.Entry{{Time: t0, PublishRequest: &api.PublishRequest{}}},
	})

	assert.NotNil(t, err)
}
