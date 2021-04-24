package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
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

	obj0   = []byte(`{"a": 1, "b": 2, "c": 3}`)
	obj1   = []byte(`{"a": "foo", "b": "bar", "c": "string"}`)
	obj2   = []byte(`{"key0": "\"quoted\"", "key\"1\"": "str"}`)
	obj3   = []byte(`{"pi": 3.1415, "e": 2,718281828459045}`)
	objFmt = []byte(`{
	"c" :	15,
	"foo":null,
	"bar":  "qwerty \"encoded\""
}`)
	objFmt1 = []byte(`{
  "a": true,
  "b": {
    "c": "foo",
    "d": [
      5,
      3.1415,
      812.48927
    ]
  }
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
	if vec.Root().Type() != vector.TypeNull {
		t.Error("null mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStr)
	if vec.Root().Type() != vector.TypeStr || !bytes.Equal(bytealg.Trim(scalarStr, bQuote), vec.Get().Bytes()) {
		t.Error("str mismatch")
	}

	vec.Reset()
	_ = vec.ParseCopy(scalarStrQ)
	if vec.Root().Type() != vector.TypeStr || !bytes.Equal(scalarStrQexp, vec.Get().Bytes()) {
		t.Error("quoted str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum0)
	i, _ := vec.Get().Int()
	if vec.Root().Type() != vector.TypeNum || i != 123456 {
		t.Error("num 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum1)
	f, _ := vec.Get().Float()
	if vec.Root().Type() != vector.TypeNum || f != 123.456 {
		t.Error("num 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum2)
	f, _ = vec.Get().Float()
	if vec.Root().Type() != vector.TypeNum || f != 3.7e-5 {
		t.Error("num 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarTrue)
	if vec.Root().Type() != vector.TypeBool || vec.Get().Bool() != true {
		t.Error("bool true mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarFalse)
	if vec.Root().Type() != vector.TypeBool || vec.Get().Bool() != false {
		t.Error("bool false mismatch")
	}
}

func testArr(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(arr0)
	v := vec.Get()
	i, _ := vec.Get("1").Int()
	if v.Type() != vector.TypeArr || v.Limit() != 5 || i != 2 {
		t.Error("arr 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr1)
	v = vec.Get()
	if v.Type() != vector.TypeArr || v.Limit() != 3 || !bytes.Equal(vec.Get("1").Bytes(), []byte("bar")) {
		t.Error("arr 1 mismatch")
	}
	v.Each(func(idx int, node *vector.Node) {
		switch idx {
		case 0:
			if v := node.String(); v != "foo" {
				t.Error(`arr 1 val 0 mismatch, need "foo", got`, v)
			}
		case 1:
			if v := node.String(); v != "bar" {
				t.Error(`arr 1 val 1 mismatch, need "bar", got`, v)
			}
		case 2:
			if v := node.String(); v != "string" {
				t.Error(`arr 1 val 2 mismatch, need "string", got`, v)
			}
		}
	})

	vec.Reset()
	_ = vec.Parse(arr2)
	v = vec.Get()
	f, _ := vec.Get("0").Float()
	if v.Type() != vector.TypeArr || v.Limit() != 2 || f != 3.14156 {
		t.Error("arr 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr3)
	v = vec.Get()
	if v.Type() != vector.TypeArr || v.Limit() != 3 || vec.Get("1").Type() != vector.TypeNull {
		t.Error("arr 3 mismatch")
	}
}

func testObj(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(obj0)
	v := vec.Get()
	i, _ := vec.Get("b").Int()
	if v.Type() != vector.TypeObj && v.Limit() != 3 || i != 2 {
		t.Error("obj 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj1)
	v = vec.Get()
	if v.Type() != vector.TypeObj && v.Limit() != 3 || vec.Get("c").String() != "string" {
		t.Error("obj 1 mismatch")
	}

	vec.Reset()
	_ = vec.ParseCopy(obj2)
	v = vec.Get()
	if v.Type() != vector.TypeObj && v.Limit() != 2 || vec.Get("key0").String() != "\"quoted\"" {
		t.Error("obj 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj3)
	v = vec.Get()
	if v.Type() != vector.TypeObj && v.Limit() != 2 {
		t.Error("obj 3 mismatch")
	}
	v.Each(func(idx int, node *vector.Node) {
		switch idx {
		case 0:
			if k := node.KeyString(); k != "pi" {
				t.Error(`obj 3 key 0 mismatch, need "pi", got`, k)
			}
			if v, _ := node.Float(); v != 3.1415 {
				t.Error(`obj 3 value 0 mismatch, need 3.1415, got`, v)
			}
		case 1:
			if k := node.KeyString(); k != "e" {
				t.Error(`obj 3 key 0 mismatch, need "e", got`, k)
			}
			if v, _ := node.Float(); v != 2 {
				t.Error(`obj 3 value 0 mismatch, need 2, got`, v)
			}
		}
	})

	vec.Reset()
	_ = vec.Parse(objFmt)
	v = vec.Get()
	if v.Type() != vector.TypeObj {
		t.Error("obj fmt mismatch")
	}
}

func testFmt(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(objFmt1)
	v := vec.Get()
	if v.Type() != vector.TypeObj {
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

func TestVector_ParseFmt(t *testing.T) {
	testFmt(t)
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

func BenchmarkVector_ParseFmt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFmt(b)
	}
}
