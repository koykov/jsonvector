package jsonvector

import (
	"bytes"
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

func (vec *Vector) GetObject(keys ...string) *Node {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeObj {
		return nil
	}
	return v.Object()
}

func (vec *Vector) GetArray(keys ...string) *Node {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeArr {
		return nil
	}
	return v.Array()
}

func (vec *Vector) GetBytes(keys ...string) []byte {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

func (vec *Vector) GetString(keys ...string) string {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

func (vec *Vector) GetBool(keys ...string) bool {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

func (vec *Vector) GetFloat(keys ...string) float64 {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Float()
}

func (vec *Vector) GetInt(keys ...string) int64 {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Int()
}

func (vec *Vector) GetUint(keys ...string) uint64 {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Uint()
}

func (vec *Vector) GetObjectByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeObj {
		return nil
	}
	return v.Object()
}

func (vec *Vector) GetArrayByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeArr {
		return nil
	}
	return v.Array()
}

func (vec *Vector) GetBytesByPath(path, sep string) []byte {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

func (vec *Vector) GetStringByPath(path, sep string) string {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

func (vec *Vector) GetBoolByPath(path, sep string) bool {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

func (vec *Vector) GetFloatByPath(path, sep string) float64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Float()
}

func (vec *Vector) GetIntByPath(path, sep string) int64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Int()
}

func (vec *Vector) GetUintByPath(path, sep string) uint64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Uint()
}

func (vec *Vector) getArr(root *Node, keys ...string) *Node {
	if len(keys) == 0 {
		return root
	}
	k, err := strconv.Atoi(keys[0])
	if err != nil || k >= root.Len() {
		return nil
	}
	i := vec.r.r[root.d+1][root.s+k]
	v := &vec.v[i]
	tail := keys[1:]
	if v.t != TypeArr && v.t != TypeObj {
		if len(tail) > 0 {
			// Attempt to get child of scalar value.
			return nil
		}
		return v
	}
	if v.t == TypeArr {
		return vec.getArr(v, tail...)
	}
	if v.t == TypeObj {
		return vec.getObj(v, tail...)
	}
	return nil
}

func (vec *Vector) getObj(root *Node, keys ...string) *Node {
	if len(keys) == 0 {
		return root
	}
	var v *Node
	for i := root.s; i < root.e; i++ {
		k := vec.r.r[root.d+1][i]
		v = &vec.v[k]
		if bytes.Equal(v.k.Bytes(), fastconv.S2B(keys[0])) {
			break
		}
	}
	if v == nil {
		return v
	}
	tail := keys[1:]
	if v.t != TypeArr && v.t != TypeObj {
		if len(tail) > 0 {
			// Attempt to get child of scalar value.
			return nil
		}
		return v
	}
	if v.t == TypeArr {
		return vec.getArr(v, tail...)
	}
	if v.t == TypeObj {
		return vec.getObj(v, tail...)
	}
	return nil
}
