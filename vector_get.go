package jsonvector

func (vec *Vector) GetString(keys ...string) string {
	v := vec.Get(keys...)
	if v == nil || v.Type() != TypeStr {
		return ""
	}
	return v.String()
}
