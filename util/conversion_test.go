package util

import (
	"testing"
	"bytes"
	"math/rand"
	crand "crypto/rand"
)

func TestBytesToString(t *testing.T) {
	lines := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		data := make([]byte, 15 + rand.Intn(35))
		var err error
		if i % 2 == 0 {
			_, err = rand.Read(data)
		} else {
			_, err = crand.Read(data)
		}

		if err != nil {
			t.Error("Error generating random data", err)
			return
		}

		lines[i] = data
	}

	for _, line := range lines {
		if !bytes.Equal(line, []byte(BytesToString(line))) {
			t.Error("BytesToString+reversed line didn't match original")
			return
		}
	}
}

func TestStringToBytes(t *testing.T) {
	lines := make([]string, 50)
	for i := 0; i < 50; i++ {
		data := make([]byte, 15 + rand.Intn(35))
		var err error
		if i % 2 == 0 {
			_, err = rand.Read(data)
		} else {
			_, err = crand.Read(data)
		}

		if err != nil {
			t.Error("Error generating random data", err)
			return
		}

		lines[i] = string(data)
	}

	for _, line := range lines {
		if !bytes.Equal([]byte(line), StringToBytes(line)) {
			t.Error("StringToBytes line didn't match original (data)")
			return
		}
	}
}