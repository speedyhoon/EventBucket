package main

import "html/template"
import "os"

func main() {
	t := template.New("").Funcs(template.FuncMap {
		"unescaped": func(x string) template.HTML {
			return template.HTML(x)
		},
	})
	t, _ = t.Parse(`{{.}}, {{unescaped .}}`)
	t.Execute(os.Stdout, "Might <Escape> That")
}
