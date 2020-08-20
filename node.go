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

func (n *Node) Type() Type {
	return n.t
}

func (n *Node) Get(keys ...string) *Node {
	if len(keys) == 0 {
		return n
	}
	if n.t != TypeObj && n.t != TypeArr {
		// Attempt to get child of scalar value.
		return nil
	}
	vec := n.vec()
	if vec == nil {
		return n
	}
	if n.t == TypeObj {
		for i := n.cs; i < n.ce; i++ {
			k := vec.r[n.d+1][i]
			c := &vec.v[k]
			if bytes.Equal(c.k.Bytes(), fastconv.S2B(keys[0])) {
				if len(keys[1:]) == 0 {
					return c
				} else {
					return c.Get(keys[1:]...)
				}
			}
		}
	}
	if n.t == TypeArr {
		k, err := strconv.Atoi(keys[0])
		if err != nil || k >= n.Len() {
			return nil
		}
		i := vec.r[n.d+1][n.cs+k]
		v := &vec.v[i]
		if len(keys[1:]) == 0 {
			return v
		} else {
			return n.Get(keys[1:]...)
		}
	}
	return nil
}

func (n *Node) Len() int {
	if n.ce != n.cs && n.ce >= n.cs {
		return n.ce - n.cs
	}
	return 1
}

func (n *Node) Object() *Object {
	if n.t != TypeObj {
		return nil
	}
	return (*Object)(n)
}

func (n *Node) Array() *Array {
	if n.t != TypeArr {
		return nil
	}
	return (*Array)(n)
}

func (n *Node) Bytes() []byte {
	if n.t != TypeStr {
		return nil
	}
	return n.v.Bytes()
}

func (n *Node) ForceBytes() []byte {
	return n.v.Bytes()
}

func (n *Node) RawBytes() []byte {
	return n.v.rawBytes()
}

func (n *Node) String() string {
	if n.t != TypeStr {
		return ""
	}
	return n.v.String()
}

func (n *Node) Bool() bool {
	if n.t != TypeBool {
		return false
	}
	return bytes.Equal(n.v.Bytes(), bTrue)
}

func (n *Node) Float() float64 {
	if n.t != TypeNum {
		return 0
	}
	f, err := strconv.ParseFloat(n.v.String(), 64)
	if err == nil {
		return f
	}
	n.Err = err
	return 0
}

func (n *Node) Int() int64 {
	if n.t != TypeNum {
		return 0
	}
	i, err := strconv.ParseInt(n.v.String(), 10, 64)
	if err == nil {
		return i
	}
	n.Err = err
	return 0
}

func (n *Node) Uint() uint64 {
	if n.t != TypeNum {
		return 0
	}
	u, err := strconv.ParseUint(n.v.String(), 10, 64)
	if err == nil {
		return u
	}
	n.Err = err
	return 0
}

func (n *Node) ChildIdx() []int {
	if n.t != TypeArr && n.t != TypeObj {
		return nil
	}
	if vec := n.vec(); vec != nil {
		return vec.regGet(n.d+1, n.cs, n.ce)
	}
	return nil
}

func (n *Node) Reset() {
	n.t = TypeUnk
	n.k.set(0, 0)
	n.v.set(0, 0)
	n.d, n.p = 0, 0
	n.cs, n.ce = 0, 0
	n.Err = nil
}

func (n *Node) vec() *Vector {
	if n.p == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(n.p))
}
