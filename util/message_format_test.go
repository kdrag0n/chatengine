package util

import (
	"testing"
	"time"
)

var (
	_result         string
	spaceTestString = "hello there,                              human    . i like spaces        ! and CON tract ion s oa hwq idh    m            do not will not i am I am i'm dont willnt maynt shouldnt should not         . a jm. .d. qw.dqw .sad .asd. asd. as.d as.d .sad as ."
	expStRes = "Hello there, human. I like spaces! And CON tract ion s oa hwq idh m don't won't I'm I'm I'm don't willn't mayn't shouldn't shouldn't. A jm. D. Qw. Dqw. Sad. Asd. Asd. As. D as. D. Sad as."
)

func BenchmarkFormatFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := FormatMsg(spaceTestString, false)
		if f != expStRes {
			b.Error("Unexpected output from FormatMsg.", f, "!=", expStRes)
		}
	}
}

func BenchmarkSleep1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Second)
	}
}

func TestFormatMsg(t *testing.T) {
	f := FormatMsg(spaceTestString, false)
	if f != expStRes {
		t.Error("Unexpected output from FormatMsg.", f, "!=", expStRes)
	}
}