package blackhole

import (
	"context"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter"
	log "github.com/sirupsen/logrus"
)

type blackholeEmitter struct{}

func Open() (emitter.Emitter, error) {
	return &blackholeEmitter{}, nil
}

func MustOpen() emitter.Emitter {
	em, err := Open()
	if err != nil {
		panic(err)
	}
	return em
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
