package main

import (
	"github.com/aptible/mini-collector/collector"
	"github.com/aptible/mini-collector/publisher"
	"github.com/aptible/mini-collector/tls"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	publisherBufferSize = 10
	// pollInterval        = 2 * time.Second // TODO
	pollInterval = 2000 * time.Millisecond // TODO
)

func getEnvOrFatal(k string) string {
	val, ok := os.LookupEnv(k)
	if !ok {
		log.Fatalf("%s must be set", k)
	}
	return val
}

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	// TODO: Volumes / configuration
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

	_, debug := os.LookupEnv("MINI_COLLECTOR_DEBUG")
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	cgroupPath, ok := os.LookupEnv("MINI_COLLECTOR_CGROUP_PATH")
	if !ok {
		cgroupPath = "/sys/fs/cgroup"
	}

	mountPath, ok := os.LookupEnv("MINI_COLLECTOR_MOUNT_PATH")
	if !ok {
		mountPath = ""
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
		log.Infof("enabling tls")
		dialOption = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	} else {
		dialOption = grpc.WithInsecure()
	}

	publisher := publisher.Open(
		serverAddress,
		dialOption,
		tags,
		20,
	)

	c := collector.NewCollector(cgroupPath, containerId, mountPath)

	lastState := collector.MakeNoContainerState()

MainLoop:
	for {
		select {
		case <-time.After(time.Until(lastState.Time.Add(pollInterval))):
			var point collector.Point
			point, lastState = c.GetPoint(lastState)
			err := publisher.Queue(lastState.Time, point)
			if err != nil {
				log.Warnf("publisher is failling behind: %v", err)
			}
		case <-termChan:
			// Exit
			log.Infof("exiting")
			break MainLoop
		}
	}

	publisher.Close()
}
