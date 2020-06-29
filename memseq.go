package jsonvector

import (
	"reflect"
	"unsafe"
)

type memseq struct {
	o uint64
	l int
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
	return *(*[]byte)(unsafe.Pointer(&h))
}
