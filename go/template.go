package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/speedyhoon/text/template"
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
	template                     string
	SubTemplate                  string //Which template to load within the main template
	skipCSP                      bool
	Status                       int
}

func (p page) csp() string {
	if p.skipCSP {
		return open
	}
	return lock
}

type markupEnv struct {
	Page        page
	Menu        []menu
	IsDarkTheme bool
}

var (
	templates      *template.Template
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
				Name:      "Enter Totals",
				Link:      urlEnterTotals,
				RangeMenu: true,
			}, {
				Name: "Event Report",
				Link: urlEventReport,
			}, {
				Name: "Shooters Report",
				Link: urlShootersReport,
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
		}, {
			Name: "Settings",
			Link: urlSettings,
		}, {
			Name: "Licence",
			Link: urlLicence,
		}},
	}
)

//TODO remove if !debug
func init() {
	if err := loader(); !debug && err != nil {
		warn.Fatal(err)
	}
}

func loader() (err error) {
	templates, err = template.New("").Funcs(template.FuncMap{
		"a": func(attribute string, value interface{}) string {
			//"a": func(attribute string, value interface{}) template.HTMLAttr {
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
			case []option:
				if len(value.([]option)) > 0 {
					output = attribute
				}
			//TODO remove default if !debug
			default:
				warn.Printf("attribute type %T not defined\n%v %v\n", value, value, len(value.([]option)))
			}
			//return template.HTMLAttr(output)
			return output
		},
		"ageGroup": func(index uint) string {
			if index < uint(len(dataListAgeGroup())) {
				return dataListAgeGroup()[index].Label
			}
			return ""
		},
		"ordinal": func(x int) string {
			return ordinal(uint(x), false)
		},
		"findRange": func(rangeID uint, ranges []Range) Range {
			for _, r := range ranges {
				if rangeID == r.ID {
					return r
				}
			}
			return Range{}
		},
	}).ParseFiles(
		filepath.Join(runDir, "h"),
	)
	return
}

func templater(w http.ResponseWriter, page page) {
	//Gzip response, even if requester doesn't support gzip
	gz := gzip.NewWriter(w)
	defer gz.Close()
	wz := gzipResponseWriter{Writer: gz, ResponseWriter: w}

	//Add HTTP headers so browsers don't cache the HTML resource because it may contain different content every request.
	headers(wz, html, nocache, cGzip, page.csp())

	if page.Status != 0 {
		wz.WriteHeader(page.Status)
	}

	//Convert page.Title to the lowercase HTML template file name
	page.SubTemplate = strings.Replace(strings.ToLower(page.Title), " ", "", -1)

	//Add page content just generated to the default page environment (which has CSS and JS, etc).
	masterTemplate.Page = page

	//TODO optionally remove during build time if debug == true
	if debug {
		if err := loader(); err != nil {
			warn.Println(err)
			http.Error(wz, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := templates.ExecuteTemplate(wz, masterTemplate.Page.template, masterTemplate); err != nil {
		warn.Println(err)
		http.Error(wz, err.Error(), http.StatusInternalServerError)
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

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
