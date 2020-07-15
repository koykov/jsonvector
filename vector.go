package jsonvector

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
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
	v []Val
	l int
	// Registry.
	r  [][]int
	rl int
}

var (
	bNull   = []byte("null")
	bTrue   = []byte("true")
	bFalse  = []byte("false")
	bQuote  = []byte(`"`)
	bEQuote = []byte(`\"`)

	ErrEmptySrc = errors.New("can't parse empty source")
	ErrUnexpId  = errors.New("unexpected identifier")
	ErrUnexpEOF = errors.New("unexpected end of file")
	ErrEOA      = errors.New("end of array")
	ErrEOO      = errors.New("end of object")
)

func NewVector() *Vector {
	return &Vector{}
}

func (vec *Vector) Parse(s []byte, copy bool) (err error) {
	if len(s) == 0 {
		err = ErrEmptySrc
		return
	}
	if copy {
		vec.s = append(vec.s[:0], s...)
	} else {
		vec.s = s
	}
	vec.s = s
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.s))
	vec.a = uint64(h.Data)
	vec.p = vec.ptr()

	offset := 0
	for offset < len(vec.s) {
		val := vec.newVal(0)
		i := vec.l - 1
		vec.reg(0, i)
		offset, err = vec.parse(0, offset, val)
		if err != nil {
			return err
		}
		vec.v[i] = *val
	}

	return
}

func (vec *Vector) ParseStr(s string, copy bool) error {
	return vec.Parse(fastconv.S2B(s), copy)
}

func (vec *Vector) Len() int {
	return vec.l
}

func (vec *Vector) Get(keys ...string) *Val {
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
		return vec.getA(r, keys...)
	}
	if r.t == TypeObj {
		return vec.getO(r, keys...)
	}
	return r
}

func (vec *Vector) getA(root *Val, keys ...string) *Val {
	if len(keys) == 0 {
		return root
	}
	k, err := strconv.Atoi(keys[0])
	if err != nil || k >= root.Len() {
		return nil
	}
	i := vec.r[root.d+1][k]
	v := &vec.v[i]
	tail := keys[1:]
	if v.t != TypeArr && v.t != TypeObj {
		if len(tail) > 0 {
			// Attempt to get child of scalar value.
			return nil
		}
		return v
	}
	if v.t == TypeArr {
		return vec.getA(v, tail...)
	}
	if v.t == TypeObj {
		return vec.getO(v, tail...)
	}
	return nil
}

func (vec *Vector) getO(root *Val, keys ...string) *Val {
	if len(keys) == 0 {
		return root
	}
	var v *Val
	for i := root.cs; i < root.ce; i++ {
		k := vec.r[root.d+1][i]
		v = &vec.v[k]
		if bytes.Equal(v.k.Bytes(), fastconv.S2B(keys[0])) {
			break
		}
	}
	if v == nil {
		return v
	}
	tail := keys[1:]
	if v.t != TypeArr && v.t != TypeObj {
		if len(tail) > 0 {
			// Attempt to get child of scalar value.
			return nil
		}
		return v
	}
	if v.t == TypeArr {
		return vec.getA(v, tail...)
	}
	if v.t == TypeObj {
		return vec.getO(v, tail...)
	}
	return nil
}

func (vec *Vector) newVal(depth int) (r *Val) {
	if vec.l < len(vec.v) {
		r = &vec.v[vec.l]
		r.Reset()
		vec.l++
	} else {
		r = &Val{t: TypeUnk}
		vec.v = append(vec.v, *r)
		vec.l++
	}
	r.p = vec.p
	r.d = depth
	return
}

