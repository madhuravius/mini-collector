package publisher

import (
	"fmt"
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/collector"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	maxBackoffDuration    = 5 * time.Second
	defaultPublishTimeout = 4 * time.Second
	defaultBufferSize     = 10
)

type publisher struct {
	publishTimeout time.Duration

	serverAddress string
	dialOption    grpc.DialOption
	tags          map[string]string

	publishChannel chan *api.PublishRequest
	doneChannel    chan interface{}
	cancel         context.CancelFunc

	clientFactory clientFactory
}

func createGrpcConnection(ctx context.Context, serverAddress string, dialOption grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		serverAddress,
		dialOption,
		grpc.WithBackoffMaxDelay(maxBackoffDuration),
	)
}

func Open(config *Config) (Publisher, error) {
	return open(config, createGrpcConnection, api.NewAggregatorClient)
}

func mustOpen(config *Config, cnf connectionFactory, clf clientFactory) Publisher {
	pub, err := open(config, cnf, clf)
	if err != nil {
		panic(err)
	}
	return pub
}

func open(config *Config, cnf connectionFactory, clf clientFactory) (Publisher, error) {
	ctx, cancel := context.WithCancel(context.Background())

	publishTimeout := (*config).PublishTimeout
	if publishTimeout == 0 {
		publishTimeout = defaultPublishTimeout
	}

	bufferSize := (*config).BufferSize
	if bufferSize == 0 {
		bufferSize = defaultBufferSize
	}

	serverAddress := (*config).ServerAddress
	if serverAddress == "" {
		return nil, fmt.Errorf("ServerAddress is required")
	}

	p := &publisher{
		publishTimeout: publishTimeout,

		serverAddress: (*config).ServerAddress,
		dialOption:    (*config).DialOption,
		tags:          (*config).Tags,

		publishChannel: make(chan *api.PublishRequest, bufferSize),
		doneChannel:    make(chan interface{}),
		cancel:         cancel,

		clientFactory: clf,
	}

	conn, err := cnf(ctx, p.serverAddress, p.dialOption)
	if err != nil {
		return nil, fmt.Errorf("connectionFactory failed: %v", err)
	}

	go func() {
		// This is here to support tests, where we return nil for the
		// *grpc.ClientConn as a stub: since *grpc.ClientConn is a
		// struct (as opposed to an interface), we can't do any better.
		if conn != nil {
			defer conn.Close()
		}

		p.startConnection(ctx, conn)
	}()

	return p, nil
}

func (p *publisher) startConnection(ctx context.Context, connection *grpc.ClientConn) {
	client := p.clientFactory(connection)

	md := metadata.New(p.tags)

	baseCtx := metadata.NewOutgoingContext(ctx, md)

PublishLoop:
	for {
		select {
		case payload := <-p.publishChannel:
			err := func() error {
				localCtx, cancel := context.WithTimeout(baseCtx, p.publishTimeout)
				defer cancel()

				_, err := client.Publish(localCtx, payload, grpc.FailFast(false))
				if err != nil {
					// Wait on the context no matter what.
					// This ensures that even if grpc
					// returns quicly despite
					// grpc.FailFast(false) being set, we
					// don't accidentally go into a hot
					// loop.
					<-localCtx.Done()
					return err
				}

				return nil
			}()

			if err != nil {
				// Try to requeue the request. But, if the buffer is
				// full, just drop it (favor more recent data points).
				select {
				case p.publishChannel <- payload:
					log.Infof("requeued point [%v]: %v", (*payload).UnixTime, err)
				default:
					log.Warnf("dropped point [%v]: %v", (*payload).UnixTime, err)
				}

				continue
			}

			log.Debugf("delivered point [%v]", (*payload).UnixTime)
		case <-ctx.Done():
			log.Debugf("shutdown loop loop")
			break PublishLoop
		}
	}

	log.Debugf("shutdown publisher")
	p.doneChannel <- nil
}

func (p *publisher) Queue(ctx context.Context, ts time.Time, point collector.Point) error {
	payload := buildPublishRequest(ts, point)

	select {
	case p.publishChannel <- &payload:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *publisher) Close() {
	p.cancel()
	<-p.doneChannel
}
