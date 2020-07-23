package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/bytealg"
)

var (
	unesc       = []byte(`Lorem \"ipsum\" dolor \"sit\" amet.`)
	unescExpect = []byte(`Lorem "ipsum" dolor "sit" amet.`)

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

	obj0   = []byte(`{"a": 1, "b": 2, "c": 3}`)
	obj1   = []byte(`{"a": "foo", "b": "bar", "c": "string"}`)
	obj2   = []byte(`{"key0": "\"quoted\"", "key\"1\"": "str"}`)
	obj3   = []byte(`{"pi": 3.1415, "e": 2,718281828459045}`)
	objFmt = []byte(`{
	"c" :	15,
	"foo":null,
	"bar":  "qwerty \"encoded\""
}`)

	cmpx0       = []byte(`{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}`)
	cmpxBeauty0 = []byte(`{
	"glossary": {
		"title": "example glossary",
		"GlossDiv": {
			"title": "S",
			"GlossList": {
				"GlossEntry": {
					"ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
						"para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": [
							"GML",
							"XML"
						]
					},
					"GlossSee": "markup"
				}
			}
		}
	}
}`)
	cmpx1       = []byte(`{"firstName":"John","lastName":"Smith","isAlive":true,"age":27,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`)
	cmpxBeauty1 = []byte(`{
	"firstName": "John",
	"lastName": "Smith",
	"isAlive": true,
	"age": 27,
	"address": {
		"streetAddress": "21 2nd Street",
		"city": "New York",
		"state": "NY",
		"postalCode": "10021-3100"
	},
	"phoneNumbers": [
		{
			"type": "home",
			"number": "212 555-1234"
		},
		{
			"type": "office",
			"number": "646 555-4567"
		},
		{
			"type": "mobile",
			"number": "123 456-7890"
		}
	],
	"children": [],
	"spouse": null
}`)
	cmpx2       = []byte(`{"$schema":"http://json-schema.org/schema#","title":"Product","type":"object","required":["id","name","price"],"properties":{"id":{"type":"number","description":"Product identifier"},"name":{"type":"string","description":"Name of the product"},"price":{"type":"number","minimum":0},"tags":{"type":"array","items":{"type":"string"}},"stock":{"type":"object","properties":{"warehouse":{"type":"number"},"retail":{"type":"number"}}}}}`)
	cmpxBeauty2 = []byte(`{
	"$schema": "http://json-schema.org/schema#",
	"title": "Product",
	"type": "object",
	"required": [
		"id",
		"name",
		"price"
	],
	"properties": {
		"id": {
			"type": "number",
			"description": "Product identifier"
		},
		"name": {
			"type": "string",
			"description": "Name of the product"
		},
		"price": {
			"type": "number",
			"minimum": 0
		},
		"tags": {
			"type": "array",
			"items": {
				"type": "string"
			}
		},
		"stock": {
			"type": "object",
			"properties": {
				"warehouse": {
					"type": "number"
				},
				"retail": {
					"type": "number"
				}
			}
		}
	}
}`)
	cmpx3       = []byte(`{"id":1,"name":"Foo","price":123,"tags":["Bar","Eek"],"stock":{"warehouse":300,"retail":20}}`)
	cmpxBeauty3 = []byte(`{
	"id": 1,
	"name": "Foo",
	"price": 123,
	"tags": [
		"Bar",
		"Eek"
	],
	"stock": {
		"warehouse": 300,
		"retail": 20
	}
}`)
	cmpx4       = []byte(`{"first name":"John","last name":"Smith","age":25,"address":{"street address":"21 2nd Street","city":"New York","state":"NY","postal code":"10021"},"phone numbers":[{"type":"home","number":"212 555-1234"},{"type":"fax","number":"646 555-4567"}],"sex":{"type":"male"}}`)
	cmpxBeauty4 = []byte(`{
	"first name": "John",
	"last name": "Smith",
	"age": 25,
	"address": {
		"street address": "21 2nd Street",
		"city": "New York",
		"state": "NY",
		"postal code": "10021"
	},
	"phone numbers": [
		{
			"type": "home",
			"number": "212 555-1234"
		},
		{
			"type": "fax",
			"number": "646 555-4567"
		}
	],
	"sex": {
		"type": "male"
	}
}`)
	cmpx5       = []byte(`{"quiz":{"sport":{"q1":{"question":"Which one is correct team name in NBA?","options":["New York Bulls","Los Angeles Kings","Golden State Warriros","Huston Rocket"],"answer":"Huston Rocket"}},"maths":{"q1":{"question":"5 + 7 = ?","options":["10","11","12","13"],"answer":"12"},"q2":{"question":"12 - 8 = ?","options":["1","2","3","4"],"answer":"4"}}}}`)
	cmpxBeauty5 = []byte(`{
	"quiz": {
		"sport": {
			"q1": {
				"question": "Which one is correct team name in NBA?",
				"options": [
					"New York Bulls",
					"Los Angeles Kings",
					"Golden State Warriros",
					"Huston Rocket"
				],
				"answer": "Huston Rocket"
			}
		},
		"maths": {
			"q1": {
				"question": "5 + 7 = ?",
				"options": [
					"10",
					"11",
					"12",
					"13"
				],
				"answer": "12"
			},
			"q2": {
				"question": "12 - 8 = ?",
				"options": [
					"1",
					"2",
					"3",
					"4"
				],
				"answer": "4"
			}
		}
	}
}`)

	badTrash        = []byte(`foo bar`)
	badScalarStr    = []byte(`"unclosed string example`)
	badNumDiv       = []byte("3,14151")
	badUnparsedTail = []byte(`{"a": 1, "b": 2}foo`)

	buf []byte
	bb  bytes.Buffer
	vec = NewVector()
)

