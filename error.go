package jsonvector

import "errors"

var (
	ErrEmptySrc     = errors.New("can't parse empty source")
	ErrUnparsedTail = errors.New("unparsed tail")
	ErrUnexpId      = errors.New("unexpected identifier")
	ErrUnexpEOF     = errors.New("unexpected end of file")
	ErrUnexpEOS     = errors.New("unexpected end of string")
	ErrIncompatType = errors.New("incompatible type")
	ErrNotFound     = errors.New("node not found")
	ErrInternal     = errors.New("internal vector error")
)
