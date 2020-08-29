package jsonvector

import (
	"io"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

type Type int

const (
	TypeUnk = iota
	TypeNull
	TypeObj
	TypeArr
	TypeStr
	TypeNum
	TypeBool
)

type Vector struct {
	s []byte
	p uintptr
	a uint64
	v []Node
	l int
	e int
	r registry

	ss []string
}

var (
	bNull  = []byte("null")
	bTrue  = []byte("true")
	bFalse = []byte("false")
	bQuote = []byte(`"`)
	bSlash = []byte(`\`)
	bFmt   = []byte(" \t\n\r")
)

func NewVector() *Vector {
	return &Vector{}
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

func (vec *Vector) Len() int {
	return vec.l
}

func (vec *Vector) ErrorOffset() int {
	return vec.e
}

func (vec *Vector) Root() *Node {
	return vec.Get()
}

func (vec *Vector) Get(keys ...string) *Node {
	if len(keys) == 0 {
		if vec.Len() > 0 {
			return &vec.v[0]
		}
		return nil
	}

	r := &vec.v[0]
	if r.t != TypeObj && r.t != TypeArr {
		if len(keys) > 1 {
			// Attempt to get child of scalar value.
			return nil
		}
		return r
	}

	if r.t == TypeArr {
		return vec.getArr(r, keys...)
	}
	if r.t == TypeObj {
		return vec.getObj(r, keys...)
	}
	return r
}

func (vec *Vector) GetByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	return vec.Get(vec.ss...)
}

// Format vector in human readable representation.
func (vec *Vector) Beautify(w io.Writer) error {
	r := &vec.v[0]
	return vec.beautify(w, r, 0)
}

func (vec *Vector) Exists(key string) bool {
	n := vec.Get()
	return n.Exists(key)
}

func (vec *Vector) newNode(depth int) (r *Node) {
	if vec.l < len(vec.v) {
		r = &vec.v[vec.l]
		r.Reset()
		vec.l++
	} else {
		r = &Node{t: TypeUnk}
		vec.v = append(vec.v, *r)
		vec.l++
	}
	r.p, r.d = vec.p, depth
	return
}

func (vec *Vector) Reset() {
	if vec.l == 0 {
		return
	}
	_ = vec.v[vec.l-1]
	for i := 0; i < vec.l; i++ {
		vec.v[i].p = 0
	}
	vec.s, vec.a, vec.l, vec.e = nil, 0, 0, 0
	vec.r.reset()
}

func (vec *Vector) ptr() uintptr {
	return uintptr(unsafe.Pointer(vec))
}
