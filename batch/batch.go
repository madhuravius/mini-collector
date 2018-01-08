package batch

import (
	"github.com/sirupsen/logrus"
)

type Batch struct {
	Id      uint64
	Entries []Entry
}

func (b *Batch) Fields() logrus.Fields {
	return logrus.Fields{"batch": b.Id}
}
