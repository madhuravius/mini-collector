package publisher

import (
	"context"
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/collector"
	"google.golang.org/grpc"
	"time"
)

type Publisher interface {
	Queue(ctx context.Context, ts time.Time, point collector.Point) error
	Close()
}

type Config struct {
	ServerAddress string
	DialOption    grpc.DialOption
	Tags          map[string]string

	BufferSize     int
	PublishTimeout time.Duration
}

type clientFactory = func(cc grpc.ClientConnInterface) api.AggregatorClient
type connectionFactory = func(ctx context.Context, serverAddress string, dialOption grpc.DialOption) (*grpc.ClientConn, error)
