// Auto-generated code. DO NOT EDIT.
package publisher

import (
	"github.com/aptible/mini-collector/api"
	"github.com/aptible/mini-collector/collector"
	"time"
)

func buildPublishRequest(ts time.Time, point collector.Point) api.PublishRequest {
	return api.PublishRequest{
		UnixTime: ts.Unix(),
		Running:  point.Running,

		MilliCpuUsage: point.MilliCpuUsage,

		MemoryTotalMb: point.MemoryTotalMb,

		MemoryRssMb: point.MemoryRssMb,

		MemoryLimitMb: point.MemoryLimitMb,

		DiskUsageMb: point.DiskUsageMb,

		DiskLimitMb: point.DiskLimitMb,

		DiskReadKbps: point.DiskReadKbps,

		DiskWriteKbps: point.DiskWriteKbps,

		DiskReadIops: point.DiskReadIops,

		DiskWriteIops: point.DiskWriteIops,
	}
}
