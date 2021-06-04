package utils

import (
	"time"
)

var CurrentStaticTime = time.Time{}

func CurrentTimeMs() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func CurrentTime() time.Time {
	if !CurrentStaticTime.IsZero() {
		return CurrentStaticTime
	}
	return time.Now().UTC()
}

func TSToTimeIfNotSet(current *time.Time, tsMs *int64) time.Time {
	if (current == nil || current.IsZero()) && tsMs != nil {
		secs := *tsMs / time.Second.Milliseconds()
		msecs := *tsMs - secs*time.Second.Milliseconds()
		ts := time.Unix(secs, msecs*time.Millisecond.Nanoseconds()).UTC()
		return ts
	}
	return *current
}

func TSToTimePtrIfNotSet(current *time.Time, tsMs *int64) *time.Time {
	if (current == nil || current.IsZero()) && tsMs != nil {
		secs := *tsMs / time.Second.Milliseconds()
		msecs := *tsMs - secs*time.Second.Milliseconds()
		ts := time.Unix(secs, msecs*time.Millisecond.Nanoseconds()).UTC()
		return &ts
	}
	return current
}
