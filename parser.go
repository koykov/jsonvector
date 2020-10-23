package jsonvector

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if len(s) == 0 {
		err = ErrEmptySrc
		return
	}
	s = bytealg.Trim(s, bFmt)
	if copy {
		// Copy input data.
		vec.b = append(vec.b[:0], s...)
		vec.s = vec.b
	} else {
		// Use input data as source.
		vec.s = s
	}

	// Get source data address and raw parser pointer.
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.s))
	vec.a = uint64(h.Data)
	vec.p = vec.ptr()

	offset := 0
	// Create root node and register it.
	val := vec.newNode(0)
	i := vec.l - 1
	vec.r.reg(0, i)

	// Parse source data.
	offset, err = vec.parseGeneric(0, offset, val)
	if err != nil {
		vec.e = offset
		return err
	}
	vec.v[i] = *val

	// Check unparsed tail.
	if offset < len(vec.s) {
		vec.e = offset
		return ErrUnparsedTail
	}

	return
}

// Generic parser helper.
func (vec *Vector) parseGeneric(depth, offset int, v *Node) (int, error) {
	var err error
	v.s = vec.r.len(depth)
	switch {
	case vec.s[offset] == 'n':
		// Check null node.
		if len(vec.s[offset:]) > 3 && bytes.Equal(bNull, vec.s[offset:offset+4]) {
			v.t = TypeNull
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case vec.s[offset] == '{':
		// Check open object node.
		v.t = TypeObj
		offset, err = vec.parseObj(depth+1, offset, v)
	case vec.s[offset] == '[':
		// Check open array node.
		v.t = TypeArr
		offset, err = vec.parseArr(depth+1, offset, v)
	case vec.s[offset] == '"':
		// Check string node.
		v.t = TypeStr
		// Save offset of string value.
		v.v.o = vec.a + uint64(offset+1)
		// Get index of end of string value.
		e := bytealg.IndexAt(vec.s, bQuote, offset+1)
		if e < 0 {
			return len(vec.s), ErrUnexpEOS
		}
		v.v.e = false
		if vec.s[e-1] != '\\' {
			// Good case - string isn't escaped.
			v.v.l = e - offset - 1
			offset = e + 1
		} else {
			// Walk over double quotas and look for unescaped.
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
			// Extra check of escaping sequences.
			v.v.e = bytes.IndexByte(v.v.rawBytes(), '\\') >= 0
		}
	case isDigit(vec.s[offset]):
		// Check number node.
		if len(vec.s[offset:]) > 0 {
			// Get the edges of number.
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
		// Check bool (true) node.
		if len(vec.s[offset:]) > 3 && bytes.Equal(bTrue, vec.s[offset:offset+4]) {
			v.t = TypeBool
			v.v.set(vec.a+uint64(offset), 4)
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case vec.s[offset] == 'f':
		// Check bool (false) node.
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

// Object parsing helper.
func (vec *Vector) parseObj(depth, offset int, v *Node) (int, error) {
	v.s = vec.r.len(depth)
	offset++
	var (
		err error
		eof bool
	)
	for offset < len(vec.s) {
		if vec.s[offset] == '}' {
			// Edge case: empty object.
			offset++
			break
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
		// Parse key.
		if vec.s[offset] != '"' {
			// Key should be a string wrapped with double quotas.
			return offset, ErrUnexpId
		}
		offset++
		// Register new node.
		c := vec.newNode(depth)
		i := vec.l - 1
		v.e = vec.r.reg(depth, i)
		// Fill up key's offset and length.
		c.k.o = vec.a + uint64(offset)
		e := bytealg.IndexAt(vec.s, bQuote, offset)
		if e < 0 {
			return len(vec.s), ErrUnexpEOS
		}
		c.k.e = false
		if vec.s[e-1] != '\\' {
			// Key is an unescaped string, good case.
			c.k.l = e - offset
			offset = e + 1
		} else {
			// Key contains escaped bytes.
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
			// Extra check of escaped sequences in the key.
			c.k.e = bytes.IndexByte(c.k.rawBytes(), '\\') >= 0
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
		// Check division symbol.
		if vec.s[offset] == ':' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
		// Parse value.
		// Value may be an arbitrary type.
		if offset, err = vec.parseGeneric(depth, offset, c); err != nil {
			return offset, err
		}
		// Save node to the vector.
		vec.v[i] = *c
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
		if vec.s[offset] == '}' {
			// End of the object caught.
			offset++
			break
		}
		if vec.s[offset] == ',' {
			// Object elements separator caught.
			offset++
		} else {
			return offset, ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
	}
	return offset, err
}

// Array parsing helper.
func (vec *Vector) parseArr(depth, offset int, v *Node) (int, error) {
	v.s = vec.r.len(depth)
	offset++
	var (
		err error
		eof bool
	)
	for offset < len(vec.s) {
		if vec.s[offset] == ']' {
			// Edge case: empty array.
			offset++
			break
		}
		// Register new node.
		c := vec.newNode(depth)
		i := vec.l - 1
		v.e = vec.r.reg(depth, i)
		// Parse the value.
		if offset, err = vec.parseGeneric(depth, offset, c); err != nil {
			return offset, err
		}
		// Save node to the vector.
		vec.v[i] = *c
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
		if vec.s[offset] == ']' {
			// End of the array caught.
			offset++
			break
		}
		if vec.s[offset] == ',' {
			// Array elements separator caught.
			offset++
		} else {
			return offset, ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, ErrUnexpEOF
		}
	}
	return offset, nil
}

// Skip formatting symbols like tabs, spaces, ...
//
// Returns the next non-format symbol index.
func (vec *Vector) skipFmt(offset int) (int, bool) {
	if offset >= len(vec.s) {
		return offset, true
	}
	for bytes.IndexByte(bFmt, vec.s[offset]) != -1 {
		offset++
	}
	return offset, false
}

// Check if given byte is a part of the number.
func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

// Check if given is a part of the number, including dot.
func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
