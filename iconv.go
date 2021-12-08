//go:build iconv
// +build iconv

package fxml

import (
	iconv "github.com/djimenez/iconv-go"
)

type (
	Converter interface {
		Init(string) error
		ToString([]byte) string
	}
	IConv struct {
		conv func([]byte) string
		ic   iconv.Converter
	}
)

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

func (ic IConv) ToString(v []byte) string {
	return ic.conv(v)
}

var conv IConv
