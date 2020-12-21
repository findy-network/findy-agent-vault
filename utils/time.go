package utils

import (
	"time"
)

func CurrentTimeMs() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func CurrentTime() time.Time {
	return time.Now().UTC()
}

func TimestampToTime(tsMs *int64) *time.Time {
	var t *time.Time
	if tsMs != nil {
		secs := *tsMs / time.Second.Milliseconds()
		msecs := *tsMs - secs*time.Second.Milliseconds()
		ts := time.Unix(secs, msecs*time.Millisecond.Nanoseconds()).UTC()
		t = &ts
	}
	return t
}
