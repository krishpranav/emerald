package cart

import "fmt"

type Header struct {
	Entry     [4]byte
	Title     string
	GameCode  string
	MakerCode string
}

func New(src []byte) *Header {
	return &Header{
		Entry:     [4]byte{src[0], src[1], src[2], src[3]},
		Title:     string(src[0xa0 : 0xa0+12]),
		GameCode:  string(src[0xac : 0xac+4]),
		MakerCode: string(src[0xb0 : 0xb0+2]),
	}
}

func (h *Header) String() string {
	str := `Title: %s
GameCode: %s
MakerCode: %s`
	return fmt.Sprintf(str, h.Title, h.GameCode, h.MakerCode)
}
