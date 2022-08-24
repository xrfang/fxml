package fxml_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	//{
	//    "Name": {
	//        "Space": "",
	//        "Local": "kml"
	//    },
	//    "Attr": [
	//        {
	//            "Name": {
	//                "Space": "",
	//                "Local": "xmlns"
	//            },
	//            "Value": "http://www.opengis.net/kml/2.2"
	//        },
	//        {
	//            "Name": {
	//                "Space": "xmlns",
	//                "Local": "gx"
	//            },
	//            "Value": "http://www.google.com/kml/ext/2.2"
	//        },
	//        {
	//            "Name": {
	//                "Space": "xmlns",
	//                "Local": "xsi"
	//            },
	//            "Value": "http://www.w3.org/2001/XMLSchema-instance"
	//        }
	//    ],
	//    "Text": "\n\t",
	//    "Children": [
	//        {
	//            "Name": {
	//                "Space": "",
	//                "Local": "Document"
	//            },
	//            "Attr": [
	//                {
	//                    "Name": {
	//                        "Space": "",
	//                        "Local": "id"
	//                    },
	//                    "Value": "example"
	//                }
	//            ],
	//            "Text": "\n\t  ",
	//            "Children": [
	//                {
	//                    "Name": {
	//                        "Space": "",
	//                        "Local": "name"
	//                    },
	//                    "Text": "Map of the region"
	//                }
	//            ]
	//        }
	//    ]
	//}
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
	//Output: <?xml version="1.0" encoding="UTF-8"?><kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><Document id="example"><name>Map of the region</name></Document></kml>
}

// render XMLTree as XML string
func ExampleXMLTree_ToString() {
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
	str, err := xt.ToString(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
	//Output: <?xml version="1.0" encoding="UTF-8"?><kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><Document id="example"><name>Map of the region</name></Document></kml>
}

// traverse a XMLTree
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

// walk part of a XMLTree and modify node text
func ExampleXMLTree_Walk() {
	bs := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
	<style>This is root
	  <LineStyle>
        <color>red</color>
        <width>1</width>
	  </LineStyle>
	  <PolyStyle>
		<color>green</color>
		<width>2</width>
	  </PolyStyle>
	</style>`)
	xt, _ := fxml.Parse(bs)
	xt.Walk(func(ni fxml.XNodInfo, x *fxml.XMLTree) fxml.XWalkResult {
		fmt.Printf("%s|%+v\n", strings.Join(ni.Path, "/"), ni)
		x.Text = strings.ToUpper(x.Text)
		if ni.Path[len(ni.Path)-1] == "color" {
			return fxml.WRSkip
		}
		if ni.Path[len(ni.Path)-1] == "PolyStyle" {
			return fxml.WRTerm
		}
		return fxml.WRCont
	})
	je := json.NewEncoder(os.Stdout)
	je.Encode(xt)
	//Output:
	//style|{Path:[style] Index:0 RIndex:0}
	//style/|{Path:[style ] Index:0 RIndex:-3}
	//style/LineStyle|{Path:[style LineStyle] Index:1 RIndex:-2}
	//style/LineStyle/color|{Path:[style LineStyle color] Index:0 RIndex:-2}
	//style/PolyStyle|{Path:[style PolyStyle] Index:2 RIndex:-1}
	//{"Name":{"Space":"","Local":"style"},"Children":[{"Name":{"Space":"","Local":""},"Text":"THIS IS ROOT"},{"Name":{"Space":"","Local":"LineStyle"},"Children":[{"Name":{"Space":"","Local":"color"},"Text":"RED"},{"Name":{"Space":"","Local":"width"},"Text":"1"}]},{"Name":{"Space":"","Local":"PolyStyle"},"Children":[{"Name":{"Space":"","Local":"color"},"Text":"green"},{"Name":{"Space":"","Local":"width"},"Text":"2"}]}]}
}
