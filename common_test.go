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
	dst, _ = assertParseStage(tb, dst, err, errOffset)
	return dst
}

func assertParseStage(tb testing.TB, dst *Vector, err error, errOffset int) (*Vector, *stage) {
	key := getTBName(tb)
	st := getStage(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			if !errors.Is(err1, err) || dst.ErrorOffset() != errOffset {
				tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			}
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
		}
	}
	return dst, st
}

func assertParseMulti(tb testing.TB, dst *Vector, buf *bytes.Buffer, err error, errOffset int) *Vector {
	key := getTBName(tb)
	mst := getStageMulti(key)
	if mst == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	for i := 0; i < len(mst.buf); i++ {
		st := &mst.buf[i]
		err1 := dst.ParseCopy(st.origin)
		if err1 != nil {
			if err != nil {
				if !errors.Is(err1, err) || dst.ErrorOffset() != errOffset {
					tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
				}
			} else {
				tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
			}
		}
		root := dst.RootTop()
		buf.Reset()
		_ = root.Beautify(buf)
		if fmt1 := buf.Bytes(); !bytes.Equal(fmt1, st.fmt) {
			tb.Fatalf("node mismatch, need '%s'\ngot '%s'", string(st.fmt), string(fmt1))
		}
	}
	return dst
}

func assertType(tb testing.TB, vec *Vector, path string, typ vector.Type) {
	if typ1 := vec.Dot(path).Type(); typ1 != typ {
		tb.Error("type mismatch, need", typ, "got", typ1)
	}
}

func assertLen(tb testing.TB, vec *Vector, path string, len int) {
	if node := vec.Dot(path); node.Limit() != len {
		tb.Error("length mismatch, need", len, "got", node.Limit())
	}
}

func assertNode(tb testing.TB, vec *Vector, path string, val any) {
	node := vec.Dot(path)
	var eq bool
	switch val.(type) {
	case string:
		eq = node.String() == val.(string)
	case int:
		i, _ := node.Int()
		eq = int(i) == val.(int)
	case float64:
		f, _ := node.Float()
		eq = f == val.(float64)
	case bool:
		eq = node.Bool() == val.(bool)
	}
	if !eq {
		tb.Error("value mismatch, need", val, "got", node)
	}
}

func assertFmt(tb testing.TB, vec *Vector, buf *bytes.Buffer) {
	key := getTBName(tb)
	st := getStage(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	vec.Reset()
	buf.Reset()
	_ = vec.ParseCopy(st.origin)
	err := vec.Beautify(buf)
	if err != nil {
		tb.Error(key, err)
	}
	if !bytes.Equal(buf.Bytes(), st.fmt) {
		tb.Error(key, "fmt mismatch")
	}
}

func bench(b *testing.B, fn func(vec *Vector)) {
	vec := NewVector()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec = assertParse(b, vec, nil, 0)
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
	vec := NewVector()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vec = assertParseMulti(b, vec, buf, nil, 0)
		fn(vec)
	}
}
