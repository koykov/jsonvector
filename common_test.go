package jsonvector

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stage struct {
	key string

	origin, fmt, flat []byte
}

type multiStage struct {
	key string
	buf []stage
}

var (
	stages         []stage
	stagesReg      = map[string]int{}
	multiStages    []multiStage
	multiStagesReg = map[string]int{}
)

func init() {
	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" && !strings.Contains(filepath.Base(path), ".fmt.json") {
			st := stage{}
			st.key = strings.Replace(filepath.Base(path), ".json", "", 1)
			st.origin, _ = os.ReadFile(path)
			if st.fmt, _ = os.ReadFile(strings.Replace(path, ".json", ".fmt.json", 1)); len(st.fmt) > 0 {
				st.fmt = bytealg.Trim(st.fmt, btNl)
			}
			if st.flat, _ = os.ReadFile(strings.Replace(path, ".json", ".flat.json", 1)); len(st.flat) > 0 {
				st.flat = bytealg.Trim(st.flat, btNl)
			}
			stages = append(stages, st)
			stagesReg[st.key] = len(stages) - 1
			return nil
		}

		if info.IsDir() && path != "testdata" {
			mstg := multiStage{key: filepath.Base(path)}
			_ = filepath.Walk(path, func(path1 string, info1 os.FileInfo, err1 error) error {
				if filepath.Ext(path1) == ".json" && !strings.Contains(filepath.Base(path1), ".fmt.json") {
					st := stage{}
					st.key = strings.Replace(filepath.Base(path1), ".json", "", 1)
					st.origin, _ = os.ReadFile(path1)
					if st.fmt, _ = os.ReadFile(strings.Replace(path1, ".json", ".fmt.json", 1)); len(st.fmt) > 0 {
						st.fmt = bytealg.Trim(st.fmt, btNl)
					}
					if st.flat, _ = os.ReadFile(strings.Replace(path1, ".json", ".flat.json", 1)); len(st.flat) > 0 {
						st.flat = bytealg.Trim(st.flat, btNl)
					}
					mstg.buf = append(mstg.buf, st)
					return nil
				}
				return nil
			})
			multiStages = append(multiStages, mstg)
			multiStagesReg[mstg.key] = len(multiStages) - 1
		}
		return nil
	})
}

func getStage(key string) *stage {
	i, ok := stagesReg[key]
	if !ok {
		return nil
	}
	return &stages[i]
}

func getStageMulti(key string) *multiStage {
	i, ok := multiStagesReg[key]
	if !ok {
		return nil
	}
	return &multiStages[i]
}

func getTBName(tb testing.TB) string {
	key := tb.Name()
	return key[strings.Index(key, "/")+1:]
}

func assertParse(tb testing.TB, dst *Vector, err error, errOffset int) *Vector {
	key := getTBName(tb)
	st := getStage(key)
	require.NotNil(tb, st, "stage not found")
	dst = assertParseStage(tb, st, dst, err, errOffset)
	return dst
}

func assertParseStage(tb testing.TB, st *stage, dst *Vector, err error, errOffset int) *Vector {
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			assert.True(tb, errors.Is(err1, err), `error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			assert.Equal(tb, errOffset, dst.ErrorOffset(), "error offset mismatch")
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
		}
	}
	return dst
}

func assertParseMulti(tb testing.TB, dst *Vector, buf *bytes.Buffer, err error, errOffset int) *Vector {
	key := getTBName(tb)
	mst := getStageMulti(key)
	require.NotNil(tb, mst, "stage not found")
	return assertParseStageMulti(tb, mst, dst, buf, err, errOffset)
}

func assertParseStageMulti(tb testing.TB, mst *multiStage, dst *Vector, buf *bytes.Buffer, err error, errOffset int) *Vector {
	dst.Reset()
	for i := 0; i < len(mst.buf); i++ {
		st := &mst.buf[i]
		err1 := dst.ParseCopy(st.origin)
		if err1 != nil {
			if err != nil {
				assert.True(tb, errors.Is(err1, err), `error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
				assert.Equal(tb, errOffset, dst.ErrorOffset(), "error offset mismatch")
			} else {
				tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
			}
		}
		root := dst.RootTop()
		buf.Reset()
		_ = root.Beautify(buf)
		assert.True(tb, bytes.Equal(st.fmt, buf.Bytes()), "node mismatch")
	}
	return dst
}

func assertType(tb testing.TB, vec *Vector, path string, typ vector.Type) {
	assert.Equal(tb, typ, vec.Dot(path).Type(), "type mismatch")
}

func assertLen(tb testing.TB, vec *Vector, path string, len int) {
	assert.Equal(tb, len, vec.Dot(path).Limit(), "length mismatch")
}

func assertNode(tb testing.TB, vec *Vector, path string, val any) {
	node := vec.Dot(path)
	switch v := val.(type) {
	case string:
		assert.True(tb, v == node.String(), "value mismatch")
	case int:
		i, _ := node.Int()
		assert.True(tb, v == int(i), "value mismatch")
	case float64:
		f, _ := node.Float()
		assert.True(tb, v == f, "value mismatch")
	case bool:
		assert.True(tb, v == node.Bool(), "value mismatch")
	default:
		tb.Errorf("unsupported type %T", val)
	}
}

func assertFmt(tb testing.TB, vec *Vector, buf *bytes.Buffer) {
	key := getTBName(tb)
	st := getStage(key)
	require.NotNil(tb, st, "stage not found")
	vec.Reset()
	buf.Reset()
	_ = vec.ParseCopy(st.origin)
	err := vec.Beautify(buf)
	assert.NoError(tb, err)
	assert.True(tb, bytes.Equal(st.fmt, buf.Bytes()))
}

func bench(b *testing.B, fn func(vec *Vector)) {
	key := getTBName(b)
	st := getStage(key)
	if st == nil {
		b.Fatal("stage not found")
	}

	vec := NewVector()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec = assertParseStage(b, st, vec, nil, 0)
		fn(vec)
	}
}

func benchFmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec := Acquire()
		assertFmt(b, vec, buf)
		Release(vec)
	}
}

func benchMulti(b *testing.B, buf *bytes.Buffer, fn func(vec *Vector)) {
	key := getTBName(b)
	mst := getStageMulti(key)
	if mst == nil {
		b.Fatal("stage not found")
	}

	vec := NewVector()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec = assertParseStageMulti(b, mst, vec, buf, nil, 0)
		fn(vec)
	}
}
