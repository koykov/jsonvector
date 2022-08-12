package jsonvector

import (
	"bytes"
	"testing"

	"github.com/koykov/vector"
)

func TestScalar(t *testing.T) {
	vec := NewVector()
	t.Run("scalarNull", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNull)
	})
	t.Run("scalarString", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeStr)
		assertNode(t, vec, "", "foo bar string")
	})
	t.Run("scalarStringQuoted", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeStr)
		assertNode(t, vec, "", `foo "bar" string`)
	})
	t.Run("scalarNumber", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNum)
		assertNode(t, vec, "", 123456)
	})
	t.Run("scalarFloat", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNum)
		assertNode(t, vec, "", 123.456)
	})
	t.Run("scalarFloatScientific", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeNum)
		assertNode(t, vec, "", 3.7e-5)
	})
	t.Run("scalarTrue", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertNode(t, vec, "", true)
	})
	t.Run("scalarFalse", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeBool)
		assertNode(t, vec, "", false)
	})
}

func TestArray(t *testing.T) {
	vec := NewVector()
	t.Run("arrayNumber", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeArr)
		assertLen(t, vec, "", 5)
		assertNode(t, vec, "0", 1)
		assertNode(t, vec, "2", 3)
		assertNode(t, vec, "4", 5)
	})
	t.Run("arrayString", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeArr)
		assertLen(t, vec, "", 3)
		assertNode(t, vec, "0", "foo")
		assertNode(t, vec, "1", "bar")
		assertNode(t, vec, "2", "string")
	})
	t.Run("arrayStringQuoted", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeArr)
		assertLen(t, vec, "", 2)
		assertNode(t, vec, "0", `quoted "str" value`)
		assertNode(t, vec, "1", "foo")
	})
	t.Run("arrayFloat", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeArr)
		assertLen(t, vec, "", 2)
		assertNode(t, vec, "0", 3.14156)
		assertNode(t, vec, "1", 6.23e-4)
	})
}

func TestObject(t *testing.T) {
	vec := NewVector()
	t.Run("objectNumber", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 3)
		assertNode(t, vec, "a", 1)
		assertNode(t, vec, "b", 2)
		assertNode(t, vec, "c", 3)
	})
	t.Run("objectString", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 3)
		assertNode(t, vec, "a", "foo")
		assertNode(t, vec, "b", "bar")
		assertNode(t, vec, "c", "string")
	})
	t.Run("objectStringQuoted", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 2)
		assertNode(t, vec, "key0", `"quoted"`)
		assertNode(t, vec, `key"1"`, "str")
	})
	t.Run("objectFloat", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 2)
		assertNode(t, vec, "pi", 3.1415)
		assertNode(t, vec, "e", 2.718281828459045)
	})
	t.Run("objectFmt", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 3)
		assertNode(t, vec, "c", 15)
		assertType(t, vec, "foo", vector.TypeNull)
		assertNode(t, vec, "bar", `qwerty "encoded"`)
	})
	t.Run("objectFmt1", func(t *testing.T) {
		vec = assertParse(t, vec, nil, 0)
		assertType(t, vec, "", vector.TypeObj)
		assertLen(t, vec, "", 2)
		assertNode(t, vec, "a", true)
		assertNode(t, vec, "b.c", "foo")
		assertLen(t, vec, "b.d", 3)
		assertNode(t, vec, "b.d.0", 5)
		assertNode(t, vec, "b.d.1", 3.1415)
		assertNode(t, vec, "b.d.2", 812.48927)
	})
}

func TestError(t *testing.T) {
	vec := NewVector()
	t.Run("badToken", func(t *testing.T) { assertParse(t, vec, vector.ErrUnexpId, 0) })
	t.Run("badUnclosedString", func(t *testing.T) { assertParse(t, vec, vector.ErrUnexpEOS, 24) })
	t.Run("badFloatSeparator", func(t *testing.T) { assertParse(t, vec, vector.ErrUnparsedTail, 1) })
	t.Run("badUnparsedTail", func(t *testing.T) { assertParse(t, vec, vector.ErrUnparsedTail, 16) })
}

func TestMulti(t *testing.T) {
	vec := NewVector()
	var buf bytes.Buffer
	t.Run("multi0", func(t *testing.T) {
		vec = assertParseMulti(t, vec, &buf, nil, 0)
	})
}

