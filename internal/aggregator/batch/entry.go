package batch

import (
	"github.com/aptible/mini-collector/protobufs"
	"time"
)

type Entry struct {
	Time time.Time
	Tags map[string]string
	*protobufs.PublishRequest
}
