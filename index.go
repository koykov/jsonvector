package jsonvector

// Nodes index.
//
// Contains indexes of nodes in the vector divided by depth.
// Y-axis means depth, X-axis means position in index.
type index struct {
	// Index tree.
	t [][]int
	// Index depth.
	d int
}

// Register new index for given depth.
func (r *index) reg(depth, idx int) int {
	if len(r.t) <= depth {
		for len(r.t) <= depth {
			r.t = append(r.t, nil)
			r.d = len(r.t)
		}
	}
	r.t[depth] = append(r.t[depth], idx)
	return len(r.t[depth])
}

// Get count of indexes registered on depth.
func (r *index) len(depth int) int {
	if len(r.t) <= depth {
		return 0
	}
	return len(r.t[depth])
}

// Get subset [s:e] of indexes registered on depth.
func (r *index) get(depth, s, e int) []int {
	l := r.len(depth)
	if l > s {
		return r.t[depth][s:e]
	}
	return nil
}

// Reset the index.
func (r *index) reset() {
	for i := 0; i < len(r.t); i++ {
		r.t[i] = r.t[i][:0]
	}
	r.d = 0
}
