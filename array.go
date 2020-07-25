package jsonvector

type Array []Value

func (a *Array) Get(idx int) *Value {
	if idx < len(*a) {
		return &((*a)[idx])
	}
	return nil
}
