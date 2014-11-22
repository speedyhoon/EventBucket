package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/template"
)

var data = []byte(`[
	{"name": "ANisus", "sex":"M"},
	{"name": "Sofia", "sex":"F"},
	{"name": "Anna", "sex":"F"}
]`)
var tmpl = `
{{$hasFemale := cell ""}}
{{range .}}
	{{if ne .sex $hasFemale.Get}}
		Different
		{{$_ := $hasFemale.Set .sex}}
	{{end}}
	{{.name}}
{{end}}
`

func main() {
	var v []interface{}
	json.Unmarshal(data, &v) // Hehe, ignored the error

	var t = template.New("t")

	_, err := t.Funcs(template.FuncMap{"eq": eq, "cell": NewCell}).Parse(tmpl)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, v)

	if err != nil {
		panic(err)
	}
}

type Cell struct{ v interface{} }

func NewCell(v ...interface{}) (*Cell, error) {
	switch len(v) {
	case 0:
		return new(Cell), nil
	case 1:
		return &Cell{v[0]}, nil
	default:
		return nil, fmt.Errorf("wrong number of args: want 0 or 1, got %v", len(v))
	}
}

func (c *Cell) Set(v interface{}) *Cell { c.v = v; return c }
func (c *Cell) Get() interface{}        { return c.v }

func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}
