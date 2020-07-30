package jsonvector

type Array struct {
	Value
}

func (a *Array) At(idx int) *Value {
	vec := a.vec()
	if vec == nil {
		return nil
	}
	ci := a.childIdx()
	if idx < len(ci) {
		return &vec.v[ci[idx]]
	}
	return nil
}
