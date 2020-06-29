package jsonvector

type Val struct {
	t    Type
	k, v memseq
	r    []int
}

func (v *Val) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.r = v.r[:0]
}
