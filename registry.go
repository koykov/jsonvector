package jsonvector

// Index registry storage.
//
// Contains indexes of nodes in the vector divided by depth.
type registry struct {
	r [][]int
	l int
}

// Register new index in registry for given depth.
func (r *registry) reg(depth, idx int) int {
	if len(r.r) <= depth {
		for len(r.r) <= depth {
			r.r = append(r.r, nil)
			r.l = len(r.r)
		}
	}
	r.r[depth] = append(r.r[depth], idx)
	return len(r.r[depth])
}

// Get count of indexes registered on depth.
func (r *registry) len(depth int) int {
	if len(r.r) <= depth {
		return 0
	}
	return len(r.r[depth])
}

// Get subset [s:e] of indexes registered on depth.
func (r *registry) get(depth, s, e int) []int {
	l := r.len(depth)
	if l > s {
		return r.r[depth][s:e]
	}
	return nil
}

// Reset the registry.
func (r *registry) reset() {
	for i := 0; i < len(r.r); i++ {
		r.r[i] = r.r[i][:0]
	}
	r.l = 0
}
