package jsonvector

import (
	"io"

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
	vec.SetBit(vector.FlagInit, true)
	vec.Helper = helper
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// ParseString parses source string.
func (vec *Vector) ParseString(s string) error {
	return vec.parse(byteconv.S2B(s), false)
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// ParseCopyString copies source string and parse it.
func (vec *Vector) ParseCopyString(s string) error {
	return vec.parse(byteconv.S2B(s), true)
}

// ParseFile reads file contents and parse it.
func (vec *Vector) ParseFile(path string) error {
	err := vec.Vector.ParseFile(path)
	if err != vector.ErrNotImplement {
		return err
	}
	return vec.parse(vec.Buf(), false)
}

// ParseReader reads source from r and parse it.
func (vec *Vector) ParseReader(r io.Reader) error {
	err := vec.Vector.ParseReader(r)
	if err != vector.ErrNotImplement {
		return err
	}
	return vec.parse(vec.Buf(), false)
}
