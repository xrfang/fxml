//go:build !iconv
// +build !iconv

package fxml

import (
	"fmt"
)

type (
	/*
		Charset conversion interface.  This package only outputs UTF-8 text,
		all other charsets are either not accepted or silently converted to
		UTF-8.  The behavior is defined by the actual converter.

		By default, the ``null'' converter is used, which does not accept non
		UTF-8 text, i.e. the ``ProcInst'' of XML document must be:

		    <?xml version="1.0" encoding="UTF-8"?>

		To allow other charsets, use the ``iconv'' converter:

		    go build -tags iconv

		Beware that the ``iconv'' converter uses ``libiconv'' through CGo. See:
		https://github.com/djimenez/iconv-go
	*/
	Converter interface {
		Init(string) error
		ToString([]byte) string
	}
	// "null" converter, which output UTF-8 text as-is, and do not accept
	// non UTF-8 characters.  This converter does not require ``libiconv''.
	//
	// Note: the Converter interface (i.e. this struct) is used internally
	// by the parser, do NOT use it in your code.
	NConv struct{}
)

// Initialize the converter.  This function is called by the parser internally.
func (ic *NConv) Init(charset string) error {
	if charset != "UTF-8" {
		return fmt.Errorf("invalid charset '%s' (UTF-8 only)", charset)
	}
	return nil
}

// Convert a byte slice to string.  This function is called by the parser internally.
func (ic NConv) ToString(v []byte) string {
	return string(v)
}

var conv NConv
