package jsonvector

type Object Node

func (o *Object) Look(key string) *Node {
	vec := (*Node)(o).vec()
	if vec == nil {
		return nil
	}
	ci := (*Node)(o).ChildIdx()
	for _, i := range ci {
		c := vec.v[i]
		if key == c.k.String() {
			return &c
		}
	}
	return nil
}
