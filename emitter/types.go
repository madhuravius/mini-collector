package emitter

import (
	"context"
	"github.com/aptible/mini-collector/batch"
)

type Emitter interface {
	Emit(ctx context.Context, batch batch.Batch) error
}
