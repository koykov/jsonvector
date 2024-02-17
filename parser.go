package jsonvector

import (
	"bytes"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

var (
	// Byte constants.
	bNull  = []byte("null")
	bTrue  = []byte("true")
	bFalse = []byte("false")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if vec.Helper == nil {
		vec.Helper = helper
	}

	s = bytealg.TrimFmt4(s)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	offset := 0
	// Create root node and register it.
	root, i := vec.GetNode(0)

	// Parse source data.
	offset, err = vec.parseGeneric(0, offset, root)
	if err != nil {
		vec.SetErrOffset(offset)
		return err
	}
	vec.PutNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		return vector.ErrUnparsedTail
	}

	return
}

// Generic parser helper.
func (vec *Vector) parseGeneric(depth, offset int, node *vector.Node) (int, error) {
	var err error
	node.SetOffset(vec.Index.Len(depth))
	src := vec.Src()
	n := len(src)
	_ = src[n-1]
	switch {
	case src[offset] == 'n':
		// Check null node.
		if len(src[offset:]) > 3 && bytes.Equal(bNull, src[offset:offset+4]) {
			node.SetType(vector.TypeNull)
			offset += 4
		} else {
			return offset, vector.ErrUnexpId
		}
	case src[offset] == '{':
		// Check open object node.
		node.SetType(vector.TypeObj)
		offset, err = vec.parseObj(depth+1, offset, node)
	case src[offset] == '[':
		// Check open array node.
		node.SetType(vector.TypeArr)
		offset, err = vec.parseArr(depth+1, offset, node)
	case src[offset] == '"':
		// Check string node.
		node.SetType(vector.TypeStr)
		// Save offset of string value.
		node.Value().TakeAddr(src).SetOffset(offset + 1)
		// Get index of end of string value.
		e := bytealg.IndexByteAtBytes(src, '"', offset+1)
		if e < 0 {
			return n, vector.ErrUnexpEOS
		}
		node.Value().SetBit(flagEscape, true)
		if src[e-1] != '\\' {
			// Good case - string isn't escaped.
			node.Value().SetLen(e - offset - 1)
			offset = e + 1
		} else {
			// Walk over double quotas and look for unescaped.
			//_ = vec.Src()[vec.SrcLen()-1]
			for i := e; i < n; {
				i = bytealg.IndexByteAtBytes(src, '"', i+1)
				if i < 0 {
					e = n - 1
					break
				}
				e = i
				if src[e-1] != '\\' {
					break
				}
			}
			node.Value().SetLen(e - offset - 1)
			node.Value().SetBit(flagEscape, true)
			offset = e + 1
		}
		if !node.Value().CheckBit(flagEscape) {
			// Extra check of escaping sequences.
			node.Value().SetBit(flagEscape, bytealg.HasByteBytes(node.Value().RawBytes(), '\\'))
		}
	case isDigit(src[offset]):
		// Check number node.
		if len(src[offset:]) > 0 {
			// Get the edges of number.
			i := offset
			for isDigitDot(src[i]) {
				i++
				if i == n {
					break
				}
			}
			node.SetType(vector.TypeNum)
			node.Value().Init(src, offset, i-offset)
			offset = i
		} else {
			vec.SetErrOffset(offset)
			return offset, vector.ErrUnexpEOF
		}
	case src[offset] == 't':
		// Check bool (true) node.
		if len(src[offset:]) > 3 && bytes.Equal(bTrue, src[offset:offset+4]) {
			node.SetType(vector.TypeBool)
			node.Value().Init(src, offset, 4)
			offset += 4
		} else {
			return offset, vector.ErrUnexpId
		}
	case src[offset] == 'f':
		// Check bool (false) node.
		if len(src[offset:]) > 4 && bytes.Equal(bFalse, src[offset:offset+5]) {
			node.SetType(vector.TypeBool)
			node.Value().Init(src, offset, 5)
			offset += 5
		} else {
			return offset, vector.ErrUnexpId
		}
	default:
		err = vector.ErrUnexpId
	}
	return offset, err
}

