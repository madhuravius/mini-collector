package writer

import (
	"github.com/aptible/mini-collector/batch"
)

type Writer interface {
	Write(batch batch.Batch) error
}

type CloseWriter interface {
	Writer
	Close() error
}
