package jsonvector

import "errors"

var (
	ErrEmptySrc     = errors.New("can't parse empty source")
	ErrUnparsedTail = errors.New("unparsed tail")
	ErrUnexpId      = errors.New("unexpected identifier")
	ErrUnexpEOF     = errors.New("unexpected end of file")
	ErrUnexpEOS     = errors.New("unexpected end of string")
	ErrEOA          = errors.New("end of array")
	ErrEOO          = errors.New("end of object")
)
