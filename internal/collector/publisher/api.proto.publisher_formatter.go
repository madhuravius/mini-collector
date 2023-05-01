// Auto-generated code. DO NOT EDIT.
package publisher

import (
	"time"
	"github.com/aptible/mini-collector/protobufs"
	"github.com/aptible/mini-collector/internal/collector"
)

func buildPublishRequest(ts time.Time, point collector.Point) protobufs.PublishRequest {
	return protobufs.PublishRequest{
		UnixTime:      ts.Unix(),
		Running:       point.Running,

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
		PidsCurrent: point.PidsCurrent,
		PidsLimit: point.PidsLimit,
		
	}
}
