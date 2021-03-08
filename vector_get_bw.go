package jsonvector

// Old vector.Get*ByPath() methods.
// Kept for backward compatibility.

import (
	"github.com/koykov/bytealg"
)

// Look and get object by given path.
func (vec *Vector) GetObjectByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeObj {
		return nil
	}
	return v.Object()
}

// Look and get array by given path.
func (vec *Vector) GetArrayByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeArr {
		return nil
	}
	return v.Array()
}

// Look and get bytes by given path.
func (vec *Vector) GetBytesByPath(path, sep string) []byte {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

// Look and get string by given path.
func (vec *Vector) GetStringByPath(path, sep string) string {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

// Look and get bool by given path.
func (vec *Vector) GetBoolByPath(path, sep string) bool {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

// Look and get float by given path.
func (vec *Vector) GetFloatByPath(path, sep string) (float64, error) {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Float()
}

// Look and get integer by given path.
func (vec *Vector) GetIntByPath(path, sep string) (int64, error) {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Int()
}

// Look and get unsigned integer by given path.
func (vec *Vector) GetUintByPath(path, sep string) (uint64, error) {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Uint()
}
