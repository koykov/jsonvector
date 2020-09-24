package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/bytealg"
)

var (
	scalarNull    = []byte("null")
	scalarStr     = []byte(`"foo bar string"`)
	scalarStrQ    = []byte(`"foo \"bar\" string"`)
	scalarStrQexp = []byte(`foo "bar" string`)
	scalarNum0    = []byte("123456")
	scalarNum1    = []byte("123.456")
	scalarNum2    = []byte("3.7e-5")
	scalarTrue    = []byte("true")
	scalarFalse   = []byte("false")

	arr0 = []byte(`[1, 2, 3, 4, 5]`)
	arr1 = []byte(`["foo", "bar", "string"]`)
	arr2 = []byte(`[3.14156, 6.23e-4]`)
	arr3 = []byte(`["quoted \"str\" value", null, "foo"]`)

	obj0    = []byte(`{"a": 1, "b": 2, "c": 3}`)
	obj1    = []byte(`{"a": "foo", "b": "bar", "c": "string"}`)
	obj2    = []byte(`{"key0": "\"quoted\"", "key\"1\"": "str"}`)
	obj2exp = []byte(`{"key0": "\"quoted\"", "key\"1\"": "str"}`)
	obj3    = []byte(`{"pi": 3.1415, "e": 2,718281828459045}`)
	objFmt  = []byte(`{
	"c" :	15,
	"foo":null,
	"bar":  "qwerty \"encoded\""
}`)

	badTrash        = []byte(`foo bar`)
	badScalarStr    = []byte(`"unclosed string example`)
	badNumDiv       = []byte("3,14151")
	badUnparsedTail = []byte(`{"a": 1, "b": 2}foo`)

	vec = NewVector()
)

func testScalar(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(scalarNull)
	if vec.v[0].t != TypeNull {
		t.Error("null mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStr)
	if vec.v[0].t != TypeStr || !bytes.Equal(bytealg.Trim(scalarStr, bQuote), vec.Get().Bytes()) {
		t.Error("str mismatch")
	}

	vec.Reset()
	_ = vec.ParseCopy(scalarStrQ)
	if vec.v[0].t != TypeStr || !bytes.Equal(scalarStrQexp, vec.Get().Bytes()) {
		t.Error("quoted str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum0)
	if vec.v[0].t != TypeNum || vec.Get().Int() != 123456 {
		t.Error("num 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum1)
	if vec.v[0].t != TypeNum || vec.Get().Float() != 123.456 {
		t.Error("num 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum2)
	if vec.v[0].t != TypeNum || vec.Get().Float() != 3.7e-5 {
		t.Error("num 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarTrue)
	if vec.v[0].t != TypeBool || vec.Get().Bool() != true {
		t.Error("bool true mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarFalse)
	if vec.v[0].t != TypeBool || vec.Get().Bool() != false {
		t.Error("bool false mismatch")
	}
}

func testArr(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(arr0)
	v := vec.Get()
	if v.Type() != TypeArr || v.Len() != 5 || vec.Get("1").Int() != 2 {
		t.Error("arr 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr1)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 3 || !bytes.Equal(vec.Get("1").Bytes(), []byte("bar")) {
		t.Error("arr 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr2)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 2 || vec.Get("0").Float() != 3.14156 {
		t.Error("arr 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr3)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 3 || vec.Get("1").Type() != TypeNull {
		t.Error("arr 3 mismatch")
	}
}

func testObj(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(obj0)
	v := vec.Get()
	if v.Type() != TypeObj && v.Len() != 3 || vec.Get("b").Int() != 2 {
		t.Error("obj 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj1)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 3 || vec.Get("c").String() != "string" {
		t.Error("obj 1 mismatch")
	}

	vec.Reset()
	_ = vec.ParseCopy(obj2)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 2 || vec.Get("key0").String() != "\"quoted\"" {
		t.Error("obj 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj3)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 2 {
		t.Error("obj 3 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(objFmt)
	v = vec.Get()
	if v.Type() != TypeObj {
		t.Error("obj fmt mismatch")
	}
}

func TestVector_ParseScalar(t *testing.T) {
	testScalar(t)
}

func TestVector_ParseArr(t *testing.T) {
	testArr(t)
}

func TestVector_ParseObj(t *testing.T) {
	testObj(t)
}

func TestErr(t *testing.T) {
	var err error

	vec.Reset()
	err = vec.Parse(badTrash)
	if err != ErrUnexpId && vec.ErrorOffset() != 0 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badScalarStr)
	if err != ErrUnexpEOS || vec.ErrorOffset() != 24 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badNumDiv)
	if err != ErrUnparsedTail && vec.ErrorOffset() != 1 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badUnparsedTail)
	if err != ErrUnparsedTail && vec.ErrorOffset() != 16 {
		t.Error("error assertion failed")
	}
}

func BenchmarkVector_ParseScalar(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testScalar(b)
	}
}

func BenchmarkVector_ParseArr(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testArr(b)
	}
}

func BenchmarkVector_ParseObj(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testObj(b)
	}
}
