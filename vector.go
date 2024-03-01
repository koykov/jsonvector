package jsonvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

// Vector represents URL parser.
type Vector struct {
	vector.Vector
}

// NewVector makes new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = helper
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// ParseStr parses source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(byteconv.S2B(s), false)
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// ParseCopyStr copies source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(byteconv.S2B(s), true)
}
