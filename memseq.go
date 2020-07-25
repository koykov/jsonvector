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

func unescape(p []byte) []byte {
	l, i := len(p), 0
	for {
		i = bytealg.IndexAt(p, bSlash, i)
		if i < 0 || i+1 == l {
			break
		}
		switch p[i+1] {
		case '\\':
			copy(p[i:], p[i+1:])
			i++
		case '"', '/':
			copy(p[i:], p[i+1:])
		case 'n':
			p[i] = '\n'
			copy(p[i+1:], p[i+2:])
		case 'r':
			p[i] = '\r'
			copy(p[i+1:], p[i+2:])
		case 't':
			p[i] = '\t'
			copy(p[i+1:], p[i+2:])
		case 'b':
			p[i] = '\b'
			copy(p[i+1:], p[i+2:])
		case 'f':
			p[i] = '\f'
			copy(p[i+1:], p[i+2:])
		}
		l--
		p = p[:l]
	}
	return p
}
