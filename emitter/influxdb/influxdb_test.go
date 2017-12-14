package influxdb

import (
	"context"
	"fmt"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter/blackhole"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	writeTimeout = 5 * time.Millisecond
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

func TestEmitEmptyBatchDoesNotDoAnything(t *testing.T) {
	c := &successClient{writes: make(chan client.BatchPoints)}
	em := open("em", c, "myDb", blackhole.MustOpen())
	defer em.Close()

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()

	err := em.Emit(ctx, batch.Batch{})
	if assert.Nil(t, err) {
		select {
		case <-c.writes:
			t.Fatalf("received write")
		case <-ctx.Done():
			// what we expected here
		}
	}
}

func TestEmitSendsToInfluxdb(t *testing.T) {
	c := &successClient{writes: make(chan client.BatchPoints)}
	em := open("em", c, "myDb", blackhole.MustOpen())
	defer em.Close()

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()

	t0 := time.Unix(10, 0)

	err := em.Emit(ctx, batch.Batch{
		Entries: []batch.Entry{batch.Entry{Time: t0}},
	})

	if assert.Nil(t, err) {
		select {
		case write := <-c.writes:
			assert.Equal(t, 1, len(write.Points()))
			point := (*write.Points()[0])
			assert.Equal(t, t0, point.Time())
		case <-ctx.Done():
			t.Fatalf("did not receive write")
		}
	}
}

func TestEmitterPassesToNextEmitter(t *testing.T) {
	n := 2

	c0 := &successClient{writes: make(chan client.BatchPoints)}
	em0 := open("em0", c0, "mydb", blackhole.MustOpen())
	defer em0.Close()

	c1 := &errorClient{}
	em1 := open("em0", c1, "mydb", em0)
	defer em1.Close()

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()

	for i := 0; i < n; i++ {
		err := em1.Emit(ctx, batch.Batch{
			Entries: []batch.Entry{batch.Entry{}},
		})
		assert.Nil(t, err)
	}

	for i := 0; i < n; i++ {
		// Expect to have received n+1 writes.
		select {
		case <-c0.writes:
			// What we want
		case <-ctx.Done():
			t.Fatalf("did not receive write %d", i)
		}
	}
}
