package jsonvector

// AppendEscape escapes byte array and add result to dst.
func AppendEscape(dst, b []byte) []byte {
	var o int
	l := len(b)
	if l == 0 {
		return dst
	}
	_ = b[l-1]
	for i := 0; i < l; i++ {
		c := b[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != 0x7F {
			continue
		}

		dst = append(dst, b[o:i]...)
		switch c {
		case '"':
			dst = append(dst, '\\', '"')
		case '\\':
			dst = append(dst, '\\', '\\')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\b':
			dst = append(dst, '\\', 'b')
		case 0x7F: // DEL
			dst = append(dst, `\u007f`...)
		default:
			dst = append(dst, `\u00`...)
			dst = appendHex(dst, c)
		}
		o = i + 1
	}

	return append(dst, b[o:]...)
}

func appendHex(dst []byte, c byte) []byte {
	const hex = "0123456789abcdef"
	return append(dst, hex[c>>4], hex[c&0x0F])
}
