package jsonvector

import (
	"unsafe"
)

// Skip formatting symbols like tabs, spaces, ...
//
// Returns index of next non-format symbol.
// DEPRECATED: use skipFmtTable instead.
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
func skipFmtTable(src []byte, n, offset uint32) (uint32, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if n-offset > 512 {
		offset, _ = skipFmtBin8(src, n, offset)
	}
	for ; skipTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

// Binary based approach of skipFmt.
func skipFmtBin8(src []byte, n, offset uint32) (uint32, bool) {
	_ = src[n-1]
	_ = skipTable[255]
	if *(*uint64)(unsafe.Pointer(&src[offset])) == binNlSpace7 {
		offset += 8
		for offset < n && *(*uint64)(unsafe.Pointer(&src[offset])) == binSpace8 {
			offset += 8
		}
	}
	return offset, false
}

var (
	skipTable   = [256]bool{}
	binNlSpace7 uint64
	binSpace8   uint64
)

func init() {
	skipTable[' '] = true
	skipTable['\t'] = true
	skipTable['\n'] = true
	skipTable['\t'] = true

	binNlSpace7Bytes, binSpace8Bytes := []byte("\n       "), []byte("        ")
	binNlSpace7, binSpace8 = *(*uint64)(unsafe.Pointer(&binNlSpace7Bytes[0])), *(*uint64)(unsafe.Pointer(&binSpace8Bytes[0]))
}

var _ = skipFmt
