package jsonvector

import (
	"github.com/koykov/vector"
)

const (
	flagEscape = 0
)

type JsonHelper struct{}

var (
	jsonHelper = &JsonHelper{}
)

func (h *JsonHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = unescape(b)
		p.SetLimit(len(b))
	}
	return b
}
