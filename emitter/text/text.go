package text

import (
	"context"
	"fmt"
	"github.com/aptible/mini-collector/batch"
)

type textEmitter struct{}

func New() (*textEmitter, error) {
	return &textEmitter{}, nil
}

func (t *textEmitter) Emit(ctx context.Context, batch batch.Batch) error {
	for _, entry := range batch.Entries {
		fmt.Printf("%+v\n", entry)
	}

	return nil
}

func (t *textEmitter) Close() {
}
