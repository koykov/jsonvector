package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/vector"
)

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

func BenchmarkComplex(b *testing.B) {
	b.Run("complex0", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "glossary.GlossDiv.GlossList.GlossEntry.ID", "SGML")
			assertNode(b, vec, "glossary.GlossDiv.GlossList.GlossEntry.Abbrev", "ISO 8879:1986")
			assertNode(b, vec, "glossary.GlossDiv.GlossList.GlossEntry.GlossDef.para", "A meta-markup language, used to create markup languages such as DocBook.")
			assertNode(b, vec, "glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.1", "XML")
		})
	})
	b.Run("complex1", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "lastName", "Smith")
			assertNode(b, vec, "address.state", "NY")
			assertNode(b, vec, "phoneNumbers.1.number", "646 555-4567")
		})
	})
	b.Run("complex2", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "$schema", "http://json-schema.org/schema#")
			assertNode(b, vec, "properties.name.type", "string")
			assertNode(b, vec, "properties.price.minimum", 0)
			assertNode(b, vec, "properties.stock.properties.warehouse.type", "number")
		})
	})
	b.Run("complex3", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "id", 1)
			assertNode(b, vec, "stock.warehouse", 300)
			assertNode(b, vec, "price", 123)
		})
	})
	b.Run("complex4", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "last name", "Smith")
			assertNode(b, vec, "address.street address", "21 2nd Street")
			assertNode(b, vec, "phone numbers.1.number", "646 555-4567")
		})
	})
	b.Run("complex5", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			vec = assertParse(b, vec, nil, 0)
			assertType(b, vec, "", vector.TypeObj)
			assertNode(b, vec, "quiz.sport.q1.question", "Which one is correct team name in NBA?")
			assertNode(b, vec, "quiz.sport.q1.options.2", "Golden State Warriros")
			assertNode(b, vec, "quiz.sport.q1.answer", "Huston Rocket")
			assertNode(b, vec, "quiz.maths.q2.question", "12 - 8 = ?")
			assertNode(b, vec, "quiz.maths.q2.answer", "4")
		})
	})
}

func BenchmarkComplexFmt(b *testing.B) {
	b.Run("complex0", func(b *testing.B) { benchFmt(b) })
	b.Run("complex1", func(b *testing.B) { benchFmt(b) })
	b.Run("complex2", func(b *testing.B) { benchFmt(b) })
	b.Run("complex3", func(b *testing.B) { benchFmt(b) })
	b.Run("complex4", func(b *testing.B) { benchFmt(b) })
	b.Run("complex5", func(b *testing.B) { benchFmt(b) })
}
