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
	Data                         map[string]interface{}
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
	masterTemplatePath    = htmlDirectory + "master"
	formsTemplatePath     = htmlDirectory + "forms"
	reusablesTemplatePath = htmlDirectory + "reusables"
)

var (
	templates      = make(map[string]*template.Template)
	masterTemplate = markupEnv{
		Menu: []menu{{
			Name: "Home",
			Link: urlHome,
		}, {
			Name: "Events",
			Link: urlEvents,
			SubMenu: []menu{{
				Name: "Entries",
				Link: urlEntries,
			}, {
				Name: "Settings",
				Link: urlEventSettings,
			}, {
				Name: "Scoreboard",
				Link: urlScoreboard,
			}, {
				Name: "Scorecards",
				Link: urlScorecards,
			}, {
				Name: "Total Scores",
				Link: urlTotalScores,
			}, {
				Name: "Report",
				Link: urlEventReport,
			}},
		}, {
			Name: "Clubs",
			Link: urlClubs,
			SubMenu: []menu{{
				Name: "Club",
				Link: urlClub,
			}, {
				Name: "Settings",
				Link: urlClubSettings,
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
		}, {
			Name: "Licence",
			Link: urlLicence,
		}},
	}
)

func templater(w http.ResponseWriter, page page) {
	//Add HTTP headers so browsers don't cache the HTML resource because it can contain different content every request.
	headers(w, nocache)

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
			"grade": func(index uint64) string {
				if index < uint64(len(dataListGrades())) {
					return dataListGrades()[index].Label
				}
				return ""
			},
			"ageGroup": func(index uint64) string {
				if index < uint64(len(dataListAgeGroup())) {
					return dataListAgeGroup()[index].Label
				}
				return ""
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

//AddQuotes returns value with or without surrounding single or double quote characters suitable for a [[//dev.w3.org/html5/html-author/#attributes][HTML5 attribute]] value.
func addQuotes(value string) string {
	//TODO escape any rune character over X code point
	//	value = html.EscapeString(value)
	//	escaper:= = strings.NewReplacer(
	//		`&`, "&amp;",
	//`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	//		`<`, "&lt;",
	//		`>`, "&gt;",
	//`"`, "&#34;", // "&#34;" is shorter than "&quot;".
	//	)
	value = strings.Replace(value, `&`, "&amp;", -1) //will destroy any existing escaped characters like &#62;
	double := strings.Count(value, `"`)
	single := strings.Count(value, `'`)
	if single > 0 && single >= double {
		return `"` + strings.Replace(value, `"`, "&#34;", -1) + `"`
	}
	if double > 0 || strings.ContainsAny(value, " `=<>") {
		return `'` + strings.Replace(value, `'`, "&#39;", -1) + `'`
	}
	/*//Contains a single quote and a double quote character.
	if strings.Contains(value, "'") && strings.Contains(value, `"`) {
		warn.Printf("HTML attribute value %v contains both single & double quotes", value)
	}
	//Space, single quote, accent, equals, less-than sign, greater-than sign.
	if strings.ContainsAny(value, " '`=<>") {
		return `"` + value + `"`
	}
	//Double quote
	if strings.Contains(value, `"`) {
		return "'" + value + "'"
	}*/
	return value
}
