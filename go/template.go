package main

import (
	"html/template"
	"net/http"
	"strings"
)

type menu struct {
	Name, Link string
}

type page struct {
	Title, menu, MenuID string
	Data                M
}

type markupEnv struct {
	JS, CurrentYear, CSS, PNG string
	Page                      page
	Menu                      []menu
	SubMenu                   []menu
}

const (
	titleSeparator              = " - "
	htmlDirectory               = "./h/"
	masterTemplatePath          = htmlDirectory + "_master"
	formsTemplatePath           = htmlDirectory + "forms"
	networkAdaptersTemplatePath = htmlDirectory + "_networkAdapters"
)

var (
	templates      = make(map[string]*template.Template)
	masterTemplate = markupEnv{
		CSS:         dirCSS,
		CurrentYear: currentYear,
		JS:          dirJS,
		Menu:        mainMenu,
		PNG:         dirPNG,
	}
)

func templater(w http.ResponseWriter, page page) {
	//Add HTTP headers so browsers don't cache the HTML resource because it can contain different content every request.
	headers(w, []string{nocache})

	//Add page content just generated to the default page environment (which has CSS and JS, etc).
	masterTemplate.Page = page
	subMenu, ok := subMenus[page.menu]
	if ok {
		masterTemplate.SubMenu = subMenu
	} else {
		masterTemplate.SubMenu = []menu{}
	}

	//Convert page.Title to the HTML template file name (located within htmlDirectory), e.g. Events > events, Club Settings > clubSettings
	pageName := strings.Split(page.Title, titleSeparator)[0]
	pageName = strings.Replace(strings.Title(pageName), " ", "", -1)
	pageName = strings.ToLower(string([]rune(pageName)[0])) + string([]rune(pageName)[1:])

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
			"attr": func(attribute, value string) template.HTMLAttr {
				var output string
				if value != "" {
					output = attribute + "=" + addQuotes(value)
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
		}).ParseFiles(htmlDirectory+pageName, formsTemplatePath, networkAdaptersTemplatePath, masterTemplatePath))
		html = templates[pageName]
	}

	err := html.ExecuteTemplate(w, "master", masterTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
