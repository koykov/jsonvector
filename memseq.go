package jsonvector

import (
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

type memseq struct {
	o uint64
	l int
	e bool
}

func (m *memseq) set(o uint64, l int) {
	m.o, m.l = o, l
}

func (m *memseq) Bytes() []byte {
	h := reflect.SliceHeader{
		Data: uintptr(m.o),
		Len:  m.l,
		Cap:  m.l,
	}
	p := *(*[]byte)(unsafe.Pointer(&h))
	if m.e {
		p = unescape(p)
		m.l = len(p)
		m.e = false
	}
	return p
}

func (m *memseq) String() string {
	return fastconv.B2S(m.Bytes())
}

func unescape(p []byte) []byte {
	l, i := len(p), 0
	for {
		i = bytealg.IndexAt(p, bEQuote, i)
		if i < 0 {
			break
		}
		copy(p[i:], p[i+1:])
		l--
	}
	return p[:l]
}
