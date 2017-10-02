package util

import (
	"time"
)

// Clock returns the number of milliseconds that have elapsed since the program
// was started.
var MilliClock func() time.Duration
