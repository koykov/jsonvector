package jsonvector

import (
	"bytes"
	"github.com/koykov/bytealg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/vector"
)

type stage struct {
	key string

	origin, fmt []byte
}

var (
	stages []stage
)

func init() {
	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" && !strings.Contains(filepath.Base(path), ".fmt.json") {
			st := stage{}
			st.key = strings.Replace(filepath.Base(path), ".json", "", 1)
			st.origin, _ = ioutil.ReadFile(path)
			if st.fmt, _ = ioutil.ReadFile(strings.Replace(path, ".json", ".fmt.json", 1)); len(st.fmt) > 0 {
				st.fmt = bytealg.Trim(st.fmt, btNl)
			}
			stages = append(stages, st)
		}
		return nil
	})
}

func getStage(key string) (st *stage) {
	for i := 0; i < len(stages); i++ {
		st1 := &stages[i]
		if st1.key == key {
			st = st1
		}
	}
	return st
}

func getTBName(tb testing.TB) string {
	key := tb.Name()
	return key[strings.Index(key, "/")+1:]
}

func assertParse(tb testing.TB, dst *Vector, err error, errOffset int) *Vector {
	key := getTBName(tb)
	st := getStage(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	dst.Reset()
	err1 := dst.ParseCopy(st.origin)
	if err1 != nil {
		if err != nil {
			if err != err1 || dst.ErrorOffset() != errOffset {
				tb.Fatalf(`error mismatch, need "%s" at %d, got "%s" at %d`, err.Error(), errOffset, err1.Error(), dst.ErrorOffset())
			}
		} else {
			tb.Fatalf(`err "%s" caught by offset %d`, err1.Error(), dst.ErrorOffset())
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

func assertNode(tb testing.TB, vec *Vector, path string, val interface{}) {
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
