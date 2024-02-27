package jsonvector

import (
	"github.com/koykov/vector"
	"sync"
)

// Pool represents JSON vectors pool.
type Pool struct {
	p sync.Pool
}

var (
	// P is a default instance of the pool.
	// Just call urlvector.Acquire() and urlvector.Release().
	P Pool
	// Suppress go vet warnings.
	_, _, _ = Acquire, AcquireNoClear, Release
)

// Get old vector from the pool or create new one.
func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			vec.Helper = helper
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

// Acquire returns vector from default pool instance.
func Acquire() *Vector {
	return P.Get()
}

// AcquireNoClear returns vector and skip clear step.
func AcquireNoClear() *Vector {
	vec := P.Get()
	vec.SetBit(vector.FlagNoClear, true)
	return vec
}

// Release puts vector back to default pool instance.
func Release(vec *Vector) {
	P.Put(vec)
}
