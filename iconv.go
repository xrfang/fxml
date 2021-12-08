//go:build iconv
// +build iconv

package fxml

import (
	iconv "github.com/djimenez/iconv-go"
)

type (
	// Charset conversion interface.  See description in nconv.go
	Converter interface {
		Init(string) error
		ToString([]byte) string
	}
	// "iconv" converter, which output UTF-8 text as-is, and convert other
	// charsets to UTF-8 using libiconv.  This converter uses CGo. See:
	//
	//    https://github.com/djimenez/iconv-go
	IConv struct {
		conv func([]byte) string
		ic   iconv.Converter
	}
)

// Initialize the converter.  This function is called by the parser internally.
func (ic *IConv) Init(charset string) error {
	if charset == "UTF-8" {
		ic.conv = func(v []byte) string { return string(v) }
		return nil
	}
	c, err := iconv.NewConverter(charset, "UTF-8")
	if err != nil {
		return err
	}
	ic.conv = func(v []byte) string {
		output, err := c.ConvertString(string(v))
		if err != nil {
			return err.Error()
		}
		return output
	}
	return nil
}

// Convert a byte slice to string.  This function is called by the parser internally.
func (ic IConv) ToString(v []byte) string {
	return ic.conv(v)
}

var conv IConv
