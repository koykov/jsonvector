package jsonvector

import (
	"bytes"
	"strconv"

	"github.com/koykov/fastconv"
)

type Val struct {
	t      Type
	p      *Vector
	k, v   memseq
	cs, ce int
	Err    error
}

func (v *Val) Type() Type {
	return v.t
}

func (v *Val) Get(keys ...string) *Val {
	if len(keys) == 0 {
		return v
	}
	for i := v.cs; i < v.ce; i++ {
		c := v.p.v[i]
		if bytes.Equal(c.k.Bytes(), fastconv.S2B(keys[0])) {
			if len(keys[1:]) == 0 {
				return &c
			} else {
				return c.Get(keys[1:]...)
			}
		}
	}
	return nil
}

func (v *Val) Len() int {
	if v.ce != v.cs && v.ce >= v.cs {
		return v.ce - v.cs
	}
	return 1
}

func (v *Val) Bytes() []byte {
	return v.v.Bytes()
}

func (v *Val) String() string {
	return v.v.String()
}

func (v *Val) Bool() bool {
	return bytes.Equal(v.v.Bytes(), bTrue)
}

func (v *Val) Float() float64 {
	f, err := strconv.ParseFloat(v.v.String(), 64)
	if err == nil {
		return f
	}
	v.Err = err
	return 0
}

func (v *Val) Int() int64 {
	i, err := strconv.ParseInt(v.v.String(), 10, 64)
	if err == nil {
		return i
	}
	v.Err = err
	return 0
}

func (v *Val) Uint() uint64 {
	u, err := strconv.ParseUint(v.v.String(), 10, 64)
	if err == nil {
		return u
	}
	v.Err = err
	return 0
}

func (v *Val) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.cs, v.ce = 0, 0
}
