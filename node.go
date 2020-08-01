package jsonvector

import (
	"bytes"
	"strconv"
	"unsafe"

	"github.com/koykov/fastconv"
)

type Node struct {
	t      Type
	d      int
	p      uintptr
	k, v   memseq
	cs, ce int
	Err    error
}

func (v *Node) Type() Type {
	return v.t
}

func (v *Node) Get(keys ...string) *Node {
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

func (v *Node) Len() int {
	if v.ce != v.cs && v.ce >= v.cs {
		return v.ce - v.cs
	}
	return 1
}

func (v *Node) Object() *Object {
	if v.t != TypeObj {
		return nil
	}
	return &Object{*v}
}

func (v *Node) Array() *Array {
	if v.t != TypeArr {
		return nil
	}
	return &Array{*v}
}

func (v *Node) Bytes() []byte {
	if v.t != TypeStr {
		return nil
	}
	return v.v.Bytes()
}

func (v *Node) ForceBytes() []byte {
	return v.v.Bytes()
}

func (v *Node) RawBytes() []byte {
	return v.v.rawBytes()
}

func (v *Node) String() string {
	if v.t != TypeStr {
		return ""
	}
	return v.v.String()
}

func (v *Node) Bool() bool {
	if v.t != TypeBool {
		return false
	}
	return bytes.Equal(v.v.Bytes(), bTrue)
}

func (v *Node) Float() float64 {
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

func (v *Node) Int() int64 {
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

func (v *Node) Uint() uint64 {
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

func (v *Node) childIdx() []int {
	if v.t != TypeArr && v.t != TypeObj {
		return nil
	}
	if vec := v.vec(); vec != nil {
		return vec.regGet(v.d+1, v.cs, v.ce)
	}
	return nil
}

func (v *Node) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.d, v.p = 0, 0
	v.cs, v.ce = 0, 0
	v.Err = nil
}

func (v *Node) vec() *Vector {
	if v.p == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(v.p))
}
