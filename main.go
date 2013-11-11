package main

import (
	"log"
	"os"
//	"text/template"
	"html/template"
	"net/http"
	"fmt"
	"strings"
	"bytes"
)

type Page struct {
	Css, Ico, IcoType, Title, PageName, Menu, Source, Js string
}

func runit() map[string]func(http.ResponseWriter,string){
	return map[string]func(http.ResponseWriter,string){
		"": controller,
		"t": home,
	}
}

func main() {
	http.HandleFunc("/", server())
	//FEATURE add setting to change the port number
	http.ListenAndServe(":8080", nil)
}

func redirectHandler(path string, referer string) func(http.ResponseWriter, *http.Request) {
	fmt.Print("\nRedirect user to lowercase path instead\n")
	return func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
}

//Server() handles http requests. It checks the requested page exists,
//before passing the request to controller()
//func server(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
func server() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
fmt.Print("\nRequestURI=\t"+r.RequestURI)
fmt.Print("\nURL.Path=\t"+r.URL.Path)
		//TODO ability to handle /favicon.ico
		//research returning .htm files within directories
		tempUrl := strings.ToLower(r.URL.Path)
		if tempUrl != r.URL.Path {
			redirectHandler(tempUrl, "")
		}
		tempUrl = strings.Trim(tempUrl, "/")
		fmt.Print("\nTempUrl=\t"+tempUrl)
		model_func, ok := runit()[tempUrl];
		if ok {
			model_func(w, "hardcoded value")
		}else{
			//TODO return a 404 page and display any similar pages from the DB
			redirectHandler("404", tempUrl)
		}
	}

}


func controller(w http.ResponseWriter, input string){
	renderThese := settingsPage()
	outputHtml := ""
	for _, temp2 := range renderThese{
//		generateViews(pane(), temp2())
//		renderTemplate(w http.ResponseWriter, pane(), temp2())
		outputHtml += generator2(pane(), temp2())
//		outputHtml += pane(temp2, title())
	}

	pageSize := Page{
		Css: "Css",
		Ico: "Ico",
		IcoType: "IcoType",
		Title: "Title",
		PageName: "PageName",
		Menu: "Menu",
		Source: outputHtml,
		Js: "Js",
	}
	generator3(w, page(), pageSize)
	fmt.Print(outputHtml)
//	letter := htmltemp()
//	recipients := settingsPage()
//	generater(letter, recipients )
}

func settingsPage() map[string]func()View{
	return map[string]func()View{
		"date": date,
		"ranges": ranges,
		"grades": grades,
	}
}
func generateViews(letter string, recipients View){
	t := template.Must(template.New("letter").Parse(letter))
	err := t.Execute(os.Stdout, recipients)
	if err != nil {
		log.Println("executing template:", err)
	}
}





func home(w http.ResponseWriter, input string){
	fmt.Print("inside home= "+input)
}
func generator3(w http.ResponseWriter, letter string, recipients Page)string{
	myhtml := template.New("letter").Funcs(template.FuncMap {
		"XTC": func(x string) template.HTML {
			return template.HTML(x)
		},
	})
	t := template.Must(myhtml.Parse(letter))
//	t := template.New("letter").Parse(letter)
	err := t.Execute(w, recipients)
	if err != nil {
		log.Println("executing template:", err)
	}
	return "output?"
}
func generator2(letter string, recipients View)string{
	t := template.Must(template.New("letter").Parse(letter))
	var doc bytes.Buffer
	err := t.Execute(&doc, recipients)
	output := doc.String()
	if err != nil {
		log.Println("executing template:", err)
	}
	return output
}
func generator(letter string, recipients []Page){
	t := template.Must(template.New("letter").Parse(letter))
	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}
















type View struct {
	Title, Source string
}
func date() View{
	return View{
		"Date-Time","date html goes here!!",
	}
}
func ranges() View{
	return View{
		"Ranges","<input value=Boo>",
	}
}
func grades() View{
	return View{
		"Grades","<select><option>G'day Mate</option></select>",
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
		<h1>{{XTC .PageName}}</h1> {{XTC .Menu}}
	</div>
	{{XTC .Source}}
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
//{{printf "%s" .Source}}

//Panes display a large heading with content in the trailing div
//and sections follow with sub headings as <h3> tags
//func pane(title, source string)string{
func pane()string{
//return `
//<h2>` + title + `</h2>
//<div>
//	` + source + `
//</div>
//`
return `<h2>{{.Title}}</h2>
<div>
	{{XTC .Source}}
</div>`
//	{{printf "%s" .Source}}
}

//Sections MUST follow below sections
//func section(title, source string)string{
func section()string{
//return `<div>
//	<h3>` + title + `</h3>
//	` + source + `
//</div>
//`
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
