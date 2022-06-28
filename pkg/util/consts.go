package util

type Region uint32

const (
	BIOS Region = 0x0
	EWRAM Region = 0x2
)

var RegionSize = map[string]uint32 {
	"BIOS": 0x00004000,
}