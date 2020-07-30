package jsonvector

type Object struct {
	Value
}

func (o *Object) Look(key string) *Value {
	vec := o.vec()
	if vec == nil {
		return nil
	}
	ci := o.childIdx()
	for _, i := range ci {
		c := vec.v[i]
		if key == c.k.String() {
			return &c
		}
	}
	return nil
}
