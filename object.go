package jsonvector

import "bytes"

type Object []Value

func (o *Object) Look(key []byte) *Value {
	for _, v := range *o {
		if bytes.Equal(key, v.k.Bytes()) {
			return &v
		}
	}
	return nil
}
