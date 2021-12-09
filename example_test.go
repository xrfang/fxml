package fxml_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/xrfang/fxml"
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
	xt, err := fxml.Parse(b)
	if err != nil {
		panic(err)
	}
	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "    ")
	je.Encode(xt)
	//Output:
	// {
	// 	"Name": {
	// 		"Space": "",
	// 		"Local": "kml"
	// 	},
	// 	"Attr": [
	// 		{
	// 			"Name": {
	// 				"Space": "",
	// 				"Local": "xmlns"
	// 			},
	// 			"Value": "http://www.opengis.net/kml/2.2"
	// 		},
	// 		{
	// 			"Name": {
	// 				"Space": "xmlns",
	// 				"Local": "gx"
	// 			},
	// 			"Value": "http://www.google.com/kml/ext/2.2"
	// 		},
	// 		{
	// 			"Name": {
	// 				"Space": "xmlns",
	// 				"Local": "xsi"
	// 			},
	// 			"Value": "http://www.w3.org/2001/XMLSchema-instance"
	// 		}
	// 	],
	// 	"Text": "\n\t",
	// 	"Children": [
	// 		{
	// 			"Name": {
	// 				"Space": "",
	// 				"Local": "Document"
	// 			},
	// 			"Attr": [
	// 				{
	// 					"Name": {
	// 						"Space": "",
	// 						"Local": "id"
	// 					},
	// 					"Value": "example"
	// 				}
	// 			],
	// 			"Text": "\n\t  ",
	// 			"Children": [
	// 				{
	// 					"Name": {
	// 						"Space": "",
	// 						"Local": "name"
	// 					},
	// 					"Text": "Map of the region"
	// 				}
	// 			]
	// 		}
	// 	]
	// }
}

// render XMLTree as XML string
func ExampleXMLTree_Encode() {
	b := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
	<kml xmlns="http://www.opengis.net/kml/2.2"
	  xmlns:gx="http://www.google.com/kml/ext/2.2"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	  <Document id="example">
		<name>Map of the region</name>
	  </Document>
	</kml>`)
	xt, err := fxml.Parse(b)
	if err != nil {
		panic(err)
	}
	xt.Encode(os.Stdout, true)
	//Output: <?xml version="1.0" encoding="UTF-8"?><kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	//&#x9;<Document id="example">
	//&#x9;  <name>Map of the region</name></Document></kml>
}

// travers a XMLTree
func ExampleXMLTree_Traverse() {
	b := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
	<kml xmlns="http://www.opengis.net/kml/2.2"
	  xmlns:gx="http://www.google.com/kml/ext/2.2"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	  <Document id="example">
		<name>Map of the region</name>
	  </Document>
	</kml>`)
	xt, err := fxml.Parse(b)
	if err != nil {
		panic(err)
	}
	xt.Traverse(func(p string, xt fxml.XMLTree) bool {
		fmt.Println(p)
		return true
	})
	//Output:
	//kml
	//kml/Document
	//kml/Document/name
}
