package influxdb

import (
	"fmt"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter/writer"
	client "github.com/influxdata/influxdb/client/v2"
	"time"
)

const (
	timeout = 30 * time.Second
)

type influxdbClient interface {
	Write(bp client.BatchPoints) error
	Close() error
}

type influxdbWriter struct {
	database string

	client influxdbClient
}

func Open(config *Config) (writer.CloseWriter, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
		Timeout:  timeout,
	})

	if err != nil {
		return nil, fmt.Errorf("invalid InfluxDb configuration: %v", err)
	}

	return &influxdbWriter{
		database: config.Database,
		client:   client,
	}, nil
}

func (w *influxdbWriter) Write(batch batch.Batch) error {
	bp := buildBatchPoints(w.database, batch.Entries)
	return w.client.Write(bp)
}

func (w *influxdbWriter) Close() error {
	return w.client.Close()
}
