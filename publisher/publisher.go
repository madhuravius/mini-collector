package publisher

import (
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/collector"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type publisher struct {
	connectTimeout time.Duration
	publishTimeout time.Duration

	serverAddress string
	dialOption    grpc.DialOption
	tags          map[string]string

	publishChannel chan *api.PublishRequest
	doneChannel    chan interface{}
	cancel         context.CancelFunc
}

func Open(serverAddress string, dialOption grpc.DialOption, tags map[string]string, queueSize int) Publisher {
	ctx, cancel := context.WithCancel(context.Background())

	p := &publisher{
		connectTimeout: 5 * time.Second,
		publishTimeout: 2 * time.Second,

		serverAddress: serverAddress,
		dialOption:    dialOption,
		tags:          tags,

		publishChannel: make(chan *api.PublishRequest, queueSize),
		doneChannel:    make(chan interface{}),
		cancel:         cancel,
	}

	go p.startPublisher(ctx)

	return p
}

func (p *publisher) startPublisher(ctx context.Context) {
StartLoop:
	for {
		select {
		case <-ctx.Done():
			log.Debugf("shutdown publisher loop")
			break StartLoop
		default:
			p.startConnection(ctx)
		}
	}

	log.Debugf("shutdown publisher")
	p.doneChannel <- nil
}

func (p *publisher) startConnection(ctx context.Context) {
	connection, err := func() (*grpc.ClientConn, error) {
		dialCtx, cancel := context.WithTimeout(ctx, p.connectTimeout)
		defer cancel()
		return grpc.DialContext(
			dialCtx,
			p.serverAddress,
			p.dialOption,
			grpc.WithBlock(),
			grpc.WithBackoffMaxDelay(5*time.Second),
		)
	}()

	if err != nil {
		log.Errorf("could not connect to [%v]: %v", p.serverAddress, err)
		return
	}
	defer connection.Close()

	client := api.NewAggregatorClient(connection)

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
				return err
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
			log.Debugf("shutdown connection loop")
			break PublishLoop
		}
	}

	log.Debugf("shutdown connection")
}

func (p *publisher) Queue(ctx context.Context, ts time.Time, point collector.Point) error {
	payload := api.PublishRequest{
		UnixTime:      uint64(ts.Unix()),
		MilliCpuUsage: point.MilliCpuUsage,
		MemoryTotalMb: point.MemoryTotalMb,
		MemoryRssMb:   point.MemoryRssMb,
		MemoryLimitMb: point.MemoryLimitMb,
		DiskUsageMb:   point.DiskUsageMb,
		DiskLimitMb:   point.DiskLimitMb,
		Running:       point.Running,
	}

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
