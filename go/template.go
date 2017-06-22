package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
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
	JS                           []string
	Section                      string //Which template to load within the main template
	skipCSP                      bool
}

func (p page) csp() string {
	if p.skipCSP {
		return "open"
	}
	return "lock"
}

type markupEnv struct {
	Page  page
	Menu  []menu
	Theme bool
}

var (
	NewTemplates   *template.Template
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
			Name: "Licence",
			Link: urlLicence,
		}},
	}
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//TODO remove if !debug
func init() {
	if err := loader(); err != nil {
		warn.Println(err)
		os.Exit(5)
	}
}

func loader() (err error) {
	NewTemplates, err = template.New("").Funcs(template.FuncMap{
		"i": func(inputs []field, index int) *field {
			//index will always be a positive integer so the check for index >= 0 is not required
			if index < len(inputs) {
				return &inputs[index]
			}
			return nil
		},
		"a": func(attribute string, value interface{}) string { return "" },
		/*"a": func(attribute string, value interface{}) template.HTMLAttr {
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
		},*/
		//TODO remove has
		"has": func(attribute string, value interface{}) string { return "" },
		/*"has": func(t interface{}, value string) template.HTMLAttr {
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
		},*/
		"grade": findGrade,
		"ageGroup": func(index uint) string {
			if index >= 1 && index < uint(len(dataListAgeGroup())) {
				return dataListAgeGroup()[index].Label
			}
			return ""
		},
		"ordinal": func(x int) string {
			return ordinal(uint(x), false)
		},
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
	headers(wz, "html", nocache, cGzip, page.csp())

	//Convert page.Title to the lowercase HTML template file name
	page.Section = strings.Replace(strings.ToLower(page.Title), " ", "", -1)

	//Add page content just generated to the default page environment (which has CSS and JS, etc).
	masterTemplate.Page = page

	//TODO optionally remove during build time if debug == true
	if err := loader(); err != nil {
		warn.Println(err)
		http.Error(wz, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := NewTemplates.ExecuteTemplate(wz, masterTemplate.Page.template, masterTemplate); err != nil {
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
