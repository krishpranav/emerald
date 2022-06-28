package cart

import "fmt"

type Header struct {
	Entry      [4]byte
	Title      string
	GameCode   string
	MarkerCode string
}

func New(src []byte) *Header {
	return &Header{
		Entry:      [4]byte{src[0], src[1], src[2], src[3]},
		Title:      string(src[0xa0 : 0xa0+12]),
		GameCode:   string(src[0xac : 0xac+4]),
		MarkerCode: string(src[0xac : 0xac+2]),
	}
}

func (h *Header) String() string {
	str := `Title: \nGameCode: %s`
	return fmt.Sprintf(str, h.Title, h.GameCode, h.MarkerCode)
}
