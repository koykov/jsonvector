package jsonvector

import "errors"

var (
	ErrUnparsedTail = errors.New("unparsed tail")
	ErrUnexpId      = errors.New("unexpected identifier")
	ErrUnexpEOF     = errors.New("unexpected end of file")
	ErrUnexpEOS     = errors.New("unexpected end of string")
)
