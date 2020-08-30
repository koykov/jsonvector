// Getter API of JSON node.
package jsonvector

import "github.com/koykov/bytealg"

// Look and get child object by given keys.
func (n *Node) GetObject(keys ...string) *Node {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeObj {
		return nil
	}
	return c.Object()
}

// Look and get child array by given keys.
func (n *Node) GetArray(keys ...string) *Node {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeArr {
		return nil
	}
	return c.Array()
}

// Look and get child bytes by given keys.
func (n *Node) GetBytes(keys ...string) []byte {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

// Look and get child string by given keys.
func (n *Node) GetString(keys ...string) string {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

// Look and get child bool by given keys.
func (n *Node) GetBool(keys ...string) bool {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

// Look and get child float by given keys.
func (n *Node) GetFloat(keys ...string) float64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Float()
}

// Look and get child integer by given keys.
func (n *Node) GetInt(keys ...string) int64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Int()
}

// Look and get child unsigned integer by given keys.
func (n *Node) GetUint(keys ...string) uint64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Uint()
}

// Look and get child object by given path.
func (n *Node) GetObjectByPath(path, sep string) *Node {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeObj {
		return nil
	}
	return c.Object()
}

// Look and get child array by given path.
func (n *Node) GetArrayByPath(path, sep string) *Node {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeArr {
		return nil
	}
	return c.Array()
}

// Look and get child bytes by given path.
func (n *Node) GetBytesByPath(path, sep string) []byte {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

// Look and get child string by given path.
func (n *Node) GetStringByPath(path, sep string) string {
	vec := n.vec()
	if vec == nil {
		return ""
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

// Look and get child bool by given path.
func (n *Node) GetBoolByPath(path, sep string) bool {
	vec := n.vec()
	if vec == nil {
		return false
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

// Look and get child float by given path.
func (n *Node) GetFloatByPath(path, sep string) float64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Float()
}

// Look and get child integer by given path.
func (n *Node) GetIntByPath(path, sep string) int64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Int()
}

// Look and get child unsigned int by given path.
func (n *Node) GetUintByPath(path, sep string) uint64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Uint()
}
