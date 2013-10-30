package main

import (
	"log"
	"os"
	"text/template"
	"net/http"
	"fmt"
	"strings"
)

type Page struct {
	Name, Gift string
	Attended   bool
}

func main() {
	http.HandleFunc("/", server())
	http.ListenAndServe(":8080", nil)
}

func server() http.HandlerFunc{
	//fmt.Print(http.Request)
	//fmt.Printf("[%q]", strings.Trim(" !!! Achtung !!! ", "! "))
	return func(w http.ResponseWriter, r *http.Request) {
		tempUrl := strings.ToLower(r.URL.Path)
		fmt.Print("\nRequestURI=\t"+r.RequestURI)
		fmt.Print("\nURL.Path=\t"+r.URL.Path)
		if tempUrl != r.URL.Path {
			fmt.Print("\nredirect user to lowercase path instead\n")
			//		return http.RedirectHandler("http://localhost:8080/nice", 301)
		}
		//	//tempUrl = fmt.Sprintf("%b", tempUrl)
		//	//fmt.Printf(%q\n, )
		//	tempUrl = strings.Trim(tempUrl, "/")
		//	fmt.Print(tempUrl +"\n")
		//	//fmt.Print("\n")
		//	//err := templates.ExecuteTemplate(w, tmpl+".html", p)
		//	//if err != nil {
		//	//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	//}
	}
}
func controller(){
	letter := htmltemp()
	recipients := model()
	generater(letter, recipients )
}

func generater(letter string, recipients []Page){
	t := template.Must(template.New("letter").Parse(letter))

	for _, r := range recipients {
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}

func model() ([]Page){
	return []Page{
		{"Aunt Mildred", "bone china tea set", true},
		{"Uncle John", "moleskin pants", false},
		{"Cousin Rodney", "", false},
	}
}

func templater() string{
	return `
Dear {{.Name}},
{{if .Attended}}
It           was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`
}

func htmltemp() string{
	return `
<!doctype html>
<html>
<head>
  <title>EventBucket{{if .Name}} - {{.Name}}{{else}} Gooo Web Framework v0.1{{end}}</title>
</head>
<body>
{{if .Attended}}
It        was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
</body>
</html>
`
}

func minify(input string) string{
	//input = Replace(input, "  ", " ", -1))
	//input = Replace(input, "	", "", -1))
	//input = Replace(input, "\n", "", -1))
	//TODO:: remove all unicode chars above 255
	return input
}


//THIS WORKS
//test := hello
//fmt.Println(test("Leesa"))
//func hello(input string) string{
//	return "Hello "+input+"!!!!!"
//}
