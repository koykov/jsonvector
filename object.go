package jsonvector

type Object struct {
	Node
}

func (o *Object) Look(key string) *Node {
	vec := o.vec()
	if vec == nil {
		return nil
	}
	ci := o.ChildIdx()
	for _, i := range ci {
		c := vec.v[i]
		if key == c.k.String() {
			return &c
		}
	}
	return nil
}
