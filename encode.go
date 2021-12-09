package fxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"unicode/utf8"
)

var (
	//copied from standard library
	escQuot = []byte("&#34;")
	escApos = []byte("&#39;")
	escAmp  = []byte("&amp;")
	escLT   = []byte("&lt;")
	escGT   = []byte("&gt;")
	escTab  = []byte("&#x9;")
	escNL   = []byte("&#xA;")
	escCR   = []byte("&#xD;")
	escFFFD = []byte("\uFFFD")

	begComment = []byte("<!--")
	endComment = []byte("-->")
)

func assert(e interface{}) {
	if e != nil {
		panic(e)
	}
}

// adapted from standard library
func escapeText(s string, escapeNewline bool) string {
	isInCharacterRange := func(r rune) (inrange bool) {
		return r == 0x09 ||
			r == 0x0A ||
			r == 0x0D ||
			r >= 0x20 && r <= 0xD7FF ||
			r >= 0xE000 && r <= 0xFFFD ||
			r >= 0x10000 && r <= 0x10FFFF
	}
	var b bytes.Buffer
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
		case '"':
			esc = escQuot
		case '\'':
			esc = escApos
		case '&':
			esc = escAmp
		case '<':
			esc = escLT
		case '>':
			esc = escGT
		case '\t':
			esc = escTab
		case '\n':
			if !escapeNewline {
				continue
			}
			esc = escNL
		case '\r':
			esc = escCR
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = escFFFD
				break
			}
			continue
		}
		b.WriteString(s[last : i-width])
		b.Write(esc)
		last = i
	}
	b.WriteString(s[last:])
	return b.String()
}

// copied from standard library
func isValidDirective(dir xml.Directive) bool {
	var (
		depth     int
		inquote   uint8
		incomment bool
	)
	for i, c := range dir {
		switch {
		case incomment:
			if c == '>' {
				if n := 1 + i - len(endComment); n >= 0 && bytes.Equal(dir[n:i+1], endComment) {
					incomment = false
				}
			}
			// Just ignore anything in comment
		case inquote != 0:
			if c == inquote {
				inquote = 0
			}
			// Just ignore anything within quotes
		case c == '\'' || c == '"':
			inquote = c
		case c == '<':
			if i+len(begComment) < len(dir) && bytes.Equal(dir[i:i+len(begComment)], begComment) {
				incomment = true
			} else {
				depth++
			}
		case c == '>':
			if depth == 0 {
				return false
			}
			depth--
		}
	}
	return depth == 0 && inquote == 0 && !incomment
}

func nstr(n xml.Name) string {
	if n.Space != "" {
		return n.Space + ":" + n.Local
	}
	return n.Local
}

func encodeToken(w io.Writer, t xml.Token) {
	write := func(ss ...string) {
		for _, s := range ss {
			_, err := w.Write([]byte(s))
			assert(err)
		}
	}
	switch t := t.(type) {
	case xml.StartElement:
		write("<", nstr(t.Name))
		if len(t.Attr) > 0 {
			for _, a := range t.Attr {
				write(" ", nstr(a.Name), `="`, escapeText(a.Value, true), `"`)
			}
		}
		write(">")
	case xml.EndElement:
		write("</", nstr(t.Name), ">")
	case xml.CharData:
		write(escapeText(string(t), false))
	case xml.Comment:
		if bytes.Contains(t, endComment) {
			panic(errors.New("encoding comment containing --> marker"))
		}
		write("<!--", string(t), "-->")
	case xml.Directive:
		if !isValidDirective(t) {
			panic(errors.New("encoding of directive containing wrong < or > markers"))
		}
		write("<!", string(t), ">")
	default:
		panic(errors.New("encoding of invalid token type"))
	}
}

func (xt XMLTree) encode(w io.Writer) {
	encodeToken(w, xml.StartElement{Name: xt.Name, Attr: xt.Attr})
	if xt.Comment != "" {
		encodeToken(w, xml.Comment(xt.Comment))
	}
	if xt.Directive != "" {
		encodeToken(w, xml.Directive(xt.Directive))
	}
	if xt.Text != "" {
		encodeToken(w, xml.CharData(xt.Text))
	}
	for _, c := range xt.Children {
		c.encode(w)
	}
	encodeToken(w, xml.EndElement{xt.Name})
}

// Output XML string to ``w''.  If ``full'' is true, prepend the standard ProcInst:
//    <?xml version="1.0" encoding="UTF-8"?>
func (xt XMLTree) Encode(w io.Writer, full bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	if full {
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	}
	xt.encode(w)
	return
}
