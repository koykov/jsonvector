package jsonvector

import (
	"bytes"
	"testing"
)

var (
	complex0     = []byte(`{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}`)
	complexPath0 = []string{"glossary", "GlossDiv", "GlossList", "GlossEntry", "ID"}
	complexFmt0  = []byte(`{
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
	complex1     = []byte(`{"firstName":"John","lastName":"Smith","isAlive":true,"age":27,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`)
	complexPath1 = []string{"phoneNumbers", "1", "number"}
	complexFmt1  = []byte(`{
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
	complex2     = []byte(`{"$schema":"http://json-schema.org/schema#","title":"Product","type":"object","required":["id","name","price"],"properties":{"id":{"type":"number","description":"Product identifier"},"name":{"type":"string","description":"Name of the product"},"price":{"type":"number","minimum":0},"tags":{"type":"array","items":{"type":"string"}},"stock":{"type":"object","properties":{"warehouse":{"type":"number"},"retail":{"type":"number"}}}}}`)
	complexPath2 = []string{"properties", "stock", "properties", "retail", "type"}
	complexFmt2  = []byte(`{
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
	complex3     = []byte(`{"id":1,"name":"Foo","price":123,"tags":["Bar","Eek"],"stock":{"warehouse":300,"retail":20}}`)
	complexPath3 = []string{"stock", "retail"}
	complexFmt3  = []byte(`{
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
	complex4     = []byte(`{"first name":"John","last name":"Smith","age":25,"address":{"street address":"21 2nd Street","city":"New York","state":"NY","postal code":"10021"},"phone numbers":[{"type":"home","number":"212 555-1234"},{"type":"fax","number":"646 555-4567"}],"sex":{"type":"male"}}`)
	complexPath4 = []string{"address", "postal code"}
	complexFmt4  = []byte(`{
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
	complex5     = []byte(`{"quiz":{"sport":{"q1":{"question":"Which one is correct team name in NBA?","options":["New York Bulls","Los Angeles Kings","Golden State Warriros","Huston Rocket"],"answer":"Huston Rocket"}},"maths":{"q1":{"question":"5 + 7 = ?","options":["10","11","12","13"],"answer":"12"},"q2":{"question":"12 - 8 = ?","options":["1","2","3","4"],"answer":"4"}}}}`)
	complexPath5 = []string{"quiz", "maths", "q2", "options", "2"}
	complexFmt5  = []byte(`{
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

	bbuf bytes.Buffer
)

func testComplex(t testing.TB, key string, src []byte, path []string, typ Type, val interface{}) {
	vec.Reset()
	err := vec.Parse(src, false)
	if err != nil {
		t.Error(key, err)
	}
	v := vec.Get(path...)
	if v.Type() != typ {
		t.Error(key, "type assertion failed:", v.Type(), "vs", typ)
	}
	if valS, ok := val.(string); ok {
		if v.String() != valS {
			t.Error(key, "value assertion failed", v.String(), "vs", valS)
		}
	}
	if valN, ok := val.(int64); ok {
		if v.Int() != valN {
			t.Error(key, "value assertion failed", v.Int(), "vs", valN)
		}
	}
}

func testComplexFmt(t testing.TB, key string, src, dst []byte) {
	vec.Reset()
	bbuf.Reset()
	err := vec.Parse(src, false)
	if err != nil {
		t.Error(key, err)
	}
	err = vec.Beautify(&bbuf)
	if err != nil {
		t.Error(key, err)
	}
	if !bytes.Equal(bbuf.Bytes(), dst) {
		t.Error(key, "fmt assertion failed")
	}
}

func TestVector_ParseComplex0(t *testing.T) {
	testComplex(t, "complex 0", complex0, complexPath0, TypeStr, "SGML")
}

func TestVector_FmtComplex0(t *testing.T) {
	testComplexFmt(t, "complex 0", complex0, complexFmt0)
}

func TestVector_ParseComplex1(t *testing.T) {
	testComplex(t, "complex 1", complex1, complexPath1, TypeStr, "646 555-4567")
}

func TestVector_FmtComplex1(t *testing.T) {
	testComplexFmt(t, "complex 1", complex1, complexFmt1)
}

func TestVector_ParseComplex2(t *testing.T) {
	testComplex(t, "complex 2", complex2, complexPath2, TypeStr, "number")
}

func TestVector_FmtComplex2(t *testing.T) {
	testComplexFmt(t, "complex 2", complex2, complexFmt2)
}

func TestVector_ParseComplex3(t *testing.T) {
	testComplex(t, "complex 3", complex3, complexPath3, TypeNum, int64(20))
}

func TestVector_FmtComplex3(t *testing.T) {
	testComplexFmt(t, "complex 3", complex3, complexFmt3)
}

func TestVector_ParseComplex4(t *testing.T) {
	testComplex(t, "complex 4", complex4, complexPath4, TypeStr, "10021")
}

func TestVector_FmtComplex4(t *testing.T) {
	testComplexFmt(t, "complex 4", complex4, complexFmt4)
}

func TestVector_ParseComplex5(t *testing.T) {
	testComplex(t, "complex 5", complex5, complexPath5, TypeStr, "3")
}

func TestVector_FmtComplex5(t *testing.T) {
	testComplexFmt(t, "complex 5", complex5, complexFmt5)
}

func BenchmarkVector_ParseComplex0(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 0", complex0, complexPath0, TypeStr, "SGML")
	}
}

func BenchmarkVector_ParseComplex1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 1", complex1, complexPath1, TypeStr, "646 555-4567")
	}
}

func BenchmarkVector_ParseComplex2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 2", complex2, complexPath2, TypeStr, "number")
	}
}

func BenchmarkVector_ParseComplex3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 3", complex3, complexPath3, TypeNum, int64(20))
	}
}

func BenchmarkVector_ParseComplex4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 4", complex4, complexPath4, TypeStr, "10021")
	}
}

func BenchmarkVector_ParseComplex5(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 5", complex5, complexPath5, TypeStr, "3")
	}
}
