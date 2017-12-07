package main

import (
	"context"
	"fmt"
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/batcher"
	"github.com/aptible/mini-collector/emitter"
	"github.com/aptible/mini-collector/emitter/influxdb"
	"github.com/aptible/mini-collector/emitter/text"
	"github.com/aptible/mini-collector/tls"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"time"
)

const (
	port = ":8000"
)

var (
	requiredTags = []string{
		"environment",
		"service",
		"container",
	}

	optionalTags = []string{
		"app",
		"database",
	}
)

type server struct {
	batcher batcher.Batcher
}

func (s *server) Publish(ctx context.Context, point *api.PublishRequest) (*api.PublishResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("no metadata")
	}

	ts := time.Unix(int64(point.UnixTime), 0)

	tags := map[string]string{}

	for _, k := range requiredTags {
		v, ok := md[k]
		if !ok {
			return nil, fmt.Errorf("missing required metadata key: %s", k)
		}
		tags[k] = v[0]
	}

	for _, k := range optionalTags {
		v, ok := md[k]
		if !ok {
			continue
		}
		tags[k] = v[0]
	}

	s.batcher.Ingest(&batch.Entry{
		Time:           ts,
		Tags:           tags,
		PublishRequest: *point,
	})
	return &api.PublishResponse{}, nil
}

func getEmitter() (emitter.Emitter, error) {
	influxDbConfiguration, ok := os.LookupEnv("AGGREGATOR_INFLUXDB_CONFIGURATION")
	if ok {
		log.Infof("using InfluxDB emitter")
		return influxdb.New(influxDbConfiguration)
	}

	_, ok = os.LookupEnv("AGGREGATOR_TEXT_CONFIGURATION")
	if ok {
		log.Infof("using text emitter")
		return text.New()
	}

	return nil, fmt.Errorf("no emitter configured")
}

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	emitter, err := getEmitter()
	if err != nil {
		log.Fatalf("failed to get emitter: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on: %s", port)

	var srv *grpc.Server

	_, enableTls := os.LookupEnv("AGGREGATOR_TLS")
	if enableTls {
		tlsConfig, err := tls.GetTlsConfig("AGGREGATOR")
		if err != nil {
			log.Fatalf("failed to load tlsConfig: %v", err)
		}

		log.Infof("enabling tls")
		srv = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	} else {
		srv = grpc.NewServer()
	}

	api.RegisterAggregatorServer(srv, &server{
		batcher: batcher.New(emitter),
	})

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
