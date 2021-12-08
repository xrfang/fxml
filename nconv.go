//go:build !iconv
// +build !iconv

package fxml

import (
	"fmt"
)

type (
	Converter interface {
		Init(string) error
		ToString([]byte) string
	}
	NConv struct{}
)

func (ic *NConv) Init(charset string) error {
	if charset != "UTF-8" {
		return fmt.Errorf("invalid charset '%s' (UTF-8 only)", charset)
	}
	return nil
}

func (ic NConv) ToString(v []byte) string {
	return string(v)
}

var conv NConv
