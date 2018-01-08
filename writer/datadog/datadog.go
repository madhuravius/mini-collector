package datadog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aptible/mini-collector/batch"
	"github.com/aptible/mini-collector/emitter/writer"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	datadogSeriesUrl = "https://app.datadoghq.com/api/v1/series"
	timeout          = 30 * time.Second
)

type datadogWriter struct {
	apiKey string
}

func Open(config *Config) (writer.CloseWriter, error) {
	if config.ApiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}

	return &datadogWriter{
		apiKey: config.ApiKey,
	}, nil
}

func (em *datadogWriter) Write(batch batch.Batch) error {
	datadogPayload := formatBatch(batch)

	jsonPayload, err := json.Marshal(datadogPayload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?api_key=%s", datadogSeriesUrl, em.apiKey)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("NewRequest failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Do failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var bodyStr string
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		bodyStr = fmt.Sprintf("Body is not available: %v", err)
	} else {
		bodyStr = string(body)
	}

	return fmt.Errorf("Datadog POST failed: %s:\n%s", resp.Status, bodyStr)
}

func (em *datadogWriter) Close() error {
	return nil
}
