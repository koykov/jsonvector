package jsonvector

import (
	"unicode/utf16"

	"github.com/koykov/vector"
)

// Unescape byte array using itself as a destination.
func Unescape(p []byte) []byte {
	l := len(p)
	var i int
	for {
		i = vector.IndexByteAt(p, '\\', i)
		if i < 0 || i+1 == l {
			break
		}
		switch p[i+1] {
		case '\\':
			// Escaped slash caught, just copy it over slash.
			copy(p[i:], p[i+1:])
			i++
		case '"', '/':
			// Caught " and /, just copy them over slash.
			copy(p[i:], p[i+1:])
		case 'n':
			// Caught new line symbol, unescape it and copy over slash.
			p[i] = '\n'
			copy(p[i+1:], p[i+2:])
		case 'r':
			// Caught caret return symbol, unescape it and copy over slash.
			p[i] = '\r'
			copy(p[i+1:], p[i+2:])
		case 't':
			// Caught tab symbol, unescape it and copy over slash.
			p[i] = '\t'
			copy(p[i+1:], p[i+2:])
		case 'b':
			// Caught backspace symbol, unescape it and copy over slash.
			p[i] = '\b'
			copy(p[i+1:], p[i+2:])
		case 'f':
			// Caught form feed symbol, unescape it and copy over slash.
			p[i] = '\f'
			copy(p[i+1:], p[i+2:])
		case 'u':
			// Caught unicode symbol.
			if l-i < 6 {
				i++
				continue
			}
			x := p[i+2 : i+6]
			u := xtouTable(x)
			r := rune(u)
			if !utf16.IsSurrogate(r) {
				// Regular utf8 symbol.
				s := string(r)
				z := len(s)
				copy(p[i:], s)
				copy(p[i+z:], p[i+6:])
				l -= 5 - z
				i += z
			} else {
				// Caught surrogate pair, see https://en.wikipedia.org/wiki/UTF-16#Code_points_from_U+010000_to_U+10FFFF
				if l-i < 12 {
					i++
					continue
				}
				if p[i+6] != '\\' || p[i+7] != 'u' {
					i++
					continue
				} else {
					x = p[i+8 : i+12]
					u1 := xtouTable(x)
					r = utf16.DecodeRune(r, rune(u1))
					s := string(r)
					z := len(s)
					copy(p[i:], s)
					copy(p[i+z:], p[i+12:])
					l -= 11 - z
					i += z
				}
			}
		}
		l--
		p = p[:l]
	}
	return p
}

// AppendUnescape unescapes byte array and add result to dst.
func AppendUnescape(dst, p []byte) []byte {
	l := len(p)
	var i int
	for {
		i = vector.IndexByteAt(p, '\\', i)
		if i < 0 || i+1 == l {
			break
		}
		switch p[i+1] {
		case '\\':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\\')
			p = p[i+2:]
			l = len(p)
			i = 0
		case '"', '/':
			dst = append(dst, p[:i]...)
			dst = append(dst, p[i+1])
			p = p[i+2:]
			l = len(p)
			i = 0
		case 'n':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\n')
			p = p[i+2:]
			l = len(p)
			i = 0
		case 'r':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\r')
			p = p[i+2:]
			l = len(p)
			i = 0
		case 't':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\t')
			p = p[i+2:]
			l = len(p)
			i = 0
		case 'b':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\b')
			p = p[i+2:]
			l = len(p)
			i = 0
		case 'f':
			dst = append(dst, p[:i]...)
			dst = append(dst, '\f')
			p = p[i+2:]
			l = len(p)
			i = 0
		case 'u':
			if l-i < 6 {
				dst = append(dst, p[:i+1]...)
				p = p[i+1:]
				l = len(p)
				i = 0
				continue
			}
			x := p[i+2 : i+6]
			u := xtouTable(x)
			r := rune(u)
			if !utf16.IsSurrogate(r) {
				s := string(r)
				dst = append(dst, p[:i]...)
				dst = append(dst, s...)
				p = p[i+6:]
				l = len(p)
				i = 0
			} else {
				if l-i < 12 {
					dst = append(dst, p[:i+1]...)
					p = p[i+1:]
					l = len(p)
					i = 0
					continue
				}
				if p[i+6] != '\\' || p[i+7] != 'u' {
					dst = append(dst, p[:i+1]...)
					p = p[i+1:]
					l = len(p)
					i = 0
					continue
				}
				x = p[i+8 : i+12]
				u1 := xtouTable(x)
				r = utf16.DecodeRune(r, rune(u1))
				s := string(r)
				dst = append(dst, p[:i]...)
				dst = append(dst, s...)
				p = p[i+12:]
				l = len(p)
				i = 0
			}
		default:
			dst = append(dst, p[:i+1]...)
			p = p[i+1:]
			l = len(p)
			i = 0
		}
	}
	dst = append(dst, p...)
	return dst
}
