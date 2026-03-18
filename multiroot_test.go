package jsonvector

import (
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/koykov/vector"
)

type multirootStage struct {
	key string
	buf []byte
}

var (
	multiroots    [][]byte
	multirootsFmt [][]byte

	multirootStages    []multirootStage
	multirootStagesReg = map[string]int{}
)

func init() {
	_ = filepath.Walk("testdata/multiroot", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			b, _ := os.ReadFile(path)
			multiroots = append(multiroots, b)
			vec := NewVector()
			if err = vec.ParseCopy(b); err != nil {
				os.Exit(1)
			}
			var buf bytes.Buffer
			_ = vec.Beautify(&buf)
			multirootsFmt = append(multirootsFmt, buf.Bytes())
			return nil
		}
		return nil
	})
	for i := 10; i <= 100; i += 10 {
		key := strconv.Itoa(i)
		st := multirootStage{key: key}
		for j := 0; j < i; j++ {
			st.buf = append(st.buf, multiroots[j%len(multiroots)]...)
			if j > 0 {
				switch rand.Intn(5) {
				case 0:
					st.buf = append(st.buf, '\n')
				case 1:
					st.buf = append(st.buf, '\r')
				case 2:
					st.buf = append(st.buf, '\t')
				case 3:
					st.buf = append(st.buf, ' ')
				case 4:
					st.buf = append(st.buf, "\n\r\n\n\r\r\t\t   "...)
				}
			}
		}
		multirootStages = append(multirootStages, st)
		multirootStagesReg[key] = len(multirootStages) - 1
	}
}

func TestMultiroot(t *testing.T) {
	for i := 10; i <= 100; i += 10 {
		key := strconv.Itoa(i)
		t.Run(key, func(t *testing.T) {
			idx := multirootStagesReg[key]
			st := &multirootStages[idx]

			vec := Acquire()
			defer Release(vec)
			if err := vec.ParseCopy(st.buf); err != nil {
				t.Error(err)
			}

			for j := 0; ; j++ {
				root := vec.RootByIndex(j)
				if root.Type() == vector.TypeNull {
					break
				}
				var buf bytes.Buffer
				_ = root.Beautify(&buf)

				origin := multirootsFmt[j%len(multirootsFmt)]
				fmtv := buf.Bytes()
				if !bytes.Equal(origin, fmtv) {
					t.Errorf("key %s root %d", key, j)
				}
			}
		})
	}
}

func BenchmarkMultiroot(b *testing.B) {
	for i := 10; i <= 100; i += 10 {
		key := strconv.Itoa(i)
		b.Run(key, func(b *testing.B) {
			idx := multirootStagesReg[key]
			st := &multirootStages[idx]

			vec := Acquire()
			defer Release(vec)
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				vec.Reset()
				if err := vec.Parse(st.buf); err != nil {
					b.Error(err)
				}
			}
		})
	}
}
