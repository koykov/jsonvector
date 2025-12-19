package jsonvector

import "github.com/koykov/bytealg"

func skipfmt(p []byte, off int) (int, bool) {
	if p[off] > 0x20 {
		return off, false
	}
	return bytealg.SkipBytesFmt4(p, off)
}
