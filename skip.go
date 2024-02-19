package jsonvector

// Skip formatting symbols like tabs, spaces, ...
//
// Returns index of next non-format symbol.
func skipFmt(src []byte, n, offset int) (int, bool) {
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

// Table based approach of skipFmt.
func skipFmtTable(src []byte, n, offset int) (int, bool) {
	if src[offset] > ' ' {
		return offset, false
	}
	_ = src[n-1]
	_ = trimTable[255]
	for ; trimTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

var trimTable = [256]bool{}

func init() {
	trimTable[' '] = true
	trimTable['\t'] = true
	trimTable['\n'] = true
	trimTable['\t'] = true
}

var _ = skipFmt
