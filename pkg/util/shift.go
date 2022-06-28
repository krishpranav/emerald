package util

func LSL(val uint32, shiftAmount uint) uint32 { return val << shiftAmount }

func LSR(val uint32, shiftAmount uint) uint32 { return val >> shiftAmount }

func ASR(val uint32, shiftAmount uint) uint32 {
	msb := val & 0x8000_0000
	for i := uint(0); i < shiftAmount; i++ {
		val = (val >> 1) | msb
	}

	return val 
}