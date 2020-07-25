package jsonvector

import (
	"strconv"
	"unicode/utf16"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

func unescape(p []byte) []byte {
	l, i := len(p), 0
	for {
		i = bytealg.IndexAt(p, bSlash, i)
		if i < 0 || i+1 == l {
			break
		}
		switch p[i+1] {
		case '\\':
			copy(p[i:], p[i+1:])
			i++
		case '"', '/':
			copy(p[i:], p[i+1:])
		case 'n':
			p[i] = '\n'
			copy(p[i+1:], p[i+2:])
		case 'r':
			p[i] = '\r'
			copy(p[i+1:], p[i+2:])
		case 't':
			p[i] = '\t'
			copy(p[i+1:], p[i+2:])
		case 'b':
			p[i] = '\b'
			copy(p[i+1:], p[i+2:])
		case 'f':
			p[i] = '\f'
			copy(p[i+1:], p[i+2:])
		case 'u':
			if l-i < 6 {
				i++
				continue
			}
			x := p[i+2 : i+6]
			u, err := strconv.ParseUint(fastconv.B2S(x), 16, 16)
			if err != nil {
				i++
				continue
			}
			r := rune(u)
			if !utf16.IsSurrogate(r) {
				s := string(r)
				z := len(s)
				copy(p[i:], s)
				copy(p[i+z:], p[i+6:])
				l -= 5 - z
				i += z
			} else {
				if l-i < 12 {
					i++
					continue
				}
				if p[i+6] != '\\' || p[i+7] != 'u' {
					i++
					continue
				} else {
					x = p[i+8 : i+12]
					u1, err := strconv.ParseUint(fastconv.B2S(x), 16, 16)
					if err != nil {
						i++
						continue
					}
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
