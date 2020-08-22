package jsonvector

import "github.com/koykov/bytealg"

func (n *Node) GetObject(keys ...string) *Node {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeObj {
		return nil
	}
	return c.Object()
}

func (n *Node) GetArray(keys ...string) *Node {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeArr {
		return nil
	}
	return c.Array()
}

func (n *Node) GetBytes(keys ...string) []byte {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

func (n *Node) GetString(keys ...string) string {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

func (n *Node) GetBool(keys ...string) bool {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

func (n *Node) GetFloat(keys ...string) float64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Float()
}

func (n *Node) GetInt(keys ...string) int64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Int()
}

func (n *Node) GetUint(keys ...string) uint64 {
	c := n.Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Uint()
}

func (n *Node) GetObjectByPath(path, sep string) *Node {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeObj {
		return nil
	}
	return c.Object()
}

func (n *Node) GetArrayByPath(path, sep string) *Node {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeArr {
		return nil
	}
	return c.Array()
}

func (n *Node) GetBytesByPath(path, sep string) []byte {
	vec := n.vec()
	if vec == nil {
		return nil
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

func (n *Node) GetStringByPath(path, sep string) string {
	vec := n.vec()
	if vec == nil {
		return ""
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

func (n *Node) GetBoolByPath(path, sep string) bool {
	vec := n.vec()
	if vec == nil {
		return false
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

func (n *Node) GetFloatByPath(path, sep string) float64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Float()
}

func (n *Node) GetIntByPath(path, sep string) int64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Int()
}

func (n *Node) GetUintByPath(path, sep string) uint64 {
	vec := n.vec()
	if vec == nil {
		return 0
	}
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	c := n.Get(vec.ss...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Uint()
}
