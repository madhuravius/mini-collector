package aggregator

import (
	"context"
	"github.com/aptible/mini-collector/internal/aggregator/batch"
)

type Emitter interface {
	Emit(ctx context.Context, batch batch.Batch) error
	Close()
}
