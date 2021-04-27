package jsonvector

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

// Parser object.
type Vector struct {
	vector.Vector
}

var (
	// Byte constants.
	bNull  = []byte("null")
	bTrue  = []byte("true")
	bFalse = []byte("false")
	bQuote = []byte(`"`)
	bSlash = []byte(`\`)
	bFmt   = []byte(" \t\n\r")
)

// Make new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = jsonHelper
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// Parse source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// Copy source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// Copy source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

// Format vector in human readable representation.
func (vec *Vector) Beautify(w io.Writer) error {
	r := vec.Root()
	return vec.beautify(w, r, 0)
}
