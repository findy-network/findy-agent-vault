package utils

import (
	"testing"
	"time"
)

func TestTimestampConversion(t *testing.T) {
	nowTS := CurrentTimeMs()
	now := CurrentTime()
	got := TimestampToTime(&nowTS)
	if now.Sub(*got) > time.Millisecond {
		t.Errorf("Timestamp mismatch %s - %s", now, *got)
	}
}
