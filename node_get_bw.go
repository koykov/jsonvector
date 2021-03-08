package jsonvector

// Old node.Get*ByPath() methods.
// Kept for backward compatibility.

import "github.com/koykov/bytealg"

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
func (n *Node) GetFloatByPath(path, sep string) (float64, error) {
	vec := n.vec()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Float()
}

// Look and get child integer by given path.
func (n *Node) GetIntByPath(path, sep string) (int64, error) {
	vec := n.vec()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Int()
}

// Look and get child unsigned int by given path.
func (n *Node) GetUintByPath(path, sep string) (uint64, error) {
	vec := n.vec()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Uint()
}
