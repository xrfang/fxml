package fxml

import (
	"bytes"
	"fmt"
	"os"
)

// parse XML data from a reader
func ExampleParse() {
	b := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
	<kml xmlns="http://www.opengis.net/kml/2.2"
	  xmlns:gx="http://www.google.com/kml/ext/2.2"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	  <Document id="example">
		<name>Map of the region</name>
	  </Document>
	</kml>`)
	xt, err := Parse(b)
	if err != nil {
		panic(err)
	}
	xt.Encode(os.Stdout, true)
}

func ExampleTraverse() {
	b := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
	<kml xmlns="http://www.opengis.net/kml/2.2"
	  xmlns:gx="http://www.google.com/kml/ext/2.2"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	  <Document id="example">
		<name>Map of the region</name>
	  </Document>
	</kml>`)
	xt, err := Parse(b)
	if err != nil {
		panic(err)
	}
	xt.Traverse(func(p string, xt XMLTree) bool {
		fmt.Println(p)
		return true
	})
	//Output:
	//kml
	//kml/Document
	//kml/Document/name
}
