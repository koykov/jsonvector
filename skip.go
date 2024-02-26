package jsonvector

import (
	"encoding/binary"
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
func skipFmtTable(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = trimTable[255]
	if n-offset > 512 {
		offset, _ = skipFmtTable1(src, n, offset)
	}
	for ; trimTable[src[offset]]; offset++ {
	}
	return offset, offset == n
}

func skipFmtTable1(src []byte, n, offset int) (int, bool) {
	_ = src[n-1]
	_ = trimTable[255]
	if *(*uint64)(unsafe.Pointer(&src[offset])) == binNlSpace7 {
		offset += 8
		for offset < n && *(*uint64)(unsafe.Pointer(&src[offset])) == binSpace8 {
			offset += 8
		}
	}
	if *(*uint32)(unsafe.Pointer(&src[offset])) == binNlSpace3 {
		offset += 4
		for offset < n && *(*uint32)(unsafe.Pointer(&src[offset])) == binSpace4 {
			offset += 4
		}
	}
	return offset, false
}

var (
	trimTable   = [256]bool{}
	binNlSpace3 = binary.LittleEndian.Uint32([]byte("\n   "))
	binSpace4   = binary.LittleEndian.Uint32([]byte("    "))
	binNlSpace7 = binary.LittleEndian.Uint64([]byte("\n       "))
	binSpace8   = binary.LittleEndian.Uint64([]byte("        "))
)

func init() {
	trimTable[' '] = true
	trimTable['\t'] = true
	trimTable['\n'] = true
	trimTable['\t'] = true
}

var _ = skipFmt
