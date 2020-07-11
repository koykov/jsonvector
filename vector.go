package jsonvector

import (
	"bytes"
	"errors"
	"reflect"
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
	s    []byte
	a    uint64
	v    []Val
	l    int
	r, c []int
}

var (
	bNull  = []byte("null")
	bTrue  = []byte("true")
	bFalse = []byte("false")
	bQuote = []byte(`"`)

	ErrEmptySrc = errors.New("can't parse empty source")
	ErrUnexpId  = errors.New("unexpected identifier")
	ErrUnexpEOF = errors.New("unexpected end of file")
	ErrEOA      = errors.New("end of array")
	ErrEOO      = errors.New("end of object")
)

func NewVector() *Vector {
	return &Vector{}
}

func (vec *Vector) Parse(s []byte) (err error) {
	if len(s) == 0 {
		err = ErrEmptySrc
		return
	}
	vec.s = s
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.s))
	vec.a = uint64(h.Data)

	offset := 0
	for offset < len(vec.s) {
		val := vec.newVal()
		i := vec.l - 1
		vec.r = append(vec.r, i)
		offset, err = vec.parse(offset, val)
		if err != nil {
			return err
		}
		vec.v[i] = *val
	}

	return
}

func (vec *Vector) Len() int {
	return len(vec.r)
}

func (vec *Vector) Get(keys ...string) *Val {
	if len(keys) == 0 {
		if vec.Len() > 0 {
			return &vec.v[0]
		}
		return nil
	}

	root := fastconv.S2B(keys[0])
	tail := keys[1:]
	_ = vec.r[len(vec.r)-1]
	for i := 0; i < len(vec.r); i++ {
		v := &vec.v[i]
		if bytes.Equal(root, v.k.Bytes()) {
			if len(tail) == 0 {
				return v
			} else {
				v.p = vec.ptr()
				return v.Get(tail...)
			}
		}
	}
	return nil
}

func (vec *Vector) newVal() (r *Val) {
	if vec.l < len(vec.v) {
		r = &vec.v[vec.l]
		r.Reset()
		vec.l++
	} else {
		r = &Val{t: TypeUnk}
		vec.v = append(vec.v, *r)
		vec.l++
	}
	return
}

func (vec *Vector) parse(offset int, v *Val) (int, error) {
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
		offset, err = vec.parseO(offset, v)
	case vec.s[offset] == '[':
		v.t = TypeArr
		offset, err = vec.parseA(offset, v)
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

func (vec *Vector) parseO(offset int, v *Val) (int, error) {
	v.cs = len(vec.c)
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
		c := vec.newVal()
		i := vec.l - 1
		vec.c = append(vec.c, i)
		v.ce = len(vec.c)
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
		offset, err = vec.parse(offset, c)
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

func (vec *Vector) parseA(offset int, v *Val) (int, error) {
	v.cs = len(vec.c)
	offset++
	var err error
	for offset < len(vec.s) {
		if vec.s[offset] == ']' {
			offset++
			break
		}
		c := vec.newVal()
		i := vec.l - 1
		vec.c = append(vec.c, i)
		v.ce = len(vec.c)
		offset, err = vec.parse(offset, c)
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
	vec.r = vec.r[:0]
	vec.c = vec.c[:0]
}

func (vec *Vector) ptr() uintptr {
	return uintptr(unsafe.Pointer(vec))
}

func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
