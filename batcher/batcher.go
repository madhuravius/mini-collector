package batcher

import (
	"context"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ingestBufferSize = 10
	emitTimeout      = time.Second
)

type batcher struct {
	emitter emitter.Emitter

	minPublishFrequency time.Duration
	maxBatchSize        int

	ingestBuffer chan *batch.Entry
}

func New(emitter emitter.Emitter, minPublishFrequency time.Duration, maxBatchSize int) Batcher {
	b := &batcher{
		emitter:             emitter,
		minPublishFrequency: minPublishFrequency,
		maxBatchSize:        maxBatchSize,
		ingestBuffer:        make(chan *batch.Entry, ingestBufferSize),
	}

	go b.start()

	return b
}

func (b *batcher) start() {
	for {
		lastPublish := time.Now()
		currentBatch := make([]batch.Entry, 0, b.maxBatchSize)

	BatchLoop:
		for {
			select {
			case newEntry := <-b.ingestBuffer:
				currentBatch = append(currentBatch, *newEntry)
				if len(currentBatch) >= b.maxBatchSize {
					break BatchLoop
				}
			case <-time.After(time.Until(lastPublish.Add(b.minPublishFrequency))):
				break BatchLoop
			}
		}

		b.emitBatch(currentBatch)
	}
}

func (b *batcher) emitBatch(batch []batch.Entry) {
	ctx, cancel := context.WithTimeout(context.Background(), emitTimeout)
	defer cancel()

	err := b.emitter.Emit(ctx, batch)
	if err != nil {
		log.Errorf("emitter did not accept batch (%d entries): %v", len(batch), err)
	}
}

func (b *batcher) Ingest(ctx context.Context, entry *batch.Entry) error {
	select {
	case b.ingestBuffer <- entry:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
