package utils

import (
	"testing"
	"time"
)

func TestTimestampConversion(t *testing.T) {
	nowTS := CurrentTimeMs()
	now := CurrentTime()
	got := TSToTimeIfNotSet(nil, &nowTS)
	if now.Sub(got) > time.Millisecond {
		t.Errorf("Timestamp mismatch %s - %s", now, got)
	}
}
