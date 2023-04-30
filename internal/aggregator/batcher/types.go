package batcher

import (
	"context"
	"github.com/aptible/mini-collector/internal/aggregator/batch"
)

type Batcher interface {
	Ingest(ctx context.Context, entry *batch.Entry) error
	Close()
}
