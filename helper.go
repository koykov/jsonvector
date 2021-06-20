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

func (h *JsonHelper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = unescape(b)
		p.SetLen(len(b))
	}
	return b
}
