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

func ROR(val uint32, shiftAmount uint) uint32 {
	shiftAmount %= 32
	tmp0 := (val) >> (shiftAmount)        
	tmp1 := (val) << (32 - (shiftAmount)) 
	return tmp0 | tmp1                    
}