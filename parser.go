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
	bFmt   = []byte(" \t\n\r")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bFmt)
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
	switch {
	case vec.SrcAt(offset) == 'n':
		// Check null node.
		if len(vec.Src()[offset:]) > 3 && bytes.Equal(bNull, vec.Src()[offset:offset+4]) {
			node.SetType(vector.TypeNull)
			offset += 4
		} else {
			return offset, vector.ErrUnexpId
		}
	case vec.SrcAt(offset) == '{':
		// Check open object node.
		node.SetType(vector.TypeObj)
		offset, err = vec.parseObj(depth+1, offset, node)
	case vec.SrcAt(offset) == '[':
		// Check open array node.
		node.SetType(vector.TypeArr)
		offset, err = vec.parseArr(depth+1, offset, node)
	case vec.SrcAt(offset) == '"':
		// Check string node.
		node.SetType(vector.TypeStr)
		// Save offset of string value.
		node.Value().TakeAddr(vec.Src()).SetOffset(offset + 1)
		// Get index of end of string value.
		e := bytealg.IndexByteAtLR(vec.Src(), '"', offset+1)
		if e < 0 {
			return vec.SrcLen(), vector.ErrUnexpEOS
		}
		node.Value().SetBit(flagEscape, true)
		if vec.SrcAt(e-1) != '\\' {
			// Good case - string isn't escaped.
			node.Value().SetLen(e - offset - 1)
			offset = e + 1
		} else {
			// Walk over double quotas and look for unescaped.
			_ = vec.Src()[vec.SrcLen()-1]
			for i := e; i < vec.SrcLen(); {
				i = bytealg.IndexByteAtLR(vec.Src(), '"', i+1)
				if i < 0 {
					e = vec.SrcLen() - 1
					break
				}
				e = i
				if vec.SrcAt(e-1) != '\\' {
					break
				}
			}
			node.Value().SetLen(e - offset - 1)
			node.Value().SetBit(flagEscape, true)
			offset = e + 1
		}
		if !node.Value().CheckBit(flagEscape) {
			// Extra check of escaping sequences.
			node.Value().SetBit(flagEscape, bytealg.HasByteLR(node.Value().RawBytes(), '\\'))
		}
	case isDigit(vec.SrcAt(offset)):
		// Check number node.
		if len(vec.Src()[offset:]) > 0 {
			// Get the edges of number.
			i := offset
			for isDigitDot(vec.SrcAt(i)) {
				i++
				if i == vec.SrcLen() {
					break
				}
			}
			node.SetType(vector.TypeNum)
			node.Value().Init(vec.Src(), offset, i-offset)
			offset = i
		} else {
			vec.SetErrOffset(offset)
			return offset, vector.ErrUnexpEOF
		}
	case vec.SrcAt(offset) == 't':
		// Check bool (true) node.
		if len(vec.Src()[offset:]) > 3 && bytes.Equal(bTrue, vec.Src()[offset:offset+4]) {
			node.SetType(vector.TypeBool)
			node.Value().Init(vec.Src(), offset, 4)
			offset += 4
		} else {
			return offset, vector.ErrUnexpId
		}
	case vec.SrcAt(offset) == 'f':
		// Check bool (false) node.
		if len(vec.Src()[offset:]) > 4 && bytes.Equal(bFalse, vec.Src()[offset:offset+5]) {
			node.SetType(vector.TypeBool)
			node.Value().Init(vec.Src(), offset, 5)
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
	for offset < vec.SrcLen() {
		if vec.SrcAt(offset) == '}' {
			// Edge case: empty object.
			offset++
			break
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Parse key.
		if vec.SrcAt(offset) != '"' {
			// Key should be a string wrapped with double quotas.
			return offset, vector.ErrUnexpId
		}
		offset++
		// Register new node.
		child, i := vec.GetChild(node, depth)
		// Fill up key's offset and length.
		child.Key().TakeAddr(vec.Src()).SetOffset(offset)
		e := bytealg.IndexByteAtLR(vec.Src(), '"', offset+1)
		if e < 0 {
			return vec.SrcLen(), vector.ErrUnexpEOS
		}
		child.Key().SetBit(flagEscape, false)
		if vec.SrcAt(e-1) != '\\' {
			// Key is an unescaped string, good case.
			child.Key().SetLen(e - offset)
			offset = e + 1
		} else {
			// Key contains escaped bytes.
			_ = vec.Src()[vec.SrcLen()-1]
			for i := e; i < len(vec.Src()); {
				i = bytealg.IndexByteAtLR(vec.Src(), '"', i+1)
				if i < 0 {
					e = vec.SrcLen() - 1
					break
				}
				e = i
				if vec.SrcAt(e-1) != '\\' {
					break
				}
			}
			child.Key().SetLen(e - offset)
			child.Key().SetBit(flagEscape, true)
			offset = e + 1
		}
		if !child.Key().CheckBit(flagEscape) {
			// Extra check of escaped sequences in the key.
			child.Key().SetBit(flagEscape, bytealg.HasByteLR(child.KeyBytes(), '\\'))
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Check division symbol.
		if vec.SrcAt(offset) == ':' {
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		// Parse value.
		// Value may be an arbitrary type.
		if offset, err = vec.parseGeneric(depth, offset, child); err != nil {
			return offset, err
		}
		// Save node to the vector.
		vec.PutNode(i, child)
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if vec.SrcAt(offset) == '}' {
			// End of the object caught.
			offset++
			break
		}
		if vec.SrcAt(offset) == ',' {
			// Object elements separator caught.
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
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
	for offset < vec.SrcLen() {
		if vec.SrcAt(offset) == ']' {
			// Edge case: empty array.
			offset++
			break
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if vec.SrcAt(offset) == ']' {
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
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
		if vec.SrcAt(offset) == ']' {
			// End of the array caught.
			offset++
			break
		}
		if vec.SrcAt(offset) == ',' {
			// Array elements separator caught.
			offset++
		} else {
			return offset, vector.ErrUnexpId
		}
		if offset, eof = vec.skipFmt(offset); eof {
			return offset, vector.ErrUnexpEOF
		}
	}
	return offset, nil
}

// Skip formatting symbols like tabs, spaces, ...
//
// Returns the next non-format symbol index.
func (vec *Vector) skipFmt(offset int) (int, bool) {
loop:
	if offset >= vec.SrcLen() {
		return offset, true
	}
	c := vec.SrcAt(offset)
	if c != bFmt[0] && c != bFmt[1] && c != bFmt[2] && c != bFmt[3] {
		return offset, false
	}
	offset++
	goto loop
}

// Check if given byte is a part of the number.
func isDigit(c byte) bool {
	return (c >= '0' && c <= '9') || c == '-' || c == '+' || c == 'e' || c == 'E'
}

// Check if given is a part of the number, including dot.
func isDigitDot(c byte) bool {
	return isDigit(c) || c == '.'
}
