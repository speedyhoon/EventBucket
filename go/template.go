package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/speedyhoon/cnst"
	"github.com/speedyhoon/cnst/mime"
	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/text/template"
	"github.com/speedyhoon/utl"
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
	SubTemplate                  string // Which template to load within the main template.
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
	mainTheme bool
	tmpl8     *template.Template
	tempFuncs = template.FuncMap{
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
					output = attribute + "=" + utl.AddQuotes(value)
				}
			case int:
				if value.(int) > 0 {
					output = attribute + "=" + utl.AddQuotes(value)
				}
			case uint:
				if value.(uint) > 0 {
					output = attribute + "=" + utl.AddQuotes(value)
				}
			case []frm.Option:
				if len(value.([]frm.Option)) > 0 {
					output = attribute
				}
			//#ifdef DEBUG
			default:
				log.Printf("attribute type %T not defined\n%v %v\n", value, value, len(value.([]frm.Option)))
				//#endif
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
			return utl.Ordinal(uint(x), false)
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

// #ifndef DEBUG
func init() {
	var err error
	tmpl8, err = template.New("").Funcs(tempFuncs).ParseFiles(
		filepath.Join(runDir, "h"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

//#endif

func render(w http.ResponseWriter, p page) {
	// Brotli compress response, even if AcceptEncoding doesn't contain "br".
	wz := compressResponse{
		WriteCloser:    brotli.NewWriterLevel(w, brotli.BestCompression),
		ResponseWriter: w,
	}
	defer func() {
		if err := wz.WriteCloser.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Add HTTP headers so browsers don't cache the HTML resource because it may contain different content every request.
	headers(wz, mime.HTML, nocache, cnst.Gzip, p.csp())

	if p.Status != 0 {
		wz.WriteHeader(p.Status)
	}

	if p.SubTemplate == "" {
		// Convert page.Title to the lowercase HTML template file name.
		p.SubTemplate = strings.Replace(strings.ToLower(p.Title), " ", "", -1)
	}

	//#ifdef DEBUG
	var err error
	tmpl8, err = template.New("").Funcs(tempFuncs).ParseFiles(
		filepath.Join(runDir, "h"),
	)
	if err != nil {
		log.Println(err)
		http.Error(wz, err.Error(), http.StatusInternalServerError)
		return
	}
	//#endif

	type markupEnv struct {
		Page  page
		Menu  []menu
		Theme uint8
	}

	me := markupEnv{Page: p, Menu: mainMenu}
	if mainTheme {
		me.Theme = 1
	}

	if err = tmpl8.ExecuteTemplate(wz, p.template, me); err != nil {
		log.Println(err)
		http.Error(wz, err.Error(), http.StatusInternalServerError)
	}
}

type compressResponse struct {
	io.WriteCloser
	http.ResponseWriter
}

func (w compressResponse) Write(b []byte) (int, error) {
	return w.WriteCloser.Write(b)
}
