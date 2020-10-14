package jsonvector

import (
	"bytes"

	"github.com/koykov/bytealg"
)

func Validate(s []byte) (offset int, err error) {
	if len(s) == 0 {
		err = ErrEmptySrc
		return
	}
	s = bytealg.Trim(s, bFmt)
	offset, err = validateGeneric(0, s, offset)

	if offset < len(s) {
		err = ErrUnparsedTail
	}
	return
}

func validateGeneric(depth int, s []byte, offset int) (int, error) {
	var err error

	switch {
	case s[offset] == 'n':
		if len(s[offset:]) > 3 && bytes.Equal(bNull, s[offset:offset+4]) {
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case s[offset] == '{':
		offset, err = validateObj(depth+1, s, offset)
	case s[offset] == '[':
		offset, err = validateArr(depth+1, s, offset)
	case s[offset] == '"':
		e := bytealg.IndexAt(s, bQuote, offset+1)
		if e < 0 {
			return len(s), ErrUnexpEOS
		}
		if s[e-1] != '\\' {
			offset = e + 1
		} else {
			_ = s[len(s)-1]
			for i := e; i < len(s); {
				i = bytealg.IndexAt(s, bQuote, i+1)
				if i < 0 {
					e = len(s) - 1
					break
				}
				e = i
				if s[e-1] != '\\' {
					break
				}
			}
			offset = e + 1
		}
	case isDigit(s[offset]):
		if len(s[offset:]) > 0 {
			i := offset
			for isDigitDot(s[i]) {
				i++
				if i == len(s) {
					break
				}
			}
			offset = i
		} else {
			return offset, ErrUnexpEOF
		}
	case s[offset] == 't':
		if len(s[offset:]) > 3 && bytes.Equal(bTrue, s[offset:offset+4]) {
			offset += 4
		} else {
			return offset, ErrUnexpId
		}
	case s[offset] == 'f':
		if len(s[offset:]) > 4 && bytes.Equal(bFalse, s[offset:offset+5]) {
			offset += 5
		} else {
			return offset, ErrUnexpId
		}
	default:
		err = ErrUnexpId
	}

	return offset, err
}

func validateObj(depth int, s []byte, offset int) (int, error) {
	offset++
	var err error
	for offset < len(s) {
		if s[offset] == '}' {
			offset++
			break
		}
		offset = skipFmt(s, offset)
		// Parse key.
		if s[offset] != '"' {
			// Key should be a string wrapped with double quotas.
			return offset, ErrUnexpId
		}
		offset++
		e := bytealg.IndexAt(s, bQuote, offset)
		if e < 0 {
			return len(s), ErrUnexpEOS
		}
		if s[e-1] != '\\' {
			offset = e + 1
		} else {
			_ = s[len(s)-1]
			for i := e; i < len(s); {
				i = bytealg.IndexAt(s, bQuote, i+1)
				if i < 0 {
					e = len(s) - 1
					break
				}
				e = i
				if s[e-1] != '\\' {
					break
				}
			}
			offset = e + 1
		}
		offset = skipFmt(s, offset)
		if s[offset] == ':' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = skipFmt(s, offset)
		offset, err = validateGeneric(depth, s, offset)
		if err == ErrEOO {
			err = nil
			break
		}
		offset = skipFmt(s, offset)
		if s[offset] == '}' {
			offset++
			break
		}
		if s[offset] == ',' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = skipFmt(s, offset)
	}
	return offset, err
}

func validateArr(depth int, s []byte, offset int) (int, error) {
	offset++
	var err error
	for offset < len(s) {
		if s[offset] == ']' {
			offset++
			break
		}
		offset, err = validateGeneric(depth, s, offset)
		if err == ErrEOA {
			err = nil
			break
		}
		offset = skipFmt(s, offset)
		if s[offset] == ']' {
			// End of the array caught.
			offset++
			break
		}
		if s[offset] == ',' {
			offset++
		} else {
			return offset, ErrUnexpId
		}
		offset = skipFmt(s, offset)
	}
	return offset, nil
}

func skipFmt(s []byte, offset int) int {
	for bytes.IndexByte(bFmt, s[offset]) != -1 {
		offset++
	}
	return offset
}
