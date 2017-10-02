package util

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	before := MilliClock()
	if before < time.Microsecond {
		t.Error("Clock less than 1 µs")
		return
	}

	time.Sleep(time.Millisecond * 500)
	after := MilliClock()
	if after < before || (after - before) < time.Millisecond * 475 {
		t.Error("Clock didn't tick at leat 475ms after sleeping 500ms")
	}
}
