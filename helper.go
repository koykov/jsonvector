package jsonvector

import (
	"github.com/koykov/vector"
)

var (
	jsonHelper = &JsonHelper{}
)

type JsonHelper struct{}

func (h *JsonHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.GetFlag(vector.FlagEscape) {
		return unescape(b)
	}
	return b
}
