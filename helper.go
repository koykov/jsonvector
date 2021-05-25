package jsonvector

import (
	"github.com/koykov/vector"
)

const (
	flagEscape = uint8(1)
)

type JsonHelper struct{}

var (
	jsonHelper = &JsonHelper{}
)

func (h *JsonHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckFlag(flagEscape) {
		return unescape(b)
	}
	return b
}
