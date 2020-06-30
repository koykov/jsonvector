package jsonvector

type Val struct {
	t      Type
	k, v   memseq
	cs, ce int
}

func (v *Val) Reset() {
	v.t = TypeUnk
	v.k.set(0, 0)
	v.v.set(0, 0)
	v.cs, v.ce = 0, 0
}
