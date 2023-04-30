package publisher

import (
	"context"
	"fmt"
	"github.com/aptible/mini-collector/internal/collector"
	"github.com/aptible/mini-collector/protobufs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

const (
	serverAddress = "127.0.0.1:1"
)

var testConfig *Config = &Config{
	ServerAddress:  "not-actually-used",
	PublishTimeout: 1 * time.Millisecond,
}

func createMockGrpcConnection(ctx context.Context, serverAddress string, dialOption grpc.DialOption) (*grpc.ClientConn, error) {
	return nil, nil
}

type mockSuccessClient struct {
	calls int
}

func createMockSuccessClient(cc *grpc.ClientConn) protobufs.AggregatorClient {
	return &mockSuccessClient{}
}

func (p *mockSuccessClient) Publish(ctx context.Context, req *protobufs.PublishRequest, opts ...grpc.CallOption) (*protobufs.PublishResponse, error) {
	p.calls++
	return &protobufs.PublishResponse{}, nil
}

type mockErrorClient struct {
	calls int
}

func createMockErrorClient(cc grpc.ClientConnInterface) protobufs.AggregatorClient {
	return &mockErrorClient{}
}

func (p *mockErrorClient) Publish(ctx context.Context, req *protobufs.PublishRequest, opts ...grpc.CallOption) (*protobufs.PublishResponse, error) {
	p.calls++
	return nil, fmt.Errorf("oops")
}

func TestPublisherCloseReturns(t *testing.T) {
	pub := mustOpen(testConfig, createMockGrpcConnection, createMockErrorClient)
	defer pub.Close()
}

func TestPublisherSucceeds(t *testing.T) {
	client := &mockSuccessClient{}
	f := func(cc grpc.ClientConnInterface) protobufs.AggregatorClient { return client }

	pub := mustOpen(testConfig, createMockGrpcConnection, f)
	defer pub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	pub.Queue(ctx, time.Now(), collector.Point{})
	pub.Queue(ctx, time.Now(), collector.Point{})
	<-ctx.Done()

	assert.Equal(t, 2, client.calls)
}

func TestPublisherRetriesAndRateLimits(t *testing.T) {
	client := &mockErrorClient{}
	f := func(cc grpc.ClientConnInterface) protobufs.AggregatorClient { return client }

	pub := mustOpen(testConfig, createMockGrpcConnection, f)
	defer pub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	pub.Queue(ctx, time.Now(), collector.Point{})
	<-ctx.Done()

	assert.True(t, client.calls >= 5)
	assert.True(t, client.calls <= 10)
}
