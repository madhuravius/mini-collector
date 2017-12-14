package main

import (
	"context"
	"github.com/aptible/mini-collector/collector"
	"github.com/aptible/mini-collector/publisher"
	"github.com/aptible/mini-collector/tls"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	queueTimeout = time.Second
)

func getEnvOrFatal(k string) string {
	val, ok := os.LookupEnv(k)
	if !ok {
		log.Fatalf("%s must be set", k)
	}
	return val
}

func main() {
	grpcLogrus.ReplaceGrpcLogger(log.NewEntry(log.StandardLogger()))

	// TODO: Throttling stats
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	serverAddress := getEnvOrFatal("MINI_COLLECTOR_REMOTE_ADDRESS")
	containerId := getEnvOrFatal("MINI_COLLECTOR_CONTAINER_ID")
	environmentName := getEnvOrFatal("MINI_COLLECTOR_ENVIRONMENT_NAME")
	serviceName := getEnvOrFatal("MINI_COLLECTOR_SERVICE_NAME")

	tags := map[string]string{
		"environment": environmentName,
		"service":     serviceName,
		"container":   containerId,
	}

	appName, ok := os.LookupEnv("MINI_COLLECTOR_APP_NAME")
	if ok {
		tags["app"] = appName
	}

	databaseName, ok := os.LookupEnv("MINI_COLLECTOR_DATABASE_NAME")
	if ok {
		tags["database"] = databaseName
	}

	cgroupPath, ok := os.LookupEnv("MINI_COLLECTOR_CGROUP_PATH")
	if !ok {
		cgroupPath = "/sys/fs/cgroup"
	}

	mountPath, ok := os.LookupEnv("MINI_COLLECTOR_MOUNT_PATH")
	if !ok {
		mountPath = ""
	}

	pollIntervalText, ok := os.LookupEnv("MINI_COLLECTOR_POLL_INTERVAL")
	if !ok {
		pollIntervalText = "30s"
	}

	pollInterval, err := time.ParseDuration(pollIntervalText)
	if err != nil {
		log.Fatalf("invalid poll interval (%s): %v", pollIntervalText, err)
	}

	_, debug := os.LookupEnv("MINI_COLLECTOR_DEBUG")
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	var dialOption grpc.DialOption

	_, enableTls := os.LookupEnv("MINI_COLLECTOR_TLS")
	if enableTls {
		var err error
		tlsConfig, err := tls.GetTlsConfig("MINI_COLLECTOR")
		if err != nil {
			log.Fatalf("failed to load tlsConfig: %v", err)
		}
		log.Info("tls is enabled")
		dialOption = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	} else {
		log.Warn("tls is disabled")
		dialOption = grpc.WithInsecure()
	}

	publisher, err := publisher.Open(
		&publisher.Config{
			ServerAddress: serverAddress,
			DialOption:    dialOption,
			Tags:          tags,
		},
	)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}
	defer publisher.Close()

	c := collector.NewCollector(cgroupPath, containerId, mountPath)

	log.Infof("pollInterval: %s", pollInterval)
	log.Infof("containerId: %s", containerId)
	log.Infof("mountPath: %s", mountPath)

	lastPoll := time.Now()
	lastState := collector.MakeNoContainerState(lastPoll)

MainLoop:
	for {
		nextPoll := lastPoll.Add(pollInterval)

		select {
		case <-time.After(time.Until(nextPoll)):
			lastPoll = nextPoll

			point, thisState, err := c.GetPoint(lastState)

			if err != nil {
				log.Warnf("GetPoint failed: %v", err)
				continue MainLoop
			}

			lastState = thisState

			err = func() error {
				ctx, cancel := context.WithTimeout(context.Background(), queueTimeout)
				defer cancel()
				return publisher.Queue(ctx, thisState.Time, point)
			}()

			if err != nil {
				log.Warnf("Queue failed: %v", err)
			}
		case <-termChan:
			// Exit
			log.Infof("exiting")
			break MainLoop
		}
	}
}
