package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/bytealg"
)

var (
	scalarNull  = []byte("null")
	scalarStr   = []byte(`"foo bar string"`)
	scalarStrQ  = []byte(`"foo \"bar\" string"`)
	scalarNum0  = []byte("123456")
	scalarNum1  = []byte("123.456")
	scalarNum2  = []byte("3.7e-5")
	scalarTrue  = []byte("true")
	scalarFalse = []byte("false")

	arr0 = []byte(`[1, 2, 3, 4, 5]`)
	arr1 = []byte(`["foo", "bar", "string"]`)
	arr2 = []byte(`[3.14156, 6.23e-4]`)
	arr3 = []byte(`["quoted \"str\" value", null, "foo"]`)
)

var (
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
	if vec.v[0].t != TypeStr || !bytes.Equal(bytealg.Trim(scalarStr, bQuote), vec.v[0].v.Bytes()) {
		t.Error("str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStrQ)
	if vec.v[0].t != TypeStr || !bytes.Equal(bytealg.Trim(scalarStrQ, bQuote), vec.v[0].v.Bytes()) {
		t.Error("quoted str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum0)
	if vec.v[0].t != TypeNum || !bytes.Equal(scalarNum0, vec.v[0].v.Bytes()) {
		t.Error("num 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum1)
	if vec.v[0].t != TypeNum || !bytes.Equal(scalarNum1, vec.v[0].v.Bytes()) {
		t.Error("num 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum2)
	if vec.v[0].t != TypeNum || !bytes.Equal(scalarNum2, vec.v[0].v.Bytes()) {
		t.Error("num 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarTrue)
	if vec.v[0].t != TypeBool || !bytes.Equal(bTrue, vec.v[0].v.Bytes()) {
		t.Error("bool true mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarFalse)
	if vec.v[0].t != TypeBool || !bytes.Equal(bFalse, vec.v[0].v.Bytes()) {
		t.Error("bool false mismatch")
	}
}

func testArr(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(arr0)
	if vec.v[0].t != TypeArr || vec.l != 6 || len(vec.v[0].r) != 5 {
		t.Error("arr 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr1)
	if vec.v[0].t != TypeArr || vec.l != 3 || len(vec.v[0].r) != 2 {
		t.Error("arr 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr2)
	if vec.v[0].t != TypeArr || vec.l != 2 || len(vec.v[0].r) != 1 {
		t.Error("arr 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr3)
	if vec.v[0].t != TypeArr || vec.l != 3 || len(vec.v[0].r) != 2 {
		t.Error("arr 3 mismatch")
	}
}

func TestVector_ParseScalar(t *testing.T) {
	testScalar(t)
}

func BenchmarkVector_ParseScalar(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testScalar(b)
	}
}

func TestVector_ParseArr(t *testing.T) {
	testArr(t)
}
