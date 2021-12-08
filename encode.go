package fxml

import (
	"encoding/xml"
	"io"
)

func assert(e interface{}) {
	if e != nil {
		panic(e)
	}
}

func (xt XMLTree) encode(e *xml.Encoder) {
	assert(e.EncodeToken(xml.StartElement{Name: xt.Name, Attr: xt.Attr}))
	if xt.Comment != "" {
		assert(e.EncodeToken(xml.Comment(xt.Comment)))
	}
	if xt.Directive != "" {
		assert(e.EncodeToken(xml.Directive(xt.Directive)))
	}
	if xt.Text != "" {
		assert(e.EncodeToken(xml.CharData(xt.Text)))
	}
	for _, c := range xt.Children {
		c.encode(e)
	}
	assert(e.EncodeToken(xml.EndElement{xt.Name}))
}

func (xt XMLTree) Encode(w io.Writer, full bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	e := xml.NewEncoder(w)
	e.Indent("", "    ")
	if full {
		assert(e.EncodeToken(xml.ProcInst{
			Target: "xml",
			Inst:   []byte(`version="1.0" encoding="UTF-8"`),
		}))
	}
	xt.encode(e)
	return e.Flush()
}
