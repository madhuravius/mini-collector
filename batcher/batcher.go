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
	doneChannel  chan interface{}
	cancel       context.CancelFunc
}

func New(emitter emitter.Emitter, minPublishFrequency time.Duration, maxBatchSize int) Batcher {
	ctx, cancel := context.WithCancel(context.Background())

	b := &batcher{
		emitter:             emitter,
		minPublishFrequency: minPublishFrequency,
		maxBatchSize:        maxBatchSize,
		ingestBuffer:        make(chan *batch.Entry, ingestBufferSize),
		doneChannel:         make(chan interface{}),
		cancel:              cancel,
	}

	go func() {
		b.run(ctx)
		b.doneChannel <- nil
	}()

	return b
}

func (b *batcher) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// no-op: proceed
		}

		nextBatch := func() batch.Batch {
			batchCtx, cancel := context.WithTimeout(ctx, b.minPublishFrequency)
			defer cancel()
			return b.prepareBatch(batchCtx)
		}()

		b.emitBatch(nextBatch)
	}
}

func (b *batcher) prepareBatch(ctx context.Context) batch.Batch {
	currentBatch := make(batch.Batch, 0, b.maxBatchSize)

	for {
		select {
		case newEntry := <-b.ingestBuffer:
			currentBatch = append(currentBatch, *newEntry)
			if len(currentBatch) >= b.maxBatchSize {
				return currentBatch
			}
		case <-ctx.Done():
			return currentBatch
		}
	}
}

func (b *batcher) emitBatch(batch batch.Batch) {
	log.Infof("emitting batch (%d entries)", len(batch))

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

func (b *batcher) Close() {
	b.cancel()
	<-b.doneChannel
}
