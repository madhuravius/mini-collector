package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/batcher"
	"github.com/aptible/mini-collector/emitter"
	"github.com/aptible/mini-collector/emitter/blackhole"
	"github.com/aptible/mini-collector/emitter/text"
	"github.com/aptible/mini-collector/emitter/writer"
	"github.com/aptible/mini-collector/tls"
	"github.com/aptible/mini-collector/writer/influxdb"
	"github.com/aptible/mini-collector/writer/datadog"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	ts := time.Unix(point.UnixTime, 0)

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

	err := s.batcher.Ingest(ctx, &batch.Entry{
		Time:           ts,
		Tags:           tags,
		PublishRequest: *point,
	})

	if err != nil {
		log.Warnf("Ingest failed: %v", err)
	}

	return &api.PublishResponse{}, nil
}

func stackWriters(writerFactory func() (writer.CloseWriter, error), namePrefix string, count int) (emitter.Emitter, func(), error) {
	name := fmt.Sprintf("%s %d", namePrefix, count)

	w, err := writerFactory()
	if err != nil {
		return nil, nil, fmt.Errorf("writerFactory failed: %v", err)
	}

	if count <= 1 {
		em := writer.Open(name, w, blackhole.Open())
		return em, func() {
			em.Close()
			w.Close()
		}, nil
	}

	nextCount := count - 1
	nextEmitter, closeNext, err := stackWriters(writerFactory, namePrefix, nextCount)

	if err != nil {
		w.Close()
		return nil, nil, fmt.Errorf("stackWriters(%d) failed: %v", nextCount, err)
	}

	em := writer.Open(name, w, nextEmitter)

	return em, func() {
		em.Close()
		w.Close()
		closeNext()
	}, nil
}

func getEmitter() (emitter.Emitter, func(), error) {
	// TODO: Extract this into a function shared by writers
	datadogConfiguration, ok := os.LookupEnv("AGGREGATOR_DATADOG_CONFIGURATION")
	if ok {
		log.Infof("using Datadog writer")

		config := &datadog.Config{}
		err := json.Unmarshal([]byte(datadogConfiguration), &config)
		if err != nil {
			return nil, nil, fmt.Errorf("could not decode Datadog configuration: %v", err)
		}

		writerFactory := func() (writer.CloseWriter, error) {
			return datadog.Open(config)
		}

		return stackWriters(writerFactory, "Datadog", 3)
	}

	influxDbConfiguration, ok := os.LookupEnv("AGGREGATOR_INFLUXDB_CONFIGURATION")
	if ok {
		log.Infof("using InfluxDB writer")

		config := &influxdb.Config{}
		err := json.Unmarshal([]byte(influxDbConfiguration), &config)
		if err != nil {
			return nil, nil, fmt.Errorf("could not decode InfluxDB configuration: %v", err)
		}

		writerFactory := func() (writer.CloseWriter, error) {
			return influxdb.Open(config)
		}

		return stackWriters(writerFactory, "InfluxDB", 3)
	}

	_, ok = os.LookupEnv("AGGREGATOR_TEXT_CONFIGURATION")
	if ok {
		log.Infof("using text emitter")
		em := text.Open()
		return em, em.Close, nil
	}

	return nil, nil, fmt.Errorf("no emitter configured")
}

func getBatcher(em emitter.Emitter) (batcher.Batcher, error) {
	minPublishFrequencyText, ok := os.LookupEnv("AGGREGATOR_MINIMUM_PUBLISH_FREQUENCY")
	if !ok {
		minPublishFrequencyText = "15s"
	}

	minPublishFrequency, err := time.ParseDuration(minPublishFrequencyText)
	if err != nil {
		return nil, fmt.Errorf("invalid minimum publish frequency (%s): %v", minPublishFrequencyText, err)
	}

	log.Infof("minPublishFrequency: %v", minPublishFrequency)

	// TODO: Make batchsize configurable?
	return batcher.New(em, minPublishFrequency, 1000), nil

}

func main() {
	grpcLogrus.ReplaceGrpcLogger(log.NewEntry(log.StandardLogger()))

	emitter, closeCallback, err := getEmitter()
	if err != nil {
		log.Fatalf("getEmitterStack failed: %v", err)
	}
	defer closeCallback()

	batcher, err := getBatcher(emitter)
	if err != nil {
		log.Fatalf("getBatcher failed: %v", err)
	}
	defer batcher.Close()

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

		log.Info("tls is enabled")
		srv = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	} else {
		log.Warn("tls is disabled")
		srv = grpc.NewServer()
	}

	api.RegisterAggregatorServer(srv, &server{
		batcher: batcher,
	})

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Infof("exiting")
}
