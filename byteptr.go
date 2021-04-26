package jsonvector

import (
	"reflect"
	"unsafe"

	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

var (
	jsonUnesc = &JsonUnescapeHelper{}
)

type JsonUnescapeHelper struct{}

func (h *JsonUnescapeHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.GetFlag(vector.FlagEscape) {
		return unescape(b)
	}
	return b
}

// Byte sequence.
type byteptr struct {
	// Offset in virtual memory.
	o uint64
	// Length of bytes array.
	l int
	// Escaped byte sequence.
	e bool
}

// Set new offset and length.
func (m *byteptr) set(o uint64, l int) {
	m.o, m.l = o, l
}

// Convert byte sequence to byte array.
func (m *byteptr) Bytes() []byte {
	p := m.rawBytes()
	if m.e {
		// Unescape byte array.
		p = unescape(p)
		m.l = len(p)
		m.e = false
	}
	return p
}

// Convert byte sequence to string.
func (m *byteptr) String() string {
	return fastconv.B2S(m.Bytes())
}

// Convert byte sequence to byte array.
func (m *byteptr) rawBytes() []byte {
	h := reflect.SliceHeader{
		Data: uintptr(m.o),
		Len:  m.l,
		Cap:  m.l,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}