func (vec *Vector) parse(depth, offset int, v *Val) (int, error) {
	var err error
	switch {
	case vec.s[offset] == 'n':
		if len(vec.s[offset:]) > 3 && bytes.Equal(bNull, vec.s[offset:offset+4]) {
			v.t = TypeNull
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case vec.s[offset] == '{':
		v.t = TypeObj
		offset, err = vec.parseO(depth+1, offset, v)
	case vec.s[offset] == '[':
		v.t = TypeArr
		offset, err = vec.parseA(depth+1, offset, v)
	case vec.s[offset] == ']':
		return offset, ErrEOA
	case vec.s[offset] == '"':
		v.t = TypeStr
		v.v.o = vec.a + uint64(offset+1)
		e := bytealg.IndexAt(vec.s, bQuote, offset+1)
		if vec.s[e-1] != '\\' {
			v.v.l = e - offset - 1
			offset = e + 1
		} else {
			_ = vec.s[len(vec.s)-1]
			for i := e; i < len(vec.s); {
				i = bytealg.IndexAt(vec.s, bQuote, i+1)
				e = i
				if vec.s[e-1] != '\\' {
					break
				}
			}
			v.v.l = e - offset - 1
			v.v.e = true
			offset += e + 1
		}
	case isDigit(vec.s[offset]):
		if len(vec.s[offset:]) > 0 {
			i := offset
			for isDigitDot(vec.s[i]) {
				i++
				if i == len(vec.s) {
					break
				}
			}
			v.t = TypeNum
			v.v.set(vec.a+uint64(offset), i-offset)
			offset = i
		} else {
			return offset, ErrUnexpEOF
		}
	case vec.s[offset] == 't':
		if len(vec.s[offset:]) > 3 && bytes.Equal(bTrue, vec.s[offset:offset+4]) {
			v.t = TypeBool
			v.v.set(vec.a+uint64(offset), 4)
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case vec.s[offset] == 'f':
		if len(vec.s[offset:]) > 4 && bytes.Equal(bFalse, vec.s[offset:offset+5]) {
			v.t = TypeBool
			v.v.set(vec.a+uint64(offset), 5)
			offset += 5
		} else {
			return offset, ErrUnexpId
		}
	}
	return offset, err
}

func (vec *Vector) parseO(depth, offset int, v *Val) (int, error) {
	v.cs = vec.regLen(depth)
	offset++
	var err error
	for offset < len(vec.s) {
		if vec.s[offset] == '}' {
			offset++
			break
		}
		for vec.s[offset] == ' ' {
			offset++
		}
		// Parse key.
		if vec.s[offset] != '"' {
			return offset, ErrUnexpId
		}
		offset++
		c := vec.newVal(depth)
		i := vec.l - 1
		v.ce = vec.reg(depth, i)
		c.k.o = vec.a + uint64(offset)
		e := bytealg.IndexAt(vec.s, bQuote, offset)
		if vec.s[e-1] != '\\' {
			c.k.l = e - offset
			offset = e + 1
		} else {
			_ = vec.s[len(vec.s)-1]
			for i := e; i < len(vec.s); {
				i = bytealg.IndexAt(vec.s, bQuote, i+1)
				e = i
				if vec.s[e-1] != '\\' {
					break
				}
			}
			c.k.l = e - offset
			c.k.e = true
			offset += e + 1
		}
		// Parse value.
		for vec.s[offset] == ' ' {
			offset++
		}
		if vec.s[offset] == ':' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		for vec.s[offset] == ' ' {
			offset++
		}
		offset, err = vec.parse(depth+1, offset, c)
		if err == ErrEOO {
			err = nil
			break
		}
		for vec.s[offset] == ',' || vec.s[offset] == ' ' {
			offset++
		}
		vec.v[i] = *c
	}
	return offset, err
}

func (vec *Vector) parseA(depth, offset int, v *Val) (int, error) {
	// v.cs = len(vec.c)
	v.cs = vec.regLen(depth)
	offset++
	var err error
	for offset < len(vec.s) {
		if vec.s[offset] == ']' {
			offset++
			break
		}
		c := vec.newVal(depth)
		i := vec.l - 1
		v.ce = vec.reg(depth, i)
		offset, err = vec.parse(depth+1, offset, c)
		if err == ErrEOA {
			err = nil
			break
		}
		for vec.s[offset] == ',' || vec.s[offset] == ' ' {
			offset++
		}
		vec.v[i] = *c
	}
	return offset, nil
}

func (vec *Vector) Reset() {
	if vec.l == 0 {
		return
	}
	_ = vec.v[vec.l-1]
	for i := 0; i < vec.l; i++ {
		vec.v[i].p = 0
	}
	vec.s = nil
	vec.a = 0
	vec.l = 0
	for i := 0; i < vec.rl; i++ {
		vec.r[i] = vec.r[i][:0]
	}
	vec.rl = 0
}

func (vec *Vector) ptr() uintptr {
	return uintptr(unsafe.Pointer(vec))
}

func (vec *Vector) reg(depth, idx int) int {
	if len(vec.r) <= depth {
		vec.r = append(vec.r, nil)
		vec.rl++
	}
	if vec.rl <= depth {
		vec.rl++
	}
	vec.r[depth] = append(vec.r[depth], idx)
	return len(vec.r[depth])
}

func (vec *Vector) regLen(depth int) int {
	if len(vec.r) <= depth {
		return 0
	}
	return len(vec.r[depth])
}

func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
