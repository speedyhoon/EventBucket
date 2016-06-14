package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type menu struct {
	Name, Link        string
	SubMenu           []menu
	RangeMenu, Hidden bool
}

type page struct {
	Title, Menu, MenuID, Heading string
	Data                         map[string]interface{}
	Error                        error
	template                     uint8
	JS                           []string
}

type markupEnv struct {
	//CurrentYear string
	Page  page
	Menu  []menu
	Theme bool
}

const (
	htmlDirectory         = "./h/"
	masterTemplatePath    = htmlDirectory + "master"
	masterScoreboard      = htmlDirectory + "masterScoreboard"
	formsTemplatePath     = htmlDirectory + "forms"
	reusablesTemplatePath = htmlDirectory + "reusables"

	templateScoreboard = 1
	templateNone       = 255
)

var (
	masterStuff = [][]string{
		{formsTemplatePath, reusablesTemplatePath, masterTemplatePath},
		{masterScoreboard},
	}
	templates      = make(map[string]*template.Template)
	masterTemplate = markupEnv{
		Menu: []menu{{
			Name: "Events",
			Link: urlEvents,
			SubMenu: []menu{{
				Name: "Entries",
				Link: urlEntries,
			}, {
				Name: "Event Settings",
				Link: urlEventSettings,
			}, {
				Name:      "Scoreboard",
				Link:      urlScoreboard,
				RangeMenu: true,
			}, {
				Name:      "Enter Shots",
				Link:      urlEnterShots,
				RangeMenu: true,
			}, {
				Name:      "Enter Range Totals",
				Link:      urlEnterTotals,
				RangeMenu: true,
			}, {
				Name: "Event Report",
				Link: urlEventReport,
			}, {
				Name:   "Print Entry List",
				Link:   urlEntryList,
				Hidden: true,
			}},
		}, {
			Name: "Clubs",
			Link: urlClubs,
			SubMenu: []menu{{
				Name: "Club",
				Link: urlClub,
			}},
		}, {
			Name: "Shooters",
			Link: urlShooters,
		}, {
			Name: "Archive",
			Link: urlArchive,
		}, {
			Name: "About",
			Link: urlAbout,
		}},
	}
)

var (
	templateFuncMap = template.FuncMap{
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
					output = attribute + "=" + addQuotes(value)
				}
			case int:
				if value.(int) > 0 {
					output = attribute + "=" + addQuotes(value)
				}
			case uint:
				if value.(uint) > 0 {
					output = attribute + "=" + addQuotes(value)
				}
			}
			return template.HTMLAttr(output)
		},
		"has": func(t interface{}, value string) template.HTMLAttr {
			var hasValue bool
			switch t.(type) {
			default:
				warn.Printf("unexpected type %T", t) //%T prints whatever type t has
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
		"grade": findGrade,
		"ageGroup": func(index uint) string {
			if index >= 1 && index < uint(len(dataListAgeGroup())) {
				return dataListAgeGroup()[index].Label
			}
			return ""
		},
		"ordinal": ordinal,
		"findRange": func(id interface{}, ranges []Range) Range {
			var rangeID uint
			switch id.(type) {
			case uint:
				rangeID = id.(uint)
			case string:
				rangeID, _ = stoU(id.(string))
			}
			for _, r := range ranges {
				if rangeID == r.ID {
					return r
				}
			}
			return Range{}
		},
		"N": func(start, end uint) (stream chan uint) {
			stream = make(chan uint)
			go func() {
				var i uint = start
				for ; i <= end; i++ {
					stream <- i
				}
				close(stream)
			}()
			return
		},
	}
)

func templater(w http.ResponseWriter, page page) {
	//Add HTTP headers so browsers don't cache the HTML resource because it can contain different content every request.
	headers(w, nocache)
	if page.template == 25 {
		page.template = 0
		w.Header().Set(csp, "style-src 'self'")
	} else {
		w.Header().Set(csp, "default-src 'none'; style-src 'self'; script-src 'self'; connect-src ws: 'self'; img-src 'self' data:") //font-src 'self'
	}

	//Convert page.Title to the HTML template file name (located within htmlDirectory), e.g. Events > Events, Club Settings > ClubSettings
	pageName := strings.Replace(strings.Title(page.Title), " ", "", -1)

	htmlFileNames := []string{htmlDirectory + pageName}
	if page.template != templateNone {
		htmlFileNames = append(htmlFileNames, masterStuff[page.template]...)
	}

	//Add page content just generated to the default page environment (which has CSS and JS, etc).
	masterTemplate.Page = page

	html, ok := templates[pageName]
	//debug is for dynamically re-parsing (reloading) templates on every request
	if !ok || debug {

		templates[pageName] = template.Must(template.New("main").Funcs(templateFuncMap).ParseFiles(htmlFileNames...))
		html = templates[pageName]
	}

	if err := html.ExecuteTemplate(w, "master", masterTemplate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//AddQuotes returns value with or without surrounding single or double quote characters suitable for a [[//dev.w3.org/html5/html-author/#attributes][HTML5 attribute]] value.
func addQuotes(in interface{}) string {
	value := strings.Replace(fmt.Sprintf("%v", in), `&`, "&amp;", -1) //will destroy any existing escaped characters like &#62;
	double := strings.Count(value, `"`)
	single := strings.Count(value, `'`)
	if single > 0 && single >= double {
		return `"` + strings.Replace(value, `"`, "&#34;", -1) + `"`
	}
	//Space, double quote, accent, equals, less-than sign, greater-than sign.
	if double > 0 || strings.ContainsAny(value, " \"`=<>") {
		return `'` + strings.Replace(value, `'`, "&#39;", -1) + `'`
	}
	return value
}
