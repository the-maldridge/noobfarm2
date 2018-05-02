package qdb

import (
	"time"
)

func (q *Quote) DisplayTime() string {
	return q.Submitted.Format(time.RFC1123)
}
