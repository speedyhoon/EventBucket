package main

import (
	"html/template"
	"net/http"
	"strings"
)

type menu struct {
	Name, Link string
	SubMenu    []menu
}

type page struct {
	Title, Menu, MenuID, Heading string
	Data                         M
	Error                        error
	Ajax                         bool
}

type markupEnv struct {
	CurrentYear string
	Page        page
	Menu        []menu
}

const (
	titleSeparator        = " - "
	htmlDirectory         = "./h/"
	masterTemplatePath    = htmlDirectory + "_master"
	formsTemplatePath     = htmlDirectory + "forms"
	reusablesTemplatePath = htmlDirectory + "_reusables"
)

var (
	templates      = make(map[string]*template.Template)
	masterTemplate = markupEnv{
		CurrentYear: currentYear,
		Menu:        mainMenu,
	}
)

func templater(w http.ResponseWriter, page page) {
	//Add HTTP headers so browsers don't cache the HTML resource because it can contain different content every request.
	headers(w, []string{nocache})

	//Convert page.Title to the HTML template file name (located within htmlDirectory), e.g. Events > events, Club Settings > clubSettings
	pageName := strings.Split(page.Title, titleSeparator)[0]
	pageName = strings.Replace(strings.Title(pageName), " ", "", -1)
	pageName = strings.ToLower(string([]rune(pageName)[0])) + string([]rune(pageName)[1:])

	hhh := []string{htmlDirectory + pageName, formsTemplatePath, reusablesTemplatePath}
	if !page.Ajax {
		//Add page content just generated to the default page environment (which has CSS and JS, etc).
		masterTemplate.Page = page
		hhh = append(hhh, masterTemplatePath)
	}

	html, ok := templates[pageName]
	//debug is for dynamically re-parsing (reloading) templates on every request
	if !ok || debug {

		templates[pageName] = template.Must(template.New("main").Funcs(template.FuncMap{
			"hasindex": func(inputs []field, index int) *field {
				if index < len(inputs) && index >= 0 {
					return &inputs[index]
				}
				return nil
			},
			"attr": func(attribute string, value interface{}) template.HTMLAttr {
				var output string
				switch value.(type) {
				case bool:
					if value.(bool) {
						output = attribute
					}
				case string:
					if value.(string) != "" {
						output = attribute + "=" + addQuotes(value.(string))
					}
				}
				return template.HTMLAttr(output)
			},
			"has": func(t interface{}, value string) template.HTMLAttr {
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
				}
				if hasValue {
					return template.HTMLAttr(value)
				}
				return template.HTMLAttr("")
			},
			"nameGrade": func(index uint64) string {
				return dataListGrades()[index].Label
			},
			"nameAgeGroup": func(index uint64) string {
				return dataListAgeGroup()[index].Label
			},
		}).ParseFiles(hhh...))
		html = templates[pageName]
	}

	var err error
	if page.Ajax {
		err = html.ExecuteTemplate(w, "page", page.Data)
	} else {
		err = html.ExecuteTemplate(w, "master", masterTemplate)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
