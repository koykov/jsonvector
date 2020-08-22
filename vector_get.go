package jsonvector

import "github.com/koykov/bytealg"

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
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeObj {
		return nil
	}
	return v.Object()
}

func (vec *Vector) GetArrayByPath(path, sep string) *Node {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeArr {
		return nil
	}
	return v.Array()
}

func (vec *Vector) GetBytesByPath(path, sep string) []byte {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

func (vec *Vector) GetStringByPath(path, sep string) string {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

func (vec *Vector) GetBoolByPath(path, sep string) bool {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

func (vec *Vector) GetFloatByPath(path, sep string) float64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Float()
}

func (vec *Vector) GetIntByPath(path, sep string) int64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Int()
}

func (vec *Vector) GetUintByPath(path, sep string) uint64 {
	vec.ss = bytealg.AppendSplitStr(vec.ss, path, sep, -1)
	v := vec.Get(vec.ss...)
	if v == nil || v.Type() != TypeNum {
		return 0
	}
	return v.Uint()
}
