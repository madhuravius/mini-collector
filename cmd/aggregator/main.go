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
	"github.com/aptible/mini-collector/emitter/hold"
	"github.com/aptible/mini-collector/emitter/notify"
	"github.com/aptible/mini-collector/emitter/text"
	"github.com/aptible/mini-collector/emitter/writer"
	"github.com/aptible/mini-collector/tls"
	"github.com/aptible/mini-collector/writer/datadog"
	"github.com/aptible/mini-collector/writer/influxdb"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
		"host",
	}

	optionalTags = []string{
		"app",
		"database",
	}

	logger = logrus.WithFields(logrus.Fields{
		"source": "server",
	})
)

type server struct {
	api.UnimplementedAggregatorServer
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
		PublishRequest: point,
	})

	if err != nil {
		logger.Warnf("Ingest failed: %v", err)
	}

	return &api.PublishResponse{}, nil
}

func makeFinalEmitter() (emitter.Emitter, error) {
	notifyConfig := &notify.Config{}
	ok, err := tryLoadConfiguration("AGGREGATOR_NOTIFY_CONFIGURATION", notifyConfig)
	if err != nil {
		return nil, fmt.Errorf("could not decode notify configuration: %v", err)
	}

	if ok {
		return notify.Open(notifyConfig), nil
	}

	return blackhole.Open(), nil
}

func stackWriters(writerFactory func() (writer.CloseWriter, error), namePrefix string, count int) (emitter.Emitter, func(), error) {
	name := fmt.Sprintf("%s %d", namePrefix, count)

	w, err := writerFactory()
	if err != nil {
		return nil, nil, fmt.Errorf("writerFactory failed: %v", err)
	}

	if count <= 1 {
		finalEmitter, err := makeFinalEmitter()
		if err != nil {
			w.Close()
			return nil, nil, fmt.Errorf("makeFinalEmitter failed: %v", err)
		}

		em := writer.Open(name, w, finalEmitter)

		return em, func() {
			em.Close()
			w.Close()
			finalEmitter.Close()
		}, nil
	}

	nextCount := count - 1
	nextEmitter, closeNext, err := stackWriters(writerFactory, namePrefix, nextCount)

	if err != nil {
		w.Close()
		return nil, nil, fmt.Errorf("stackWriters(%d) failed: %v", nextCount, err)
	}

	// TODO: Backoff based on count
	holdEmitter := hold.Open(5*time.Second, nextEmitter)

	em := writer.Open(name, w, holdEmitter)

	return em, func() {
		em.Close()
		w.Close()
		holdEmitter.Close()
		closeNext()
	}, nil
}

func tryLoadConfiguration(envVariable string, configStruct interface{}) (bool, error) {
	jsonConfiguration, ok := os.LookupEnv(envVariable)
	if !ok {
		return false, nil
	}

	err := json.Unmarshal([]byte(jsonConfiguration), configStruct)
	if err != nil {
		return false, err
	}

	return true, nil

}

func getEmitter() (emitter.Emitter, func(), error) {
	datadogConfig := &datadog.Config{}
	ok, err := tryLoadConfiguration("AGGREGATOR_DATADOG_CONFIGURATION", datadogConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode Datadog configuration: %v", err)
	}
	if ok {
		if datadogConfig.Timeout == "" {
			datadogConfig.Timeout = "30s"
		}
		// We're not going to use it here, but parse this just to be sure we won't hit errors later.
		// It's better for us to fail now when the aggregator is just starting up.
		_, err := time.ParseDuration(datadogConfig.Timeout)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid timeout (%s): %v", datadogConfig.Timeout, err)
		}

		retryCount, err := strconv.Atoi(datadogConfig.RetryCount)
		if err != nil {
			retryCount = 3
		}
		logger.Infof("using Datadog writer with retry count %v, timeout %v", retryCount, datadogConfig.Timeout)
		return stackWriters(func() (writer.CloseWriter, error) {
			return datadog.Open(datadogConfig)
		}, "Datadog", retryCount)
	}

	influxdbConfig := &influxdb.Config{}
	ok, err = tryLoadConfiguration("AGGREGATOR_INFLUXDB_CONFIGURATION", influxdbConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not decode InfluxDB configuration: %v", err)
	}
	if ok {
		logger.Infof("using InfluxDB writer")
		return stackWriters(func() (writer.CloseWriter, error) {
			return influxdb.Open(influxdbConfig)
		}, "InfluxDB", 3)
	}

	_, ok = os.LookupEnv("AGGREGATOR_TEXT_CONFIGURATION")
	if ok {
		logger.Infof("using text emitter")
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

	logger.Infof("minPublishFrequency: %v", minPublishFrequency)

	// TODO: Make batchsize configurable?
	return batcher.New(em, minPublishFrequency, 1000), nil

}

func main() {
	grpcLogrus.ReplaceGrpcLogger(logger)

	emitter, closeCallback, err := getEmitter()
	if err != nil {
		logger.Fatalf("getEmitterStack failed: %v", err)
	}
	defer closeCallback()

	batcher, err := getBatcher(emitter)
	if err != nil {
		logger.Fatalf("getBatcher failed: %v", err)
	}
	defer batcher.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	logger.Infof("listening on: %s", port)

	var srv *grpc.Server

	_, enableTls := os.LookupEnv("AGGREGATOR_TLS")
	if enableTls {
		tlsConfig, err := tls.GetTlsConfig("AGGREGATOR")
		if err != nil {
			logger.Fatalf("failed to load tlsConfig: %v", err)
		}

		logger.Info("tls is enabled")
		srv = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	} else {
		logger.Warn("tls is disabled")
		srv = grpc.NewServer()
	}

	api.RegisterAggregatorServer(srv, &server{
		batcher: batcher,
	})

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		termSig := <-termChan
		logger.Infof("received %s, shutting down", termSig)
		srv.GracefulStop()
	}()

	if err := srv.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}

	logger.Infof("server shutdown")
}