func TestSort(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		vec := NewVector()
		var (
			buf bytes.Buffer
			st  *stage
		)
		vec, st = assertParseStage(t, vec, nil, 0)
		vec.Root().SortKeys()
		_ = vec.Root().Beautify(&buf)
		if !bytes.Equal(buf.Bytes(), st.fmt) {
			t.Error("sort failed")
		}
	})
	t.Run("array", func(t *testing.T) {
		vec := NewVector()
		var (
			buf bytes.Buffer
			st  *stage
		)
		vec, st = assertParseStage(t, vec, nil, 0)
		vec.Root().Sort()
		_ = vec.Root().Beautify(&buf)
		if !bytes.Equal(buf.Bytes(), st.fmt) {
			t.Error("sort failed")
		}
	})
}

func BenchmarkScalar(b *testing.B) {
	b.Run("scalarNull", func(b *testing.B) { bench(b, func(vec *Vector) { assertType(b, vec, "", vector.TypeNull) }) })
	b.Run("scalarString", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeStr)
			assertNode(b, vec, "", "foo bar string")
		})
	})
	b.Run("scalarStringQuoted", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeStr)
			assertNode(b, vec, "", `foo "bar" string`)
		})
	})
	b.Run("scalarNumber", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNum)
			assertNode(b, vec, "", 123456)
		})
	})
	b.Run("scalarFloat", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNum)
			assertNode(b, vec, "", 123.456)
		})
	})
	b.Run("scalarFloatScientific", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeNum)
			assertNode(b, vec, "", 3.7e-5)
		})
	})
	b.Run("scalarTrue", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeBool)
			assertNode(b, vec, "", true)
		})
	})
	b.Run("scalarFalse", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeBool)
			assertNode(b, vec, "", false)
		})
	})
}

func BenchmarkArray(b *testing.B) {
	b.Run("arrayNumber", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeArr)
			assertLen(b, vec, "", 5)
			assertNode(b, vec, "0", 1)
			assertNode(b, vec, "2", 3)
			assertNode(b, vec, "4", 5)
		})
	})
	b.Run("arrayString", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeArr)
			assertLen(b, vec, "", 3)
			assertNode(b, vec, "0", "foo")
			assertNode(b, vec, "1", "bar")
			assertNode(b, vec, "2", "string")
		})
	})
	b.Run("arrayStringQuoted", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeArr)
			assertLen(b, vec, "", 2)
			assertNode(b, vec, "0", `quoted "str" value`)
			assertNode(b, vec, "1", "foo")
		})
	})
	b.Run("arrayFloat", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeArr)
			assertLen(b, vec, "", 2)
			assertNode(b, vec, "0", 3.14156)
			assertNode(b, vec, "1", 6.23e-4)
		})
	})
}

func BenchmarkObject(b *testing.B) {
	b.Run("objectNumber", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 3)
			assertNode(b, vec, "a", 1)
			assertNode(b, vec, "b", 2)
			assertNode(b, vec, "c", 3)
		})
	})
	b.Run("objectString", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 3)
			assertNode(b, vec, "a", "foo")
			assertNode(b, vec, "b", "bar")
			assertNode(b, vec, "c", "string")
		})
	})
	b.Run("objectStringQuoted", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 2)
			assertNode(b, vec, "key0", `"quoted"`)
			assertNode(b, vec, `key"1"`, "str")
		})
	})
	b.Run("objectFloat", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 2)
			assertNode(b, vec, "pi", 3.1415)
			assertNode(b, vec, "e", 2.718281828459045)
		})
	})
	b.Run("objectFmt", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 3)
			assertNode(b, vec, "c", 15)
			assertType(b, vec, "foo", vector.TypeNull)
			assertNode(b, vec, "bar", `qwerty "encoded"`)
		})
	})
	b.Run("objectFmt1", func(b *testing.B) {
		bench(b, func(vec *Vector) {
			assertType(b, vec, "", vector.TypeObj)
			assertLen(b, vec, "", 2)
			assertNode(b, vec, "a", true)
			assertNode(b, vec, "b.c", "foo")
			assertLen(b, vec, "b.d", 3)
			assertNode(b, vec, "b.d.0", 5)
			assertNode(b, vec, "b.d.1", 3.1415)
			assertNode(b, vec, "b.d.2", 812.48927)
		})
	})
}

func BenchmarkMulti(b *testing.B) {
	var buf bytes.Buffer
	b.Run("multi0", func(b *testing.B) {
		benchMulti(b, &buf, func(vec *Vector) {})
	})
}