// Object parsing helper.
func (vec *Vector) parseObj(depth, offset int, node *vector.Node) (int, error) {
	node.SetOffset(vec.Index.Len(depth))
	offset++
	var (
		err error
		eof bool
	)
	src := vec.Src()
	n := len(src)
	_ = src[n-1]
	for offset < n {
		if src[offset] == '}' {
			// Edge case: empty object.
			offset++
			break
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Parse key.
		if src[offset] != '"' {
			// Key should be a string wrapped with double quotas.
			return offset, vector.ErrUnexpId
		}
		offset++
		// Register new node.
		child, i := vec.GetChild(node, depth)
		// Fill up key's offset and length.
		child.Key().TakeAddr(src).SetOffset(offset)
		e := bytealg.IndexByteAtBytes(src, '"', offset+1)
		if e < 0 {
			return n, vector.ErrUnexpEOS
		}
		child.Key().SetBit(flagEscape, false)
		if src[e-1] != '\\' {
			// Key is an unescaped string, good case.
			child.Key().SetLen(e - offset)
			offset = e + 1
		} else {
			// Key contains escaped bytes.
			//_ = vec.Src()[vec.SrcLen()-1]
			for i := e; i < n; {
				i = bytealg.IndexByteAtBytes(src, '"', i+1)
				if i < 0 {
					e = n - 1
					break
				}
				e = i
				if src[e-1] != '\\' {
					break
				}
			}
			child.Key().SetLen(e - offset)
			child.Key().SetBit(flagEscape, true)
			offset = e + 1
		}
		if !child.Key().CheckBit(flagEscape) {
			// Extra check of escaped sequences in the key.
			child.Key().SetBit(flagEscape, bytealg.HasByteBytes(child.KeyBytes(), '\\'))
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Check division symbol.
		if src[offset] == ':' {
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Parse value.
		// Value may be an arbitrary type.
		if offset, err = vec.parseGeneric(depth, offset, child); err != nil {
			return offset, err
		}
		// Save node to the vector.
		vec.PutNode(i, child)
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if src[offset] == '}' {
			// End of the object caught.
			offset++
			break
		}
		if src[offset] == ',' {
			// Object elements separator caught.
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
	}
	return offset, err
}

// Array parsing helper.
func (vec *Vector) parseArr(depth, offset int, node *vector.Node) (int, error) {
	node.SetOffset(vec.Index.Len(depth))
	offset++
	var (
		err error
		eof bool
	)
	src := vec.Src()
	n := len(src)
	_ = src[n-1]
	for offset < n {
		if src[offset] == ']' {
			// Edge case: empty array.
			offset++
			break
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if src[offset] == ']' {
			// Edge case: empty array.
			offset++
			break
		}
		// Register new node.
		child, i := vec.GetChild(node, depth)
		// Parse the value.
		if offset, err = vec.parseGeneric(depth, offset, child); err != nil {
			return offset, err
		}
		// Save node to the vector.
		vec.PutNode(i, child)
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if src[offset] == ']' {
			// End of the array caught.
			offset++
			break
		}
		if src[offset] == ',' {
			// Array elements separator caught.
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(src, n, offset); eof {
			return offset, vector.ErrUnexpEOF
		}
	}
	return offset, nil
}

// Skip formatting symbols like tabs, spaces, ...
//
// Returns index of next non-format symbol.
func (vec *Vector) skipFmt(src []byte, n, offset int) (int, bool) {
	if src[offset] > ' ' {
		return offset, false
	}
	_ = src[n-1]
	for ; offset < n; offset++ {
		c := src[offset]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return offset, false
		}
	}
	return offset, true
}

// Check if given byte is a part of the number.
func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

// Check if given is a part of the number, including dot.
func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
