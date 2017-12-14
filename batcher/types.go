package batcher

import (
	"context"
	"github.com/aptible/mini-collector/batch"
)

type Batcher interface {
	Ingest(ctx context.Context, entry *batch.Entry) error
	Close()
}
