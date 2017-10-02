//+build !windows

package util

import (
	"time"
)

func init() {
	start := time.Now()
	MilliClock = func() time.Duration { return time.Since(start) }
}
