package jsonvector

import (
	"bytes"

	"github.com/koykov/fastconv"
)

type Val struct {
	t      Type
	k, v   memseq
	cs, ce int
}

func (v *Val) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.cs, v.ce = 0, 0
}

func (v *Val) Get(vec *Vector, keys ...string) *Val {
	if len(keys) == 0 {
		return v
	}
	for i := v.cs; i < v.ce; i++ {
		c := vec.v[i]
		if bytes.Equal(c.k.Bytes(), fastconv.S2B(keys[0])) {
			if len(keys[1:]) == 0 {
				return &c
			} else {
				return c.Get(vec, keys[1:]...)
			}
		}
	}
	return nil
}
