package jsonvector

import "sync"

// Vector pool.
type Pool struct {
	p sync.Pool
}

var (
	// Default instance of the pool.
	// Just call jsonvector.Acquire() and jsonvector.Release().
	P Pool
	// Suppress go vet warnings.
	_, _ = Acquire, Release
)

// Get old vector from the pool or create new one.
func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			vec.PrepareBytesFn = PrepareBytes
			return vec
		}
	}
	return NewVector()
}

// Put vector back to the pool.
func (p *Pool) Put(vec *Vector) {
	vec.Reset()
	p.p.Put(vec)
}

// Get vector from default pool instance.
func Acquire() *Vector {
	return P.Get()
}

// Put vector back to default pool instance.
func Release(vec *Vector) {
	P.Put(vec)
}
