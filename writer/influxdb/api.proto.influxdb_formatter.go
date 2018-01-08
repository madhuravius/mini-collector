// Auto-generated code. DO NOT EDIT.
package influxdb

import (
	"github.com/aptible/mini-collector/batch"
)

func entryToFields(entry *batch.Entry) map[string]interface{} {
	out := map[string]interface{}{
		"running": (*entry).Running,
	}

	// NOTE: We report everything as int64 because older versions of
	// InfluxDB do not support uint64 as a type.

	var val int64

	val = int64((*entry).MilliCpuUsage)
	if val >= 0 {
		out["milli_cpu_usage"] = val
	}

	val = int64((*entry).MemoryTotalMb)
	if val >= 0 {
		out["memory_total_mb"] = val
	}

	val = int64((*entry).MemoryRssMb)
	if val >= 0 {
		out["memory_rss_mb"] = val
	}

	val = int64((*entry).MemoryLimitMb)
	if val >= 0 {
		out["memory_limit_mb"] = val
	}

	val = int64((*entry).DiskUsageMb)
	if val >= 0 {
		out["disk_usage_mb"] = val
	}

	val = int64((*entry).DiskLimitMb)
	if val >= 0 {
		out["disk_limit_mb"] = val
	}

	val = int64((*entry).DiskReadKbps)
	if val >= 0 {
		out["disk_read_kbps"] = val
	}

	val = int64((*entry).DiskWriteKbps)
	if val >= 0 {
		out["disk_write_kbps"] = val
	}

	val = int64((*entry).DiskReadIops)
	if val >= 0 {
		out["disk_read_iops"] = val
	}

	val = int64((*entry).DiskWriteIops)
	if val >= 0 {
		out["disk_write_iops"] = val
	}

	return out
}
