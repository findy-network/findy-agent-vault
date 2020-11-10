package utils

import "time"

func CurrentTimeMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
