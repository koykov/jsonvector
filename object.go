package jsonvector

type Object Value

func (o *Object) Get(key string) *Value {
	v := (*Value)(o)
	vec := v.vec()
	if vec == nil {
		return nil
	}
	ci := v.childIdx()
	for _, i := range ci {
		c := vec.v[i]
		if key == c.k.String() {
			return &c
		}
	}
	return nil
}

func (o *Object) Len() int {
	return (*Value)(o).Len()
}
