package util

import (
	"time"
)

// CurrentTimeMillis gets the current time in milliseconds since UNIX epoch.
func CurrentTimeMillis() int64 {
	return time.Now().Truncate(time.Millisecond).UnixNano() / 1000000
}
