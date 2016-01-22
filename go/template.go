package main

import (
	"html/template"
	"net/http"
	"strings"
)

const (
	titleSeparator              = " - "
	masterTemplatePath          = "./h/_master"
	formsTemplatePath           = "./h/forms"
	networkAdaptersTemplatePath = "./h/networkAdapters"
)

var (
	templates      = make(map[string]*template.Template)
	masterTemplate = Template{
		CSS:         "dirCss",
		CurrentYear: currentYear,
		JS:          "dirJS",
		Menu: []menu{{
			Name: "Home",
			Link: urlHome,
		}, {
			Name: "Archive",
			Link: urlArchive,
		}, {
			Name: "Clubs",
			Link: urlClubs,
		}, {
			Name: "Events",
			Link: urlEvents,
		}, {
			Name: "About",
			Link: urlAbout,
		}, {
			Name: "Shooters",
			Link: urlShooters,
		}, {
			Name: "Licence",
			Link: urlLicence,
		}},
		PNG: "dirPNG",
	}
	masterTemplate2 = Template2{
		CSS:         "dirCss",
		CurrentYear: currentYear,
		JS:          "dirJS",
		Menu: []menu{{
			Name: "Home",
			Link: urlHome,
		}, {
			Name: "Archive",
			Link: urlArchive,
		}, {
			Name: "Clubs",
			Link: urlClubs,
		}, {
			Name: "Events",
			Link: urlEvents,
		}, {
			Name: "About",
			Link: urlAbout,
		}, {
			Name: "Shooters",
			Link: urlShooters,
		}, {
			Name: "Licence",
			Link: urlLicence,
		}},
		PNG: "dirPNG",
	}
)

type menu struct {
	Name, Link string
}

type page struct {
	Title string
	Data  M
}

type page2 struct {
	Title string
	Data  interface{}
}

type Template struct {
	JS, CurrentYear, CSS, PNG string
	Page                      page
	Menu                      []menu
}

type Template2 struct {
	JS, CurrentYear, CSS, PNG string
	Page                      page2
	Menu                      []menu
}

func templater(w http.ResponseWriter, page page) {

	headers(w, []string{nocache})

	pageName := strings.Split(strings.ToLower(page.Title), titleSeparator)[0]
	masterTemplate.Page = page

	//	funcs := template.New(viewController.TemplateFile + "Template").Funcs(template.FuncMap{
	//	"HTM": func(x string) template.HTML {
	//		return template.HTML(x)
	//	})

	html, ok := templates[pageName]
	if !ok || debug { //debug is for dynamically reloading templates on every request

		//		tpl := template.Must(template.New("main").Funcs(funcMap).ParseGlob("*.html"))

		//		templ, err := template.ParseFiles("../htm/"+pageName+".htm", masterTemplatePath)

		templates[pageName] = template.Must(template.New("main").Funcs(template.FuncMap{
			"hasindex": func(inputs []field, index int) *field {
				if index < len(inputs) && index >= 0 {
					return &inputs[index]
				}
				return nil
			},
			"attr": func(attribute, value string) template.HTMLAttr {
				var output string
				if value != "" {
					output = attribute + "=" + addQuotes(value)
				}
				return template.HTMLAttr(output)
			},
			/*"in": func(t interface{}) bool {
				switch t.(type) {
				default:
					warn.Printf("unexpected type %T", t) // %T prints whatever type t has
					return false
				case []option:
					return len(t.([]option)) >= 1
				case string:
					return t != ""
				}
			},*/

			"has": func(t interface{}, value string) template.HTMLAttr {
				//				t = functionOfSomeType()
				var hasValue bool
				switch t.(type) {
				default:
					warn.Printf("unexpected type %T", t) // %T prints whatever type t has
				case []option:
					hasValue = len(t.([]option)) >= 1
				case string:
					hasValue = t != ""
				case bool:
					hasValue = t.(bool)
					//fmt.Printf("boolean %t\n", t) // t has type bool
				}

				//				if len(condition) > 0 {
				if hasValue {
					//				if value == "" {
					//					return template.HTMLAttr("")
					//				}
					return template.HTMLAttr(value)
				}
				return template.HTMLAttr("")
				//				return template.HTMLAttr(value)
			},

			"withhas": func(t interface{}, value string) template.HTMLAttr {
				//				t = functionOfSomeType()
				var hasValue bool
				switch t.(type) {
				default:
					warn.Printf("unexpected type %T", t) // %T prints whatever type t has
				case []option:
					hasValue = len(t.([]option)) >= 1
				case string:
					hasValue = t != ""
				case int:
					hasValue = t.(int) > 0
					//fmt.Printf("boolean %t\n", t) // t has type bool
				}

				//				if len(condition) > 0 {
				if hasValue {
					//				if value == "" {
					//					return template.HTMLAttr("")
					//				}
					return template.HTMLAttr(value)
				}
				return template.HTMLAttr("")
				//				return template.HTMLAttr(value)
			},
			/*"might": func(value string, inputs []input, index int, t interface{}) template.HTMLAttr {
				//				var hasValue bool
				if len(inputs) <= index {
					//					variable := item[t]
					return template.HTMLAttr(value)
				}
				return template.HTMLAttr("")

				/*	for _, variable := range t {
					switch variable.(type) {
					default:
						warn.Printf("unexpected type %T", t) // %T prints whatever type t has
					case []option:
						hasValue = len(variable.([]option)) >= 1
					case string:
						hasValue = variable != ""
						//fmt.Printf("boolean %t\n", t) // t has type bool
					}
					if !hasValue {
						return template.HTMLAttr("")
					}
				}* /
				//				if len(condition) > 0 {
				//				if hasValue {
				//				if value == "" {
				//					return template.HTMLAttr("")
				//				}
				//				return template.HTMLAttr(value)
				//				}
				//				return template.HTMLAttr("")
				//				return template.HTMLAttr(value)
			},*/
			/*"D": func(value string) template.HTMLAttr {
				return template.HTMLAttr(addQuotes(value))
			},*/
			//			"field": func(fields []Field, index int) Field {
			//				fmt.Println(fields[index])
			//				return fields[index]
			//			},
			/*"HTM": func(inner string) template.HTML {
				return template.HTML(inner)
			},*/
			//			"nice": func(inp Jjj) string {
			//				return "fdfd"
			//			},
		}).ParseFiles("./h/"+pageName, formsTemplatePath, networkAdaptersTemplatePath, masterTemplatePath))
		html = templates[pageName]
	}

	err := html.ExecuteTemplate(w, "master", masterTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func templater2(w http.ResponseWriter, page page2) {
	pageName := strings.Split(strings.ToLower(page.Title), titleSeparator)[0]
	masterTemplate2.Page = page

	//	funcs := template.New(viewController.TemplateFile + "Template").Funcs(template.FuncMap{
	//	"HTM": func(x string) template.HTML {
	//		return template.HTML(x)
	//	})

	html, ok := templates[pageName]
	if !ok || debug { //debug is for dynamically reloading templates on every request

		templ, err := template.ParseFiles("../htm/"+pageName+".htm", masterTemplatePath)

		templates[pageName] = template.Must(templ.Funcs(template.FuncMap{

			"attr": func(attribute, value string) template.HTMLAttr {
				//				if value == "" {
				//					return template.HTMLAttr("")
				//				}
				return template.HTMLAttr(attribute + "=" + addQuotes(value))
				//				return template.HTMLAttr(value)
			},
		}), err)
		html = templates[pageName]
	}

	err := html.ExecuteTemplate(w, "master", masterTemplate2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
