package jsonvector

import "bytes"

type Object []Val

func (o *Object) Look(key []byte) *Val {
	for _, v := range *o {
		if bytes.Equal(key, v.k.Bytes()) {
			return &v
		}
	}
	return nil
}
