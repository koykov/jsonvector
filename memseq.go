package jsonvector

import (
	"reflect"
	"unsafe"

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
	p := m.unescBytes()
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

func (m *memseq) unescBytes() []byte {
	h := reflect.SliceHeader{
		Data: uintptr(m.o),
		Len:  m.l,
		Cap:  m.l,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}
