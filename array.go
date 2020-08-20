package jsonvector

type Array Node

func (a *Array) At(idx int) *Node {
	vec := (*Node)(a).vec()
	if vec == nil {
		return nil
	}
	ci := (*Node)(a).ChildIdx()
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
