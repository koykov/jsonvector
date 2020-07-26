package jsonvector

type Array Value

func (a *Array) Get(idx int) *Value {
	v := (*Value)(a)
	vec := v.vec()
	if vec == nil {
		return nil
	}
	ci := v.childIdx()
	if idx < len(ci) {
		return &vec.v[ci[idx]]
	}
	return nil
}

func (a *Array) Len() int {
	return (*Value)(a).Len()
}
