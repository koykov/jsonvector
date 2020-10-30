package jsonvector

import (
	"bytes"
	"strconv"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

// Json node object.
type Node struct {
	// Node type.
	t Type
	// Node depth in json object.
	d int
	// Raw pointer to vector.
	// It's safe to using uintptr here because vector guaranteed to exist while the node is alive and isn't a garbage
	// collected.
	p uintptr
	// Key/value bytes
	k, v byteptr
	// Position of first and last child nodes in index.
	s, e int
}

// Get node type.
func (n *Node) Type() Type {
	return n.t
}

// Get child node by given keys.
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
		for i := n.s; i < n.e; i++ {
			k := vec.i.t[n.d+1][i]
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
		i := vec.i.t[n.d+1][n.s+k]
		v := &vec.v[i]
		if len(keys[1:]) == 0 {
			return v
		} else {
			return n.Get(keys[1:]...)
		}
	}
	return nil
}

// Get child node by path.
func (n *Node) GetByPath(path, sep string) *Node {
	vec := n.vec()
	if vec == nil {
		return n
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	return n.Get(vec.ss...)
}

// Check if key exists in child nodes.
func (n *Node) Exists(key string) bool {
	if n.t != TypeObj {
		return false
	}
	vec := n.vec()
	if vec == nil {
		return false
	}
	for i := n.s; i < n.e; i++ {
		k := vec.i.t[n.d+1][i]
		c := &vec.v[k]
		if c.k.String() == key {
			return true
		}
	}
	return false
}

// Get length of child nodes.
func (n *Node) Len() int {
	if n.e != n.s && n.e >= n.s {
		return n.e - n.s
	}
	return 1
}

// Convert current node to object.
func (n *Node) Object() *Node {
	if n.t != TypeObj {
		return nil
	}
	return n
}

// Convert current node to array.
func (n *Node) Array() *Node {
	if n.t != TypeArr {
		return nil
	}
	return n
}

func (n *Node) Key() []byte {
	if n.k.o != 0 && n.k.l > 0 {
		return n.k.Bytes()
	}
	return nil
}

func (n *Node) KeyString() string {
	if n.k.o != 0 && n.k.l > 0 {
		return n.k.String()
	}
	return ""
}

// Get node value as bytes.
func (n *Node) Bytes() []byte {
	if n.t != TypeStr && n.t != TypeNum {
		return nil
	}
	return n.v.Bytes()
}

// Get node value as bytes even if type isn't a string.
func (n *Node) ForceBytes() []byte {
	return n.v.Bytes()
}

// Get node value as bytes without unescape.
func (n *Node) RawBytes() []byte {
	return n.v.rawBytes()
}

// Get node value string.
func (n *Node) String() string {
	if n.t != TypeStr && n.t != TypeNum {
		return ""
	}
	return n.v.String()
}

// Get node value as boolean.
func (n *Node) Bool() bool {
	if n.t != TypeBool {
		return false
	}
	return bytes.Equal(bytealg.ToLower(n.v.Bytes()), bTrue)
}

// Get node value as float.
func (n *Node) Float() (float64, error) {
	if n.t != TypeNum {
		return 0, ErrIncompatType
	}
	f, err := strconv.ParseFloat(n.v.String(), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// Get node value as integer.
func (n *Node) Int() (int64, error) {
	if n.t != TypeNum {
		return 0, ErrIncompatType
	}
	i, err := strconv.ParseInt(n.v.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Get node value as unsigned integer.
func (n *Node) Uint() (uint64, error) {
	if n.t != TypeNum {
		return 0, ErrIncompatType
	}
	u, err := strconv.ParseUint(n.v.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return u, nil
}

func (n *Node) Each(fn func(idx int, node *Node)) {
	idx := n.childIdx()
	vec := n.vec()
	if len(idx) == 0 || vec == nil {
		return
	}
	c := 0
	for _, i := range idx {
		cn := &vec.v[i]
		fn(c, cn)
		c++
	}
}

// Get indexes of child nodes.
func (n *Node) childIdx() []int {
	if vec := n.vec(); vec != nil {
		var e = n.e
		if e == 0 {
			e = n.s + 1
		}
		return vec.i.get(n.d+1, n.s, e)
	}
	return nil
}

// Reset the node.
func (n *Node) Reset() {
	n.t = TypeUnk
	n.k.set(0, 0)
	n.v.set(0, 0)
	n.d, n.p = 0, 0
	n.s, n.e = 0, 0
}

// Restore entire parser object from raw pointer.
func (n *Node) vec() *Vector {
	if n.p == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(n.p))
}
