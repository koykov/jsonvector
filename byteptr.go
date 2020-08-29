package jsonvector

import (
	"reflect"
	"unsafe"

	"github.com/koykov/fastconv"
)

type byteptr struct {
	o uint64
	l int
	e bool
}

func (m *byteptr) set(o uint64, l int) {
	m.o, m.l = o, l
}

func (m *byteptr) Bytes() []byte {
	p := m.rawBytes()
	if m.e {
		p = unescape(p)
		m.l = len(p)
		m.e = false
	}
	return p
}

func (m *byteptr) String() string {
	return fastconv.B2S(m.Bytes())
}

func (m *byteptr) rawBytes() []byte {
	h := reflect.SliceHeader{
		Data: uintptr(m.o),
		Len:  m.l,
		Cap:  m.l,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}
