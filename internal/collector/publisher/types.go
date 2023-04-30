package publisher

import (
	"context"
	"github.com/aptible/mini-collector/internal/collector"
	"github.com/aptible/mini-collector/protobufs"
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

type clientFactory = func(cc grpc.ClientConnInterface) protobufs.AggregatorClient
type connectionFactory = func(ctx context.Context, serverAddress string, dialOption grpc.DialOption) (*grpc.ClientConn, error)
