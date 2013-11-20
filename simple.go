package main

import (
	"net/http"
	"html/template"
	"fmt"
)
func main(){
	http.HandleFunc("/", home)
	http.HandleFunc("/clubs", clubs)
	http.ListenAndServe(":80", nil)
}
func generator(w http.ResponseWriter, fillin string, data map[string]string){
	my_html := template.New("fillin").Funcs(template.FuncMap {
		"XTC": func(x string) template.HTML {
			return template.HTML(x)
	}})
	t := template.Must(my_html.Parse(fillin))
	err := t.Execute(w, data)
	if err != nil {
		fmt.Print("executing template:", err)
	}
}
func clubs(w http.ResponseWriter, r *http.Request){
	generator(w, clubsView(), clubsData())
}
func clubsData()map[string]string{
	return map[string]string{
		"Source": "Clubs page!",
	}
}
func clubsView()string{
	return `<div id=topBar>
		<h1>Clubs</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}`
}
func templator(template string, fillin string, data map[string]string, w http.ResponseWriter){
	switch(template){
	case "home":
		generator(w, homeTemplate(fillin), data)
		break
	case "admin":
		generator(w, adminTemplate(fillin), data)
		break
	case "ajax":
		generator(w, fillin, data)
	}
}
func homeTemplate(body string)string{
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func adminTemplate(body string)string{
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func home(w http.ResponseWriter, r *http.Request){
	templator("home", homeView(), homeData(), w)
}
func homeData()map[string]string{
	return map[string]string{
		"Source": "<select><option>G'day Mate</option></select>",
	}
}
func homeView()string{
	return `<div id=topBar>
		<h1>{{ .PageName}}</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}`
}
