package util

import (
	"unicode"
)

var (
	rangeCJKAndPunct = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x3000, 0x303f, 1},
			{0x3400, 0x4db5, 1},
			{0x4e00, 0x9fd5, 1},
			{0xfa0e, 0xfa0f, 1},
			{0xfa11, 0xfa11, 1},
			{0xfa13, 0xfa14, 1},
			{0xfa1f, 0xfa1f, 1},
			{0xfa21, 0xfa21, 1},
			{0xfa23, 0xfa24, 1},
			{0xfa27, 0xfa29, 1},
		},
		R32: []unicode.Range32{
			{Lo: 0x20000, Hi: 0x2a6d6, Stride: 1},
			{Lo: 0x2a700, Hi: 0x2b734, Stride: 1},
			{Lo: 0x2b740, Hi: 0x2b81d, Stride: 1},
			{Lo: 0x2b820, Hi: 0x2cea1, Stride: 1},
		},
	}
)

// ContainsCJK returns whether a string contains any CJK characters.
func ContainsCJK(str string) bool {
	for _, r := range []rune(str) {
		if unicode.Is(rangeCJKAndPunct, r) || r == '。' {
			return true
		}
	}

	return false
}
