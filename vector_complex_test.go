package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/vector"
)

var (
	complex0     = []byte(`{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}`)
	complexPath0 = []string{"glossary", "GlossDiv", "GlossList", "GlossEntry", "ID"}
	complex1     = []byte(`{"firstName":"John","lastName":"Smith","isAlive":true,"age":27,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`)
	complexPath1 = []string{"phoneNumbers", "1", "number"}
	complex2     = []byte(`{"$schema":"http://json-schema.org/schema#","title":"Product","type":"object","required":["id","name","price"],"properties":{"id":{"type":"number","description":"Product identifier"},"name":{"type":"string","description":"Name of the product"},"price":{"type":"number","minimum":0},"tags":{"type":"array","items":{"type":"string"}},"stock":{"type":"object","properties":{"warehouse":{"type":"number"},"retail":{"type":"number"}}}}}`)
	complexPath2 = []string{"properties", "stock", "properties", "retail", "type"}
	complex3     = []byte(`{"id":1,"name":"Foo","price":123,"tags":["Bar","Eek"],"stock":{"warehouse":300,"retail":20}}`)
	complexPath3 = []string{"stock", "retail"}
	complex4     = []byte(`{"first name":"John","last name":"Smith","age":25,"address":{"street address":"21 2nd Street","city":"New York","state":"NY","postal code":"10021"},"phone numbers":[{"type":"home","number":"212 555-1234"},{"type":"fax","number":"646 555-4567"}],"sex":{"type":"male"}}`)
	complexPath4 = []string{"address", "postal code"}
	complex5     = []byte(`{"quiz":{"sport":{"q1":{"question":"Which one is correct team name in NBA?","options":["New York Bulls","Los Angeles Kings","Golden State Warriros","Huston Rocket"],"answer":"Huston Rocket"}},"maths":{"q1":{"question":"5 + 7 = ?","options":["10","11","12","13"],"answer":"12"},"q2":{"question":"12 - 8 = ?","options":["1","2","3","4"],"answer":"4"}}}}`)
	complexPath5 = []string{"quiz", "maths", "q2", "options", "2"}
)

func testComplex(t testing.TB, key string, src []byte, path []string, typ vector.Type, val interface{}) {
	vec.Reset()
	err := vec.Parse(src)
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
		i, _ := v.Int()
		if i != valN {
			t.Error(key, "value assertion failed", i, "vs", valN)
		}
	}
}

func TestComplex(t *testing.T) {
	vec := NewVector()
	t.Run("complex0", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "glossary.GlossDiv.GlossList.GlossEntry.ID", "SGML")
		assertNode(t, vec, "glossary.GlossDiv.GlossList.GlossEntry.Abbrev", "ISO 8879:1986")
		assertNode(t, vec, "glossary.GlossDiv.GlossList.GlossEntry.GlossDef.para", "A meta-markup language, used to create markup languages such as DocBook.")
		assertNode(t, vec, "glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.1", "XML")
	})
	t.Run("complex1", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "lastName", "Smith")
		assertNode(t, vec, "address.state", "NY")
		assertNode(t, vec, "phoneNumbers.1.number", "646 555-4567")
	})
	t.Run("complex2", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "$schema", "http://json-schema.org/schema#")
		assertNode(t, vec, "properties.name.type", "string")
		assertNode(t, vec, "properties.price.minimum", 0)
		assertNode(t, vec, "properties.stock.properties.warehouse.type", "number")
	})
	t.Run("complex3", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "id", 1)
		assertNode(t, vec, "stock.warehouse", 300)
		assertNode(t, vec, "price", 123)
	})
	t.Run("complex4", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "last name", "Smith")
		assertNode(t, vec, "address.street address", "21 2nd Street")
		assertNode(t, vec, "phone numbers.1.number", "646 555-4567")
	})
	t.Run("complex5", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertNode(t, vec, "quiz.sport.q1.question", "Which one is correct team name in NBA?")
		assertNode(t, vec, "quiz.sport.q1.options.2", "Golden State Warriros")
		assertNode(t, vec, "quiz.sport.q1.answer", "Huston Rocket")
		assertNode(t, vec, "quiz.maths.q2.question", "12 - 8 = ?")
		assertNode(t, vec, "quiz.maths.q2.answer", "4")
	})
}

func TestComplexFmt(t *testing.T) {
	vec := NewVector()
	buf := &bytes.Buffer{}
	t.Run("complex0", func(t *testing.T) { assertFmt(t, vec, buf) })
	t.Run("complex1", func(t *testing.T) { assertFmt(t, vec, buf) })
	t.Run("complex2", func(t *testing.T) { assertFmt(t, vec, buf) })
	t.Run("complex3", func(t *testing.T) { assertFmt(t, vec, buf) })
	t.Run("complex4", func(t *testing.T) { assertFmt(t, vec, buf) })
	t.Run("complex5", func(t *testing.T) { assertFmt(t, vec, buf) })
}

func BenchmarkVector_ParseComplex0(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 0", complex0, complexPath0, vector.TypeStr, "SGML")
	}
}

func BenchmarkVector_ParseComplex1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 1", complex1, complexPath1, vector.TypeStr, "646 555-4567")
	}
}

func BenchmarkVector_ParseComplex2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 2", complex2, complexPath2, vector.TypeStr, "number")
	}
}

func BenchmarkVector_ParseComplex3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 3", complex3, complexPath3, vector.TypeNum, int64(20))
	}
}

func BenchmarkVector_ParseComplex4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 4", complex4, complexPath4, vector.TypeStr, "10021")
	}
}

func BenchmarkVector_ParseComplex5(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testComplex(b, "complex 5", complex5, complexPath5, vector.TypeStr, "3")
	}
}
