package main

import (
	"net/http"
	"fmt"
	"reflect"
	"strings"
)
//TODO somehow add the option to build the offline & online version using the same code. But the offline version is seriously limited
//SETTINGS
const(
	PRODUCTION = false
	DEBUG = true

	PORT string = "80"//FEATURE add setting to change the port number

	EVERYONE = 0//group permissions
	MEMBER = 10
	ORGANISER = 25
	ADMIN = 50
	ULTRA = 99

	NODATA = false
	DATA = true
)

//TODO save router() as a global variable so that it only needs to execute once!
//TODO initalise the database so the db abstraction layer is globally accessible
func main(){
	http.HandleFunc("/", server)
	http.ListenAndServe(":"+PORT, nil)
}
           /*
request comes in and pass the func REQUEST Array
map the request to a function
GET - no parameters
GET with parameters
POST with parameters
TEMPLATE to render in
func options - to make them more reusable
refresh = "funcName" redirect to another page once completed or False to prevent redirect
update = map of fields to update provided from client - makes eventUpdate very dynamic
range = "up" OR "down" makes classRange very dynamic
class = "up" OR "down" makes classRange very dynamic



controller(http.request, permission_required, page_model_to_call, page_model_parameters )
      */
//Server() handles http requests. It checks the requested page exists, before passing the request to controller()
func server(w http.ResponseWriter, r *http.Request){
	fmt.Print("\nRequestURI=\t"+r.RequestURI)//r.URL.Path
	//TODO ability to handle /favicon.ico
	//TODO research returning .htm files within directories
	tempUrl := strings.Trim(strings.ToLower(r.URL.Path), "/")
	fmt.Print("\ntempUrl=\t"+tempUrl)
	if "/"+tempUrl != r.URL.Path {
		http.Redirect(w,r,"http://localhost/"+tempUrl, 301)//research would like to use 308 Permanent Redirect but GO makes it 301 instead
	}
	//TODO some how pass to the model the http.Request so it knows what template to use based on GET, GET() or POST()
	//TODO add in the GET or POST request parameters
	if page_func, ok := router()[tempUrl]; ok {
		ioc_new(page_func)
//		controller(w, page_func)
	}else{
		//TODO return a 404 page and display any similar pages from the DB
	}
}
func Call(name interface{}, params ... interface{}) map[reflect.Value][]reflect.Value{// (result []reflect.Value, err error)
	f := reflect.ValueOf(name)
	if len(params) != 1+f.Type().NumIn() {
		fmt.Print("whoops there was an error! The number of parameters do not match!")
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return map[reflect.Value][]reflect.Value{
		f: in,
	}
}

func ioc_new(model map[reflect.Value][]reflect.Value){
	var user_perm int
	user_perm = user_permission()
//	httpData := "data to do something with"
	for function, parameters := range model{
		//TODO it would be nice to call the needed functions in order before calling the model I'm note sure if this is possible though
		params := parameters[2:len(parameters)]
//		temp := []reflect.Value{}
		fmt.Print("\n\n", len(parameters), "<<length\n\n")

		if(user_perm >= parameters[0].Int){

//			if(user_perm > EVERYONE){
//				//TODO add user_perm to the first parameter
//			}
//			if(parameters[1] == DATA){
//				//TODO add httpData to the function call		parameters[1] = httpData
//			}
			//TODO append parameters[2:len(parameters)] to the end
			//TODO need to add controller as a layer between the model and the server/ioc!!!
			function.Call(params)
		}
	}
}
//func reflectSwitch(val reflect.Value)(interface{}){
//	switch val.Kind(){
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			m[typeField.Name] = strconv.FormatInt(val.Int(), 10)
//		case reflect.String:
//			m[typeField.Name] = val.String()
//	}
//	return
//}

func user_permission()int{
	//TODO set the users permission level when they login
	return 0
}

func router()(map[string]map [reflect.Value][]reflect.Value){
	//TODO this is called every request. It should only be called once when starting the server
	elements := map[string]map[reflect.Value][]reflect.Value{
		"foo": Call(foo,		ULTRA,		NODATA	),
		"bar": Call(bar,		EVERYONE,	DATA,		1,4,6),
	}
	return elements
}

func foo(user int) {
	fmt.Print("\nwe are running foo. User=", user)
}

func bar(data string, a, b, c int) {
	fmt.Print("\nwe are running bar; ", data, a, b, c)
}
































func controller(http.request, page_model_to_call, page_model_parameters )
//func controller(w http.ResponseWriter){
func controller(w http.ResponseWriter, modelfunc func(http.ResponseWriter)){
	//TODO should use buffers rather than ResponseWriters as they can be passed to minify easier!!!
	Css,Ico,IcoType,Title,PageName,Menu,outputHtml,Js := ioc(w,modelfunc)
//	renderThese := settingsPage()
//	outputHtml := ""
//	for _, temp2 := range renderThese{
//		outputHtml += generator2(pane(), temp2())
//	}
	pageSize := map[string]string{
		"Css": Css,
		"Ico": Ico,
		"IcoType": IcoType,
		"Title": Title,
		"PageName": PageName,
		"Menu": Menu,
		"Source": outputHtml,
		"Js": Js,
	}
	minify(generator3(w, page(), pageSize))
}

                 /*

func ioc(w http.ResponseWriter,modelfunc func() map[string]func() map[string]string)(string,string,string,string,string,string,string,string){
	Css := "ioc.css"
	Ico := "icon.png"
	IcoType := "image/png"//TODO make this automatic?
	Title := "New Title :)"
	PageName := "My New Page!"
	Menu := "Menu items go here?"
	Js := "newJSfile.js"
	renderThese := modelfunc(w)
//	renderThese := settingsPage()
	outputHtml := ""
	for _, temp2 := range renderThese{
		outputHtml += generator2(pane(), temp2())
	}
	return Css,Ico,IcoType,Title,PageName,Menu,outputHtml,Js
}
func generator3(w http.ResponseWriter, letter string, recipients  map[string]string){
	myhtml := template.New("letter").Funcs(template.FuncMap {
		//TODO when x is nil still return an empty template instead of silently failing
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
         */
func home(w http.ResponseWriter){
	fmt.Print("inside home= "+"input")
}

func settingsPage() map[string]func() map[string]string{
	return map[string]func() map[string]string{
		"date": date,
		"ranges": ranges,
		"grades": grades,
	}
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
	<title>{{if .Title}}{{.Title}}{{else}}EventBucket{{end}}</title>
</head>
<body>
	<div id=topBar>
		<h1>{{ .PageName}}</h1>{{if .Menu}} {{XTC .Menu}}{{end}}
	</div>
	{{if .Source}}{{XTC .Source}}{{end}}
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
	if !DEBUG{
		//input = Replace(input, "  ", " ", -1))
		//input = Replace(input, "	", "", -1))
		//input = Replace(input, "\n", "", -1))
		//TODO:: remove all unicode chars above 255
	}
	return input
}
