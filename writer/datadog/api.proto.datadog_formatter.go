// Auto-generated code. DO NOT EDIT.
package datadog

import (
	"fmt"
	"github.com/aptible/mini-collector/batch"
)

func formatBatch(batch batch.Batch) datadogPayload {
	allSeries := make([]datadogSeries, 0, len(batch.Entries))

	var (
		series datadogSeries
		host   string
		ok     bool
	)

	for _, entry := range batch.Entries {
		tags := make([]string, 0, len(entry.Tags))

		for k, v := range entry.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s", k, v))
		}

		series = datadogSeries{
			Metric: "enclave.running",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.Running},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.milli_cpu_usage",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.MilliCpuUsage},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.memory_total_mb",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.MemoryTotalMb},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.memory_rss_mb",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.MemoryRssMb},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.memory_limit_mb",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.MemoryLimitMb},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_usage_mb",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskUsageMb},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_limit_mb",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskLimitMb},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_read_kbps",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskReadKbps},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_write_kbps",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskWriteKbps},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_read_iops",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskReadIops},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

		series = datadogSeries{
			Metric: "enclave.disk_write_iops",
			Points: []datadogPoint{
				datadogPoint{entry.Time.Unix(), entry.DiskWriteIops},
			},
			Type: "gauge",
			Tags: tags,
		}

		host, ok = entry.Tags["host"]
		if ok {
			series.Host = host
		}

		allSeries = append(allSeries, series)

	}

	return datadogPayload{Series: allSeries}
}
