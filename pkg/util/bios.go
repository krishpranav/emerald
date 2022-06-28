package util

import "math"

func Div(num, denom int32) (uint32, uint32, uint32) {
	const I32_MIN = -1
	if denom == 0 {
		if num < 0 {
			return uint32(0xffff_ffff), uint32(num), 1
		}
		return 1, uint32(num), 1
	} else if denom == -1 && num == I32_MIN {
		return 0x8000_0000, 0, 0x8000_0000
	} else {
		result := num / denom
		mod := num % denom

		return uint32(result), uint32(mod)
	}
}