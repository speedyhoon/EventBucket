package main

import (
	"log"
	"html/template"//"text/template"
	"net/http"
	"fmt"
	"strings"
	"bytes"
)
//SETTINGS
const port string = ":80"//FEATURE add setting to change the port number

func router() map[string]func(http.ResponseWriter){
	return map[string]func(http.ResponseWriter){
		"": controller,
		"t": home,
	}
}

func main() {
	http.HandleFunc("/", server)
	http.ListenAndServe(port, nil)
}

//Server() handles http requests. It checks the requested page exists, before passing the request to controller()
func server(w http.ResponseWriter, r *http.Request){
	fmt.Print("\nRequestURI=\t"+r.RequestURI)//r.URL.Path
	//TODO ability to handle /favicon.ico
	//research returning .htm files within directories
	tempUrl := strings.Trim(strings.ToLower(r.URL.Path), "/")
	fmt.Print("\ntempUrl=\t"+tempUrl)
	if "/"+tempUrl != r.URL.Path {
		http.Redirect(w,r,"http://localhost/"+tempUrl, 301)//research would like to use 308 Permanent Redirect but GO makes it 301 instead
	}
	if model_func, ok := router()[tempUrl]; ok {
		model_func(w)
	}else{
		//TODO return a 404 page and display any similar pages from the DB
	}
}

func controller(w http.ResponseWriter){
	renderThese := settingsPage()
	outputHtml := ""
	for _, temp2 := range renderThese{
		outputHtml += generator2(pane(), temp2())
	}
	pageSize := map[string]string{
		"Css": "Css",
		"Ico": "Ico",
		"IcoType": "IcoType",
		"Title": "Title",
		"PageName": "PageName",
		"Menu": "Menu",
		"Source": outputHtml,
		"Js": "Js",
	}
	generator3(w, page(), pageSize)
}

func settingsPage() map[string]func() map[string]string{
	return map[string]func() map[string]string{
		"date": date,
		"ranges": ranges,
		"grades": grades,
	}
}

func home(w http.ResponseWriter){
	fmt.Print("inside home= "+"input")
}
func generator3(w http.ResponseWriter, letter string, recipients  map[string]string){
	myhtml := template.New("letter").Funcs(template.FuncMap {
		"XTC": func(x string) template.HTML {
			return template.HTML(x)
		},
	})
	t := template.Must(myhtml.Parse(letter))
	err := t.Execute(w, recipients)
	if err != nil {
		log.Println("executing template:", err)
	}
}
func generator2(letter string, recipients map[string]string)string{
	myhtml := template.New("letter").Funcs(template.FuncMap {
		"XTC": func(x string) template.HTML {
			return template.HTML(x)
		},
	})
	t := template.Must(myhtml.Parse(letter))
	var doc bytes.Buffer
	err := t.Execute(&doc, recipients)
	output := doc.String()
	if err != nil {
		log.Println("executing template:", err)
	}
	return output
}




func date() map[string]string{
	return map[string]string{
		"Source": "date html goes here!!",
	}
}
func ranges() map[string]string{
	return map[string]string{
		"Source": "<input value=Boo>",
	}
}
func grades() map[string]string{
	return map[string]string{
		"Source": "<select><option>G'day Mate</option></select>",
	}
}

func page()string{
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>{{.Title}}</title>
</head>
<body>
	<div id=topBar>
		<h1>{{ .PageName}}</h1> {{XTC .Menu}}
	</div>
	{{XTC .Source}}
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}

//Panes display a large heading with content in the trailing div and sections follow with sub headings as <h3> tags
func pane()string{
	return `<h2>{{.Title}}</h2>
<div>
	{{XTC .Source}}
</div>`
}

//Sections MUST follow below sections
func section()string{
	return `<div>
	<h3>{{.Title}}</h3>
	{{XTC .Source}}
</div>`
}

func minify(input string) string{
	//input = Replace(input, "  ", " ", -1))
	//input = Replace(input, "	", "", -1))
	//input = Replace(input, "\n", "", -1))
	//TODO:: remove all unicode chars above 255
	return input
}
