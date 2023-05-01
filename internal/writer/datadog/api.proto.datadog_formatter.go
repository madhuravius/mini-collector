// Auto-generated code. DO NOT EDIT.
package datadog

import (
	"fmt"
	"github.com/aptible/mini-collector/internal/aggregator/batch"
)

func formatBatch(batch batch.Batch) datadogPayload {
	allSeries := make([]datadogSeries, 0, len(batch.Entries))

	var (
		series datadogSeries
		host string
		ok bool
	)

	for _, entry := range batch.Entries {
		tags := make([]string, 0, len(entry.Tags))

		for k, v := range entry.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s", k, v))
		}

		ts := entry.Time.Unix()

		var val int64

		
		
		val = int64(entry.UnixTime)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.unix_time",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.MilliCpuUsage)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.milli_cpu_usage",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.MemoryTotalMb)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.memory_total_mb",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.MemoryRssMb)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.memory_rss_mb",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.MemoryLimitMb)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.memory_limit_mb",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskUsageMb)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_usage_mb",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskLimitMb)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_limit_mb",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskReadKbps)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_read_kbps",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskWriteKbps)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_write_kbps",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskReadIops)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_read_iops",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.DiskWriteIops)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.disk_write_iops",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.PidsCurrent)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.pids_current",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		
		
		val = int64(entry.PidsLimit)
		

		if val >= 0 {
			series = datadogSeries{
				Metric: "enclave.pids_limit",
				Points: []datadogPoint{
					datadogPoint{ts, val},
				},
				Type: "gauge",
				Tags: tags,
			}

			host, ok = entry.Tags["host_name"]
			if ok {
				series.Host = host
			}

			allSeries = append(allSeries, series)
		}
		

	}

	return datadogPayload{Series: allSeries}
}
