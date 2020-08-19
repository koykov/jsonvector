package jsonvector

type Array struct {
	Node
}

func (a *Array) At(idx int) *Node {
	vec := a.vec()
	if vec == nil {
		return nil
	}
	ci := a.ChildIdx()
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
