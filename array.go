package jsonvector

func (n *Node) At(idx int) *Node {
	if n.t != TypeArr {
		return nil
	}
	vec := n.vec()
	if vec == nil {
		return nil
	}
	ci := n.ChildIdx()
	h := -1
	for _, i := range ci {
		if i == idx {
			h = i
			break
		}
	}
	if h >= 0 {
		return &vec.v[h]
	}
	return nil
}
