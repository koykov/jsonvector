package jsonvector

func (n *Node) GetObject(keys ...string) *Node {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeObj {
		return nil
	}
	return c.Object()
}

func (n *Node) GetArray(keys ...string) *Node {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeArr {
		return nil
	}
	return c.Array()
}

func (n *Node) GetBytes(keys ...string) []byte {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

func (n *Node) GetString(keys ...string) string {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

func (n *Node) GetBool(keys ...string) bool {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

func (n *Node) GetFloat(keys ...string) float64 {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Float()
}

func (n *Node) GetInt(keys ...string) int64 {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Int()
}

func (n *Node) GetUint(keys ...string) uint64 {
	c := n.vec().Get(keys...)
	if c == nil || c.Type() != TypeNum {
		return 0
	}
	return c.Uint()
}
