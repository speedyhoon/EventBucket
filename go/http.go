package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	gt             = "GET"
	contentType    = "Content-Type"
	contentEncode  = "Content-Encoding"
	cacheControl   = "Cache-Control"
	expires        = "Expires"
	cache          = "cache"
	nocache        = "nocache"
	cGzip          = "gzip"
	brotli         = "br"
	acceptEncoding = "Accept-Encoding"
	csp            = "Content-Security-Policy"
	open           = "o"
	lock           = "l"
	html           = "h"
)

var (
	runDir string

	//Used for every HTTP request with cache headers set.
	cacheExpires string

	headerOptions = map[string][2]string{
		cGzip:   {contentEncode, cGzip},
		brotli:  {contentEncode, brotli},
		html:    {contentType, "text/html; charset=utf-8"},
		dirCSS:  {contentType, "text/css; charset=utf-8"},
		dirJS:   {contentType, "text/javascript"},
		urlSVG:  {contentType, "image/svg+xml"},
		dirWEBP: {contentType, "image/webp"},
		open:    {csp, "style-src 'self'"},
		lock:    {csp, "default-src 'none'; style-src 'self'; script-src 'self'; connect-src ws: 'self'; img-src 'self' data:"}, //font-src 'self'
	}
)

func init() {
	go maintainExpires()

	var err error
	runDir, err = os.Executable()
	if err == nil {
		runDir = filepath.Dir(runDir)
	}
}

func serveFile(fileName string, compress bool) {
	http.HandleFunc(fileName, func(w http.ResponseWriter, r *http.Request) {
		h := []string{cache, fileName}
		if compress {
			h = append(h, brotli)
		}
		headers(w, h...)
		http.ServeFile(w, r, filepath.Join(runDir, fileName))
	})
}

func serveDir(contentType string, compress bool) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing.
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			headers(w, contentType, cache)
			if compress {
				headers(w, brotli)
			}
			http.FileServer(http.Dir(runDir)).ServeHTTP(w, r)
		}))
}

//TODO security add Access-Control-Allow-Origin //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders ...string) {
	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set("X-Frame-Options", "DENY")
	for _, lookup := range setHeaders {
		switch lookup {
		case cache:
			w.Header().Set(cacheControl, "public")
			w.Header().Set(expires, cacheExpires)
			w.Header().Set("Vary", acceptEncoding)
		case nocache:
			w.Header().Set(cacheControl, "no-cache, no-store, must-revalidate")
			w.Header().Set(expires, "0")
			w.Header().Set("Pragma", "no-cache")
		default:
			//Set resource content type header or set content encoding gzip header
			if h, ok := headerOptions[lookup]; ok {
				w.Header().Set(h[0], h[1])
			}
		}
	}
}

func get404(pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(urlEvents,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != urlEvents {
				errorHandler(w, r, "")
				return
			}
			pageFunc(w, r)
		})
}

func getRedirectPermanent(url string, pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url,
		isGetMethod(func(w http.ResponseWriter, r *http.Request) {
			pageFunc(w, r)
		}))
	//Redirects back to subdirectory "url". Needed when url parameters are not wanted or needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url+"/", http.RedirectHandler(url, http.StatusMovedPermanently))
}

func getParameters(url string, pageFunc interface{}, regex *regexp.Regexp) {
	var parameters, lowerParams string
	http.HandleFunc(url,
		isGetMethod(func(w http.ResponseWriter, r *http.Request) {
			parameters = strings.TrimPrefix(r.URL.Path, url)
			lowerParams = strings.ToLower(parameters)

			if parameters != lowerParams {
				//Redirect to page with lowercase parameters.
				http.Redirect(w, r, url+lowerParams, http.StatusSeeOther)
				return
			}

			if regex.MatchString(lowerParams) {
				findHandler(pageFunc, w, r, lowerParams)
				return
			}
			errorWrapper(w, r, url)
		}))
}

