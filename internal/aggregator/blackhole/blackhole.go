package blackhole

import (
	"context"
	"github.com/aptible/mini-collector/internal/aggregator"
	"github.com/aptible/mini-collector/internal/aggregator/batch"
	log "github.com/sirupsen/logrus"
)

type blackholeEmitter struct{}

func Open() aggregator.Emitter {
	return &blackholeEmitter{}
}

func (e *blackholeEmitter) Emit(ctx context.Context, batch batch.Batch) error {
	log.WithFields(log.Fields{
		"source":  "emitter",
		"emitter": "blackhole",
	}).WithFields(
		batch.Fields(),
	).Errorf("batch lost")

	return nil
}

func (e *blackholeEmitter) Close() {
}
