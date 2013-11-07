package main

import (
	"fmt"
	"html/template"
	"os"
)
//
//type Inventory struct {
//	Material string
//	Count    uint
//}

func main() {
	// prep the template
	tmpl, err := template.New("test").Parse(
		"{{.Count}} items are made of {{.Material}} - {{.Foo}}\n"
	)
	if err != nil {
		panic(err)
	}

	// map first
	sweaterMap := map[string]string{"Count": "17", "Material": "wool"}
	err = tmpl.Execute(os.Stdout, sweaterMap)
	if err != nil {
		fmt.Println("Error!")
		fmt.Println(err)
	}

	// struct second
//	sweaters := Inventory{"wool", 17}
//	err = tmpl.Execute(os.Stdout, sweaters)
//	if err != nil {
//		fmt.Println("Error!")
//		fmt.Println(err)
//	}
}
