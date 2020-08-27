package jsonvector

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
)

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if len(s) == 0 {
		err = ErrEmptySrc
		return
	}
	s = bytealg.Trim(s, bFmt)
	if copy {
		vec.s = append(vec.s[:0], s...)
	} else {
		vec.s = s
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.s))
	vec.a = uint64(h.Data)
	vec.p = vec.ptr()

	offset := 0
	val := vec.newNode(0)
	i := vec.l - 1
	vec.reg(0, i)
	offset, err = vec.parseGeneric(0, offset, val)
	if err != nil {
		vec.e = offset
		return err
	}
	vec.v[i] = *val
	if offset < len(vec.s) {
		vec.e = offset
		return ErrUnparsedTail
	}

	return
}

func (vec *Vector) parseGeneric(depth, offset int, v *Node) (int, error) {
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
		offset, err = vec.parseObj(depth+1, offset, v)
	case vec.s[offset] == '[':
		v.t = TypeArr
		offset, err = vec.parseArr(depth+1, offset, v)
	case vec.s[offset] == ']':
		return offset, ErrEOA
	case vec.s[offset] == '"':
		v.t = TypeStr
		v.v.o = vec.a + uint64(offset+1)
		e := bytealg.IndexAt(vec.s, bQuote, offset+1)
		if e < 0 {
			return len(vec.s), ErrUnexpEOS
		}
		v.v.e = false
		if vec.s[e-1] != '\\' {
			v.v.l = e - offset - 1
			offset = e + 1
		} else {
			_ = vec.s[len(vec.s)-1]
			for i := e; i < len(vec.s); {
				i = bytealg.IndexAt(vec.s, bQuote, i+1)
				if i < 0 {
					e = len(vec.s) - 1
					break
				}
				e = i
				if vec.s[e-1] != '\\' {
					break
				}
			}
			v.v.l = e - offset - 1
			v.v.e = true
			offset = e + 1
		}
		if !v.v.e {
			v.v.e = bytes.IndexByte(v.v.rawBytes(), '\\') >= 0
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
			vec.e = offset
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
	default:
		err = ErrUnexpId
	}
	return offset, err
}

func (vec *Vector) parseObj(depth, offset int, v *Node) (int, error) {
	v.cs = vec.regLen(depth)
	offset++
	var err error
	for offset < len(vec.s) {
		if vec.s[offset] == '}' {
			// Edge case: empty object.
			offset++
			break
		}
		offset = vec.skipFmt(offset)
		// Parse key.
		if vec.s[offset] != '"' {
			return offset, ErrUnexpId
		}
		offset++
		c := vec.newNode(depth)
		i := vec.l - 1
		v.ce = vec.reg(depth, i)
		c.k.o = vec.a + uint64(offset)
		e := bytealg.IndexAt(vec.s, bQuote, offset)
		if e < 0 {
			return len(vec.s), ErrUnexpEOS
		}
		c.k.e = false
		if vec.s[e-1] != '\\' {
			c.k.l = e - offset
			offset = e + 1
		} else {
			_ = vec.s[len(vec.s)-1]
			for i := e; i < len(vec.s); {
				i = bytealg.IndexAt(vec.s, bQuote, i+1)
				if i < 0 {
					e = len(vec.s) - 1
					break
				}
				e = i
				if vec.s[e-1] != '\\' {
					break
				}
			}
			c.k.l = e - offset
			c.k.e = true
			offset = e + 1
		}
		if !c.k.e {
			c.k.e = bytes.IndexByte(c.k.rawBytes(), '\\') >= 0
		}
		// Parse value.
		offset = vec.skipFmt(offset)
		if vec.s[offset] == ':' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = vec.skipFmt(offset)
		offset, err = vec.parseGeneric(depth, offset, c)
		if err == ErrEOO {
			err = nil
			break
		}
		vec.v[i] = *c
		offset = vec.skipFmt(offset)
		if vec.s[offset] == '}' {
			offset++
			break
		}
		if vec.s[offset] == ',' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = vec.skipFmt(offset)
	}
	return offset, err
}

func (vec *Vector) parseArr(depth, offset int, v *Node) (int, error) {
	v.cs = vec.regLen(depth)
	offset++
	var err error
	for offset < len(vec.s) {
		if vec.s[offset] == ']' {
			// Edge case: empty array.
			offset++
			break
		}
		c := vec.newNode(depth)
		i := vec.l - 1
		v.ce = vec.reg(depth, i)
		offset, err = vec.parseGeneric(depth, offset, c)
		if err == ErrEOA {
			err = nil
			break
		}
		vec.v[i] = *c
		offset = vec.skipFmt(offset)
		if vec.s[offset] == ']' {
			offset++
			break
		}
		if vec.s[offset] == ',' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = vec.skipFmt(offset)
	}
	return offset, nil
}

func (vec *Vector) skipFmt(offset int) int {
	for bytes.IndexByte(bFmt, vec.s[offset]) != -1 {
		offset++
	}
	return offset
}

func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
