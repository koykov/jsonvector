package jsonvector

import (
	"io"

	"github.com/koykov/vector"
)

var (
	// Byte constants.
	btSpace  = []byte(` `)
	btQuote  = []byte(`"`)
	btComma  = []byte(`,`)
	btDotDot = []byte(`:`)
	btNl     = []byte("\n")
	btTab    = []byte("\t")
	btArrO   = []byte(`[`)
	btArrC   = []byte(`]`)
	btArrE   = []byte(`[]`)
	btObjO   = []byte(`{`)
	btObjC   = []byte(`}`)
	btObjE   = []byte(`{}`)
	btNull   = []byte(`null`)
)

func serialize(w io.Writer, node *vector.Node, depth int, indent bool) (err error) {
	switch node.Type() {
	case vector.TypeNull:
		_, err = w.Write(btNull)
	case vector.TypeNum, vector.TypeBool:
		_, err = w.Write(node.ForceBytes())
	case vector.TypeStr:
		_, err = w.Write(btQuote)
		_, err = w.Write(node.RawBytes())
		_, err = w.Write(btQuote)
	case vector.TypeArr:
		if node.Limit() == 0 {
			_, err = w.Write(btArrE)
		} else {
			_, err = w.Write(btArrO)
			if indent {
				_, err = w.Write(btNl)
			}
			node.Each(func(idx int, node *vector.Node) {
				if idx > 0 {
					_, err = w.Write(btComma)
					if indent {
						_, err = w.Write(btNl)
					}
				}
				if indent {
					writePad(w, node.Depth())
				}
				err = serialize(w, node, depth+1, indent)
			})
			if indent {
				_, err = w.Write(btNl)
			}
			if indent {
				writePad(w, node.Depth())
			}
			_, err = w.Write(btArrC)
		}
	case vector.TypeObj:
		if node.Limit() == 0 {
			_, err = w.Write(btObjE)
		} else {
			_, err = w.Write(btObjO)
			if indent {
				_, err = w.Write(btNl)
			}
			node.Each(func(idx int, node *vector.Node) {
				if idx > 0 {
					_, err = w.Write(btComma)
					if indent {
						_, err = w.Write(btNl)
					}
				}
				if indent {
					writePad(w, node.Depth())
				}
				_, err = w.Write(btQuote)
				_, err = w.Write(node.KeyBytes())
				_, err = w.Write(btQuote)
				_, err = w.Write(btDotDot)
				if indent {
					_, err = w.Write(btSpace)
				}
				err = serialize(w, node, depth+1, indent)
			})
			if indent {
				_, err = w.Write(btNl)
				writePad(w, node.Depth())
			}
			_, err = w.Write(btObjC)
		}
	}
	return
}

// Write number of tabs to w.
func writePad(w io.Writer, cnt int) {
	for i := 0; i < cnt; i++ {
		_, _ = w.Write(btTab)
	}
}
