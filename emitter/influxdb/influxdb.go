package influxdb

import (
	"context"
	"fmt"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter"
	client "github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	bufferSize   = 10
	drainTimeout = 2 * time.Second
)

type influxdbEmitter struct {
	logger *log.Entry

	client      InfluxdbClient
	nextEmitter emitter.Emitter
	database    string

	sendBuffer chan batch.Batch

	doneChannel chan interface{}
	cancel      context.CancelFunc
}

func Open(name string, config *Config, nextEmitter emitter.Emitter) (emitter.Emitter, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("invalid InfluxDb configuration: %v", err)
	}

	return open(name, client, config.Database, nextEmitter), nil
}

func open(name string, client InfluxdbClient, database string, nextEmitter emitter.Emitter) *influxdbEmitter {
	ctx, cancel := context.WithCancel(context.Background())

	em := &influxdbEmitter{
		logger: log.WithFields(log.Fields{
			"source":  "emitter",
			"emitter": name,
		}),

		client:     client,
		database:   database,
		sendBuffer: make(chan batch.Batch, bufferSize),

		nextEmitter: nextEmitter,
		doneChannel: make(chan interface{}),
		cancel:      cancel,
	}

	go func() {
		defer client.Close()
		em.run(ctx)
	}()

	return em
}

func (e *influxdbEmitter) Emit(ctx context.Context, batch batch.Batch) error {
	if len(batch.Entries) <= 0 {
		// Nothing to emit, skip it.
		e.logger.WithFields(
			batch.Fields(),
		).Debugf("skipped")
		return nil
	}

	select {
	case e.sendBuffer <- batch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *influxdbEmitter) run(ctx context.Context) {
	e.logger.Infof("starting")

	for {
		if e.runOnce(ctx) {
			break
		}
	}

	e.logger.Infof("shutting down")

	// TODO: Should prevent further Emits at this stage.
	func() {
		drainCtx, cancel := context.WithTimeout(context.Background(), drainTimeout)
		defer cancel()

		e.drainSendBuffer(drainCtx)
	}()

	e.doneChannel <- nil
}

func (e *influxdbEmitter) runOnce(ctx context.Context) bool {
	select {
	case batch := <-e.sendBuffer:
		e.sendOrDelegateToNextEmitter(ctx, batch)
		return false
	case <-ctx.Done():
		return true
	}
}

func (e *influxdbEmitter) drainSendBuffer(ctx context.Context) error {
	for {
		select {
		case batch := <-e.sendBuffer:
			e.sendOrDelegateToNextEmitter(ctx, batch)
		default:
			return nil
		}
	}
}

func (e *influxdbEmitter) sendOrDelegateToNextEmitter(ctx context.Context, batch batch.Batch) {
	err := e.sendBatch(batch)

	if err != nil {
		e.logger.WithFields(
			batch.Fields(),
		).Warnf("sendBatch failed: %v", err)

		e.nextEmitter.Emit(ctx, batch)
	}
}

func (e *influxdbEmitter) sendBatch(batch batch.Batch) error {
	e.logger.WithFields(
		batch.Fields(),
	).Debugf("sendBatch")

	bp := buildBatchPoints(e.database, batch.Entries)
	err := e.client.Write(bp)

	if err != nil {
		return fmt.Errorf("Write failed: %v", err)
	}

	e.logger.WithFields(
		batch.Fields(),
	).Infof("Write succeeded")

	return nil
}

func (e *influxdbEmitter) Close() {
	e.cancel()
	<-e.doneChannel
}
