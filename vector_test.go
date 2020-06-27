package jsonvector

import "testing"

var (
	scalarNull  = []byte("null")
	scalarStr   = []byte(`"foo bar string"`)
	scalarStrQ  = []byte(`"foo \"bar\" string"`)
	scalarNum0  = []byte("123456")
	scalarNum1  = []byte("123.456")
	scalarNum2  = []byte("3.7e-5")
	scalarTrue  = []byte("true")
	scalarFalse = []byte("false")
)

func TestVector_ParseScalar(t *testing.T) {
	vec := NewVector()
	_ = vec.Parse(scalarNull)
	if vec.v[0].t != TypeNull {
		t.Error("null mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStr)
	if vec.v[0].t != TypeStr {
		t.Error("str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStrQ)
	if vec.v[0].t != TypeStr {
		t.Error("quoted str mismatch")
	}
}
