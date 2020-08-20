package jsonvector

func (vec *Vector) GetObject(keys ...string) *Object {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeObj {
		return nil
	}
	return v.Object()
}

func (vec *Vector) GetArray(keys ...string) *Array {
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
