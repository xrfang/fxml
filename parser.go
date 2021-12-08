package fxml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"path"
	"regexp"
)

type (
	XTraverser func(string, XMLTree) bool
	XMLTree    struct {
		Name      xml.Name
		Attr      []xml.Attr `json:",omitempty"`
		Comment   string     `json:",omitempty"`
		Directive string     `json:",omitempty"`
		Text      string     `json:",omitempty"`
		Children  []XMLTree  `json:",omitempty"`
	}
)

func (xt *XMLTree) parse(xd *xml.Decoder) error {
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
			err := child.parse(xd)
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

func (xt XMLTree) traverse(pfx string, v XTraverser) bool {
	p := path.Join(pfx, xt.Name.Local)
	if !v(p, xt) {
		return false
	}
	for _, c := range xt.Children {
		if !c.traverse(p, v) {
			return false
		}
	}
	return true
}

func (xt XMLTree) Traverse(v XTraverser) bool {
	return xt.traverse("", v)
}

func (xt XMLTree) ToJSON() string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(xt); err != nil {
		panic(err)
	}
	return buf.String()
}

func FromJSON(js string) (*XMLTree, error) {
	var xt XMLTree
	err := json.Unmarshal([]byte(js), &xt)
	if err != nil {
		return nil, err
	}
	return &xt, nil
}

func Parse(r io.Reader) (*XMLTree, error) {
	d := xml.NewDecoder(r)
	var xt XMLTree
	err := xt.parse(d)
	if err != nil {
		return nil, err
	}
	if xt.Name.Space == "" && xt.Name.Local == "" && len(xt.Children) == 1 {
		xt = xt.Children[0]
	}
	return &xt, nil
}
