package jsonvector

type Array []Val

func (a *Array) Get(idx int) *Val {
	if idx < len(*a) {
		return &((*a)[idx])
	}
	return nil
}
