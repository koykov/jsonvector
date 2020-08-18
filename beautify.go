package jsonvector

import "io"

var (
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

func (vec *Vector) beautify(w io.Writer, v *Node, depth int) (err error) {
	switch v.t {
	case TypeNull:
		_, err = w.Write(btNull)
	case TypeNum, TypeBool:
		_, err = w.Write(v.ForceBytes())
	case TypeStr:
		_, err = w.Write(btQuote)
		_, err = w.Write(v.RawBytes())
		_, err = w.Write(btQuote)
	case TypeArr:
		ci := v.ChildIdx()
		if len(ci) == 0 {
			_, err = w.Write(btArrE)
		} else {
			_, err = w.Write(btArrO)
			_, err = w.Write(btNl)
			for cnt, i := range ci {
				c := vec.v[i]
				if cnt > 0 {
					_, err = w.Write(btComma)
					_, err = w.Write(btNl)
				}
				writePad(w, c.d)
				err = vec.beautify(w, &c, depth+1)
			}
			_, err = w.Write(btNl)
			writePad(w, v.d)
			_, err = w.Write(btArrC)
		}
	case TypeObj:
		ci := v.ChildIdx()
		if len(ci) == 0 {
			_, err = w.Write(btObjE)
		} else {
			_, err = w.Write(btObjO)
			_, err = w.Write(btNl)
			for cnt, i := range ci {
				c := vec.v[i]
				if cnt > 0 {
					_, err = w.Write(btComma)
					_, err = w.Write(btNl)
				}
				writePad(w, c.d)
				_, err = w.Write(btQuote)
				_, err = w.Write(c.k.rawBytes())
				_, err = w.Write(btQuote)
				_, err = w.Write(btDotDot)
				_, err = w.Write(btSpace)
				err = vec.beautify(w, &c, depth+1)
			}
			_, err = w.Write(btNl)
			writePad(w, v.d)
			_, err = w.Write(btObjC)
		}
	}
	return
}

func writePad(w io.Writer, cnt int) {
	for i := 0; i < cnt; i++ {
		_, _ = w.Write(btTab)
	}
}