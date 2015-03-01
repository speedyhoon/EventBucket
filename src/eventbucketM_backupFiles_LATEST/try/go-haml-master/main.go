package main

import (
  "github.com/realistschuckle/gohaml"
  "io/ioutil"
  "net/http"
  "text/template"
  "fmt"
)

//type Person struct {
//    Name string
//}

//var welcomeTemplate = makeTemplate()

func main() {
  fmt.Println("Please access to http://localhost:8080/ by browser.")
  http.HandleFunc("/", handleRoot)
  http.ListenAndServe("localhost:8080", nil)
}

func makeTemplate() *template.Template {
  scope := map[string]interface{}{"lang": "HAML"}
  content, _ := ioutil.ReadFile("sample.haml")
  engine, _ := gohaml.NewEngine(string(content))
  output := engine.Render(scope)

  welcomeTemplate := template.Must(template.New("").Parse(output))
  return welcomeTemplate
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
//  p := Person { Name: "@dddaisuke" }

	info := map[string]interface{}{
		"DIR_CSS": "_/c/",
		"FAVICON": "_/p/a",
		"Title": "Temp Page",
		"Menu": "Menu options",
		"Body": "Body Contents",
		"CURRENT_YEAR": "2014",
	}



//  welcomeTemplate.Execute(w, p)	//cache haml file
	makeTemplate().Execute(w, info)
}
