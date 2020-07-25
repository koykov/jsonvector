package jsonvector

import (
	"bytes"
	"strconv"
	"unsafe"

	"github.com/koykov/fastconv"
)

type Value struct {
	t      Type
	d      int
	p      uintptr
	k, v   memseq
	cs, ce int
	Err    error
}

func (v *Value) Type() Type {
	return v.t
}

func (v *Value) Get(keys ...string) *Value {
	if len(keys) == 0 {
		return v
	}
	if v.t != TypeObj && v.t != TypeArr {
		// Attempt to get child of scalar value.
		return nil
	}
	vec := v.vec()
	if vec == nil {
		return v
	}
	for i := v.cs; i < v.ce; i++ {
		c := vec.v[i]
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

func (v *Value) Len() int {
	if v.ce != v.cs && v.ce >= v.cs {
		return v.ce - v.cs
	}
	return 1
}

func (v *Value) Bytes() []byte {
	if v.t != TypeStr {
		return nil
	}
	return v.v.Bytes()
}

func (v *Value) ForceBytes() []byte {
	return v.v.Bytes()
}

func (v *Value) String() string {
	if v.t != TypeStr {
		return ""
	}
	return v.v.String()
}

func (v *Value) unescBytes() []byte {
	if v.t != TypeStr {
		return nil
	}
	return v.v.unescBytes()
}

func (v *Value) Bool() bool {
	if v.t != TypeBool {
		return false
	}
	return bytes.Equal(v.v.Bytes(), bTrue)
}

func (v *Value) Float() float64 {
	if v.t != TypeNum {
		return 0
	}
	f, err := strconv.ParseFloat(v.v.String(), 64)
	if err == nil {
		return f
	}
	v.Err = err
	return 0
}

func (v *Value) Int() int64 {
	if v.t != TypeNum {
		return 0
	}
	i, err := strconv.ParseInt(v.v.String(), 10, 64)
	if err == nil {
		return i
	}
	v.Err = err
	return 0
}

func (v *Value) Uint() uint64 {
	if v.t != TypeNum {
		return 0
	}
	u, err := strconv.ParseUint(v.v.String(), 10, 64)
	if err == nil {
		return u
	}
	v.Err = err
	return 0
}

func (v *Value) ChildIdx() []int {
	if v.t != TypeArr && v.t != TypeObj {
		return nil
	}
	if vec := v.vec(); vec != nil {
		return vec.regGet(v.d+1, v.cs, v.ce)
	}
	return nil
}

func (v *Value) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.cs, v.ce = 0, 0
}

func (v *Value) vec() *Vector {
	if v.p == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(v.p))
}
