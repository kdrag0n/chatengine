package util

import (
	"testing"
	"time"
)

func TestCurrentTimeMillis(t *testing.T) {
	tm := CurrentTimeMillis()
	if tm < time.Now().Unix() {
		t.Error("CurrentTimeMillis less than current time in seconds...")
		return
	}

	time.Sleep(time.Second)
	if CurrentTimeMillis() - tm < 950 {
		t.Error("CurrentTimeMillis didn't tick at least 950ms after sleeping 1s")
	}
}