package model

import (
	"math"
	"strconv"
	"time"
)

func timeToString(t *time.Time) string {
	if t != nil && !t.IsZero() {
		return strconv.FormatInt(t.UnixNano()/time.Millisecond.Nanoseconds(), 10)
	}
	return ""
}

func timeToStringPtr(t *time.Time) *string {
	if t != nil {
		res := timeToString(t)
		return &res
	}
	return nil
}

type Base struct {
	ID       string `faker:"uuid_hyphenated"`
	TenantID string
	Cursor   uint64
	Created  time.Time
}

func (b *Base) copy() *Base {
	baseCopy := *b
	return &baseCopy
}

func copyTime(t *time.Time) *time.Time {
	var res *time.Time
	if t != nil {
		ts := *t
		res = &ts
	}
	return res
}

func TimeToCursor(t *time.Time) uint64 {
	return uint64(math.Round(float64(t.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
}
