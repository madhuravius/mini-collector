package publisher

import (
	"context"
	"github.com/aptible/mini-collector/collector"
	"time"
)

type Publisher interface {
	Queue(ctx context.Context, ts time.Time, point collector.Point) error
	Close()
}
