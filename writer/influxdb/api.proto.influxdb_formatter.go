// Auto-generated code. DO NOT EDIT.
package influxdb

import (
	"github.com/aptible/mini-collector/batch"
)

func entryToFields(entry *batch.Entry) map[string]interface{} {
	return map[string]interface{}{
		// NOTE: Older versions of InfluxDB do not support uint64 here.

		"running": (*entry).Running,

		"milli_cpu_usage": int64((*entry).MilliCpuUsage),

		"memory_total_mb": int64((*entry).MemoryTotalMb),

		"memory_rss_mb": int64((*entry).MemoryRssMb),

		"memory_limit_mb": int64((*entry).MemoryLimitMb),

		"disk_usage_mb": int64((*entry).DiskUsageMb),

		"disk_limit_mb": int64((*entry).DiskLimitMb),

		"disk_read_kbps": int64((*entry).DiskReadKbps),

		"disk_write_kbps": int64((*entry).DiskWriteKbps),

		"disk_read_iops": int64((*entry).DiskReadIops),

		"disk_write_iops": int64((*entry).DiskWriteIops),
	}
}
