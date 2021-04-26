package jsonvector

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

// // Node type.
// type Type int
//
// const (
// 	TypeUnk Type = iota
// 	TypeNull
// 	TypeObj
// 	TypeArr
// 	TypeStr
// 	TypeNum
// 	TypeBool
// )

// Parser object.
type Vector struct {
	vector.Vector
	// // Source data to parse.
	// b, s []byte
	// // Source data pointer.
	// a uint64
	// // Self pointer.
	// p uintptr
	// // List of nodes.
	// v []Node
	// // Length of nodes array.
	// l int
	// // Error offset.
	// e int
	// // Nodes index.
	// i index
	// // Buffer of strings.
	// ss []string
}

var (
	// Byte constants.
	bNull  = []byte("null")
	bTrue  = []byte("true")
	bFalse = []byte("false")
	bQuote = []byte(`"`)
	bSlash = []byte(`\`)
	bFmt   = []byte(" \t\n\r")
)

// Make new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = jsonUnesc
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// Parse source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// Copy source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// Copy source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

// // Get length of nodes array.
// func (vec *Vector) Len() int {
// 	return vec.l
// }
//
// // Get error offset.
// func (vec *Vector) ErrorOffset() int {
// 	return vec.e
// }
//
// // Get root node.
// func (vec *Vector) Root() *Node {
// 	return vec.Get()
// }
//
// // Get node by given keys.
// func (vec *Vector) Get(keys ...string) *Node {
// 	if len(keys) == 0 {
// 		if vec.Len() > 0 {
// 			return &vec.v[0]
// 		}
// 		return nil
// 	}
//
// 	r := &vec.v[0]
// 	if r.t != TypeObj && r.t != TypeArr {
// 		if len(keys) > 1 {
// 			// Attempt to get child of scalar value.
// 			return nil
// 		}
// 		return r
// 	}
//
// 	if r.t == TypeArr {
// 		return vec.getArr(r, keys...)
// 	}
// 	if r.t == TypeObj {
// 		return vec.getObj(r, keys...)
// 	}
// 	return r
// }
//
// // Get node by path.
// func (vec *Vector) GetByPath(path, sep string) *Node {
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	return vec.Get(vec.ss...)
// }
//
// Format vector in human readable representation.
func (vec *Vector) Beautify(w io.Writer) error {
	r := vec.Root()
	return vec.beautify(w, r, 0)
}

//
// // Check if node exists.
// func (vec *Vector) Exists(key string) bool {
// 	n := vec.Get()
// 	return n.Exists(key)
// }
//
// // Get or create new node.
// func (vec *Vector) newNode(depth int) (r *Node) {
// 	if vec.l < len(vec.v) {
// 		r = &vec.v[vec.l]
// 		r.Reset()
// 		vec.l++
// 	} else {
// 		r = &Node{t: TypeUnk}
// 		vec.v = append(vec.v, *r)
// 		vec.l++
// 	}
// 	r.p, r.d = vec.p, depth
// 	return
// }
//
// // Reset node before put to the pool.
// func (vec *Vector) Reset() {
// 	if vec.l == 0 {
// 		return
// 	}
// 	_ = vec.v[vec.l-1]
// 	for i := 0; i < vec.l; i++ {
// 		vec.v[i].p = 0
// 	}
// 	vec.b, vec.s = vec.b[:0], nil
// 	vec.a, vec.l, vec.e = 0, 0, 0
// 	vec.i.reset()
// }
//
// // Get raw pointer of self parser.
// func (vec *Vector) ptr() uintptr {
// 	return uintptr(unsafe.Pointer(vec))
// }
