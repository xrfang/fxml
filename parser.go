package fxml

import (
	"encoding/xml"
	"io"
	"regexp"
)

type (
	XMLTree struct {
		Name      xml.Name
		Attr      []xml.Attr
		Comment   string
		Directive string
		Text      string
		Children  []XMLTree
	}
)

func (xt *XMLTree) Parse(xd *xml.Decoder) error {
	for {
		t, err := xd.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch t := t.(type) {
		case xml.StartElement:
			child := XMLTree{
				Name: t.Name,
				Attr: t.Attr,
			}
			err := child.Parse(xd)
			if err != nil {
				return err
			}
			xt.Children = append(xt.Children, child)
		case xml.EndElement:
			return nil
		case xml.CharData:
			xt.Text = string(t)
		case xml.Comment:
			xt.Comment = string(t)
		case xml.ProcInst:
			if t.Target == "xml" {
				rx := regexp.MustCompile(`encoding="(.+?)"`)
				enc := "UTF-8"
				encs := rx.FindSubmatch(t.Inst)
				if len(encs) == 2 {
					enc = string(encs[1])
				}
				err := conv.Init(enc)
				if err != nil {
					return err
				}
			}
		case xml.Directive:
			xt.Directive = string(t)
		}
	}
}