func findHandler(pageFunc interface{}, w http.ResponseWriter, r *http.Request, lowerParams string) {
	ids := strings.Split(lowerParams, "/")
	switch pageFunc.(type) {
	//pPrintScorecards.go
	case func(http.ResponseWriter, *http.Request, string):
		pageFunc.(func(http.ResponseWriter, *http.Request, string))(w, r, lowerParams)

	case func(http.ResponseWriter, *http.Request, string, string):
		pageFunc.(func(http.ResponseWriter, *http.Request, string, string))(w, r, ids[0], ids[1])

	//pEntries.go
	case func(http.ResponseWriter, *http.Request, Event):
		event, err := getEvent(ids[0])
		if err != nil {
			errorHandler(w, r, "event")
			return
		}
		pageFunc.(func(http.ResponseWriter, *http.Request, Event))(w, r, event)

	//pClub.go
	case func(http.ResponseWriter, *http.Request, Club):
		club, err := getClub(ids[0])
		//If club not found in the database return error club not found (404).
		if err != nil {
			errorHandler(w, r, "club")
			return
		}
		pageFunc.(func(http.ResponseWriter, *http.Request, Club))(w, r, club)

	case func(http.ResponseWriter, *http.Request, Event, sID):
		event, err := getEvent(ids[0])
		if err != nil {
			errorHandler(w, r, "event")
			return
		}
		shooterID, err := stoU(ids[1])
		if err != nil || shooterID >= uint(len(event.Shooters)) {
			errorHandler(w, r, "shooter")
			return
		}
		pageFunc.(func(http.ResponseWriter, *http.Request, Event, sID))(w, r, event, sID(shooterID))

	//pEnterShots.go
	//pEnterTotals.go
	//pScoreboard.go
	case func(http.ResponseWriter, *http.Request, Event, rID):
		event, err := getEvent(ids[0])
		if err != nil {
			errorHandler(w, r, "event")
			return
		}
		rangeID, err := stoU(ids[1])
		if err != nil {
			errorHandler(w, r, "range")
			return
		}
		pageFunc.(func(http.ResponseWriter, *http.Request, Event, rID))(w, r, event, rID(rangeID))
	}
}

func errorWrapper(w http.ResponseWriter, r *http.Request, url string) {
	errorType := "event"
	if url == urlClub {
		errorType = "club"
	}
	errorHandler(w, r, errorType)
}

func isGetMethod(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Don't accept post or put requests
		if r.Method != gt {
			//http.Redirect(w, r, url, http.StatusSeeOther)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
		h(w, r)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, errorType string) {
	//All EventBucket page urls and ids are lowercase
	lowerURL := strings.ToLower(r.URL.Path)

	//Redirect if url contains any uppercase letters.
	if r.URL.Path != lowerURL {
		http.Redirect(w, r, lowerURL, http.StatusSeeOther)
		return
	}
	lowerURL = strings.TrimSuffix(r.URL.Path, "/")

	//check if the request matches any of the pages that don't require parameters
	if strings.Count(lowerURL, "/") >= 2 {
		for _, page := range []string{urlAbout, urlArchive, urlClubs, urlLicence, urlShooters} {
			if strings.HasPrefix(lowerURL, page) {
				//redirect to page without parameters
				http.Redirect(w, r, page, http.StatusSeeOther)
				return
			}
		}
	}

	templater(w, page{
		Title:  "Error",
		Status: http.StatusNotFound,
		Data: map[string]interface{}{
			"Type": errorType,
		},
	})
}

func formError(w http.ResponseWriter, r *http.Request, f form, err error) {
	f.Error = err
	setSession(w, f)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainExpires() {
	setExpiresTime()
	for range time.NewTicker(time.Hour * 23).C {
		//Can't directly change global variables in a go routine, so call an external function.
		setExpiresTime()
	}
}

//Set expiry date 1 year, 0 months & 0 days in the future.
func setExpiresTime() {
	//Date format is the same as Go`s time.RFC1123 but uses "GMT" timezone instead of "UTC" time standard.
	cacheExpires = time.Now().UTC().AddDate(1, 0, 0).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	//w3.org: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
}
