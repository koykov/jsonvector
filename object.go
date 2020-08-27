package jsonvector

func (n *Node) Look(key string) *Node {
	if n.t != TypeObj {
		return nil
	}
	vec := n.vec()
	if vec == nil {
		return nil
	}
	ci := n.ChildIdx()
	for _, i := range ci {
		c := vec.v[i]
		if key == c.k.String() {
			return &c
		}
	}
	return nil
}