func testScalar(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(scalarNull, false)
	if vec.v[0].t != TypeNull {
		t.Error("null mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarStr, false)
	if vec.v[0].t != TypeStr || !bytes.Equal(bytealg.Trim(scalarStr, bQuote), vec.Get().Bytes()) {
		t.Error("str mismatch")
	}

	vec.Reset()
	buf = append(buf[:0], scalarStrQ...)
	_ = vec.Parse(buf, false)
	if vec.v[0].t != TypeStr || !bytes.Equal(bytealg.Trim(buf[:17], bQuote), vec.Get().Bytes()) {
		t.Error("quoted str mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum0, false)
	if vec.v[0].t != TypeNum || vec.Get().Int() != 123456 {
		t.Error("num 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum1, false)
	if vec.v[0].t != TypeNum || vec.Get().Float() != 123.456 {
		t.Error("num 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarNum2, false)
	if vec.v[0].t != TypeNum || vec.Get().Float() != 3.7e-5 {
		t.Error("num 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarTrue, false)
	if vec.v[0].t != TypeBool || vec.Get().Bool() != true {
		t.Error("bool true mismatch")
	}

	vec.Reset()
	_ = vec.Parse(scalarFalse, false)
	if vec.v[0].t != TypeBool || vec.Get().Bool() != false {
		t.Error("bool false mismatch")
	}
}

func testArr(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(arr0, false)
	v := vec.Get()
	if v.Type() != TypeArr || v.Len() != 5 || vec.Get("1").Int() != 2 {
		t.Error("arr 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr1, false)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 3 || !bytes.Equal(vec.Get("1").Bytes(), []byte("bar")) {
		t.Error("arr 1 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr2, false)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 2 || vec.Get("0").Float() != 3.14156 {
		t.Error("arr 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(arr3, false)
	v = vec.Get()
	if v.Type() != TypeArr || v.Len() != 3 || vec.Get("1").Type() != TypeNull {
		t.Error("arr 3 mismatch")
	}
}

func testObj(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(obj0, false)
	v := vec.Get()
	if v.Type() != TypeObj && v.Len() != 3 || vec.Get("b").Int() != 2 {
		t.Error("obj 0 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj1, false)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 3 || vec.Get("c").String() != "string" {
		t.Error("obj 1 mismatch")
	}

	vec.Reset()
	buf = append(buf[:0], obj2...)
	_ = vec.Parse(buf, false)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 2 || vec.Get("key0").String() != "\"quoted\"" {
		t.Error("obj 2 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(obj3, false)
	v = vec.Get()
	if v.Type() != TypeObj && v.Len() != 2 {
		t.Error("obj 3 mismatch")
	}

	vec.Reset()
	_ = vec.Parse(objFmt, false)
	v = vec.Get()
	if v.Type() != TypeObj {
		t.Error("obj fmt mismatch")
	}
}

func testComplex0(t testing.TB) {
	vec.Reset()
	_ = vec.Parse(cmpx0, false)
	v := vec.Get("glossary", "GlossDiv", "GlossList", "GlossEntry", "ID")
	if v.Type() != TypeStr && v.String() != "SGML" {
		t.Error("complex 0 mismatch ID")
	}
}

func testComplexBeauty0(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx0, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty0) {
		t.Error("complex 0 beauty mismatch")
	}
}

func testComplex1(t testing.TB) {
	var err error
	vec.Reset()
	err = vec.Parse(cmpx1, false)
	if err != nil {
		t.Error(err)
	}
	v := vec.Get("phoneNumbers", "1", "number")
	if v.Type() != TypeStr && v.String() != "646 555-4567" {
		t.Error("complex 1 mismatch phone number")
	}
}

func testComplexBeauty1(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx1, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty1) {
		t.Error("complex 1 beauty mismatch")
	}
}

func testComplex2(t testing.TB) {
	var err error
	vec.Reset()
	err = vec.Parse(cmpx2, false)
	if err != nil {
		t.Error(err)
	}
	v := vec.Get("properties", "stock", "properties", "retail", "type")
	if v.Type() != TypeStr && v.String() != "number" {
		t.Error("complex 2 mismatch property")
	}
}

func testComplexBeauty2(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx2, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty2) {
		t.Error("complex 2 beauty mismatch")
	}
}

func testComplex3(t testing.TB) {
	var err error
	vec.Reset()
	err = vec.Parse(cmpx3, false)
	if err != nil {
		t.Error(err)
	}
	v := vec.Get("stock", "retail")
	if v.Type() != TypeNum && v.Int() != 20 {
		t.Error("complex 3 mismatch stock")
	}
}

func testComplexBeauty3(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx3, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty3) {
		t.Error("complex 3 beauty mismatch")
	}
}

func testComplex4(t testing.TB) {
	var err error
	vec.Reset()
	err = vec.Parse(cmpx4, false)
	if err != nil {
		t.Error(err)
	}
	v := vec.Get("address", "postal code")
	if v.Type() != TypeStr && v.String() != "10021" {
		t.Error("complex 4 mismatch postal code")
	}
}

func testComplexBeauty4(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx4, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty4) {
		t.Error("complex 4 beauty mismatch")
	}
}

func testComplex5(t testing.TB) {
	var err error
	vec.Reset()
	err = vec.Parse(cmpx5, false)
	if err != nil {
		t.Error(err)
	}
	v := vec.Get("quiz", "maths", "q2", "options", "2")
	if v.Type() != TypeStr && v.String() != "3" {
		t.Error("complex 5 mismatch quiz option")
	}
}

func testComplexBeauty5(t testing.TB) {
	vec.Reset()
	bb.Reset()
	_ = vec.Parse(cmpx5, false)
	_ = vec.Beautify(&bb)
	if !bytes.Equal(bb.Bytes(), cmpxBeauty5) {
		t.Error("complex 5 beauty mismatch")
	}
}

func TestUnescape(t *testing.T) {
	buf = unescape(unesc)
	if !bytes.Equal(buf, unescExpect) {
		t.Error("unescape assertion failed")
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

func TestVector_ParseComplex0(t *testing.T) {
	testComplex0(t)
}

func TestVector_BeautyComplex0(t *testing.T) {
	testComplexBeauty0(t)
}

func TestVector_ParseComplex1(t *testing.T) {
	testComplex1(t)
}

func TestVector_BeautyComplex1(t *testing.T) {
	testComplexBeauty1(t)
}

func TestVector_ParseComplex2(t *testing.T) {
	testComplex2(t)
}

func TestVector_BeautyComplex2(t *testing.T) {
	testComplexBeauty2(t)
}

func TestVector_ParseComplex3(t *testing.T) {
	testComplex3(t)
}

func TestVector_BeautyComplex3(t *testing.T) {
	testComplexBeauty3(t)
}

func TestVector_ParseComplex4(t *testing.T) {
	testComplex4(t)
}

func TestVector_BeautyComplex4(t *testing.T) {
	testComplexBeauty4(t)
}

func TestVector_ParseComplex5(t *testing.T) {
	testComplex5(t)
}

func TestVector_BeautyComplex5(t *testing.T) {
	testComplexBeauty5(t)
}

func TestErr(t *testing.T) {
	var err error

	vec.Reset()
	err = vec.Parse(badTrash, false)
	if err != ErrUnexpId && vec.ErrorOffset() != 0 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badScalarStr, false)
	if err != ErrUnexpEOS || vec.ErrorOffset() != 24 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badNumDiv, false)
	if err != ErrUnparsedTail && vec.ErrorOffset() != 1 {
		t.Error("error assertion failed")
	}

	vec.Reset()
	err = vec.Parse(badUnparsedTail, false)
	if err != ErrUnparsedTail && vec.ErrorOffset() != 16 {
		t.Error("error assertion failed")
	}
}

func BenchmarkUnescape(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf = append(buf[:0], unesc...)
		buf = unescape(buf)
		if !bytes.Equal(buf, unescExpect) {
			b.Error("unescape assertion failed")
		}
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

func BenchmarkVector_ParseComplex0(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex0(b)
	}
}

func BenchmarkVector_ParseComplex1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex1(b)
	}
}

func BenchmarkVector_ParseComplex2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex2(b)
	}
}

func BenchmarkVector_ParseComplex3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex3(b)
	}
}

func BenchmarkVector_ParseComplex4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex4(b)
	}
}

func BenchmarkVector_ParseComplex5(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex5(b)
	}
}
