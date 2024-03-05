package jsonvector

import (
	"io"

	"github.com/koykov/vector"
)

const flagEscape = 0

type Helper struct{}

var helper = Helper{}

func (h Helper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = Unescape(b)
		p.SetLen(uint32(len(b)))
	}
	return b
}

func (h Helper) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (h Helper) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
