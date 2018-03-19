package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/speedyhoon/text/template"
	"github.com/speedyhoon/forms"
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

var (
	mainTheme uint8
	tmpl8     *template.Template
	tempFuncs   = template.FuncMap{
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
			case []forms.Option:
				if len(value.([]forms.Option)) > 0 {
					output = attribute
				}
			//TODO remove default if !debug
			default:
				warn.Printf("attribute type %T not defined\n%v %v\n", value, value, len(value.([]forms.Option)))
			}
			//return template.HTMLAttr(output)
			return output
		},
		"grade": findGrade,
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
		"N": func(start, end uint) (stream chan uint) {
			stream = make(chan uint)
			go func() {
				for i := start; i <= end; i++ {
					stream <- i
				}
				close(stream)
			}()
			return
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}
)

func init() {
	var err error
	tmpl8, err = template.New("").Funcs(tempFuncs).ParseFiles(
		filepath.Join(runDir, "h"),
	)
	if err != nil {
		warn.Fatal(err)
	}
}

func render(w http.ResponseWriter, p page) {
	//Gzip response, even if requester doesn't support gzip
	gz := gzip.NewWriter(w)
	defer func() {
		if err := gz.Close(); err != nil {
			warn.Println(err)
		}
	}()
	wz := gzipResponseWriter{Writer: gz, ResponseWriter: w}

	//Add HTTP headers so browsers don't cache the HTML resource because it may contain different content every request.
	headers(wz, html, nocache, cGzip, p.csp())

	if p.Status != 0 {
		wz.WriteHeader(p.Status)
	}

	if p.SubTemplate == "" {
		//Convert page.Title to the lowercase HTML template file name
		p.SubTemplate = strings.Replace(strings.ToLower(p.Title), " ", "", -1)
	}

	//TODO optionally remove during build time if debug == false
	if debug {
		var err error
		tmpl8, err = template.New("").Funcs(tempFuncs).ParseFiles(
			filepath.Join(runDir, "h"),
		)
		if err != nil {
			warn.Println(err)
			http.Error(wz, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	type markupEnv struct {
		Page  page
		Menu  []menu
		Theme uint8
	}

	if err := tmpl8.ExecuteTemplate(wz, p.template, markupEnv{Page: p, Menu: mainMenu, Theme: mainTheme}); err != nil {
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
