package hold

import (
	"context"
	"sync"
	"fmt"
	"time"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter"
	"github.com/sirupsen/logrus"
)

const (
	defaultDelegateTimeout = 2 * time.Second
)

type holdEmitter struct {
	logger *logrus.Entry

	nextEmitter emitter.Emitter
	delay time.Duration

	delegateTimeout time.Duration

	context context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
}

func Open(delay time.Duration, nextEmitter emitter.Emitter) (emitter.Emitter, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &holdEmitter{
		logger: logrus.WithFields(logrus.Fields{
			"source": "emitter",
			"emitter": fmt.Sprintf("hold for %s", delay),
		}),

		nextEmitter: nextEmitter,
		delay: delay,

		delegateTimeout: defaultDelegateTimeout,

		context: ctx,
		cancel: cancel,
		wg: sync.WaitGroup{},
	}, nil
}

func MustOpen(delay time.Duration, nextEmitter emitter.Emitter) emitter.Emitter {
	em, err := Open(delay, nextEmitter)
	if err != nil {
		panic(err)
	}
	return em
}

func (em *holdEmitter) Emit(ctx context.Context, batch batch.Batch) error {
	em.wg.Add(1)

	// We ignore the incoming context: we're going to kick off a new
	// goroutine immediatley and accept the message (technically we should
	// wait on the lock for that)
	go func() {
		defer em.wg.Done()
		em.holdThenDelegateToNextEmitter(batch)
	}()

	return nil
}

func (em *holdEmitter) Close() {
	// TODO: This will race with Emit. More broadly speaking, you can call
	// Close() *then* Emit().
	em.cancel()
	em.wg.Wait()
}

func (em *holdEmitter) holdThenDelegateToNextEmitter(batch batch.Batch) {
	select {
	case <- em.context.Done():
		// We're closing, stop holding!
	case <- time.After(em.delay):
		// We're done waiting for this one.:w
	}

	ctx, cancel := context.WithTimeout(context.Background(), em.delegateTimeout)
	defer cancel()

	err := em.nextEmitter.Emit(ctx, batch)

	if err != nil {
		em.logger.WithFields(
			batch.Fields(),
		).Errorf("nextEmitter.Emit failed: %v", err)
	}
}
