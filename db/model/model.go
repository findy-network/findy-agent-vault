package model

import (
	"math"
	"strconv"
	"time"
)

func timeToString(t *time.Time) string {
	const timeLen = 10
	if t != nil && !t.IsZero() {
		return strconv.FormatInt(t.UnixNano()/time.Millisecond.Nanoseconds(), timeLen)
	}
	return ""
}

func timeToStringPtr(t *time.Time) *string {
	if t != nil && !t.IsZero() {
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

func TimeToCursor(t *time.Time) uint64 {
	return uint64(math.Round(float64(t.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
}
