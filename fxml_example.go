package fxml

import (
	"fmt"
	"os"
)

func ExampleParse() {
	f, err := os.Open("your_file.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	xt, err := Parse(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(xt.ToJSON())
}

func ExampleTraverse() {
	var kml XMLTree //parsed from a KML document
	var pms []XMLTree
	//get all "Placemark"
	kml.Traverse(func(p string, xt XMLTree) bool {
		if xt.Name.Local == "Placemark" {
			pms = append(pms, xt)
		}
		return true
	})
}
