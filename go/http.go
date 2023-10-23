package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/speedyhoon/cnst"
	"github.com/speedyhoon/cnst/hdrs"
	"github.com/speedyhoon/cnst/mime"
)

const (
	cache    = "cache"
	nocache  = "nocache"
	open     = "o"
	lock     = "l"
	dateTime = "2006-01-02 15:04"
)

var (
	runDir string

	headerOptions = map[string][2]string{
		cnst.Gzip:   {hdrs.ContentEncoding, cnst.Gzip},
		cnst.Brotli: {hdrs.ContentEncoding, cnst.Brotli},
		mime.HTML:   {hdrs.ContentType, mime.HTML},
		mime.CSS:    {hdrs.ContentType, mime.CSS},
		mime.JS:     {hdrs.ContentType, mime.JS},
		mime.SVG:    {hdrs.ContentType, mime.SVG},
		mime.WEBP:   {hdrs.ContentType, mime.WEBP},
		open:        {hdrs.CSP, "style-src 'self'"},
		lock:        {hdrs.CSP, "default-src 'none'; style-src 'self'; script-src 'self'; connect-src ws: 'self'; img-src 'self' data:"}, //font-src 'self'
	}
)

func init() {
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
			h = append(h, cnst.Brotli)
		}
		headers(w, h...)
		http.ServeFile(w, r, filepath.Join(runDir, fileName))
	})
}

func serveDir(contentType string, compress bool) {
	http.HandleFunc(contentType, isDir(func(w http.ResponseWriter, r *http.Request) {
		headers(w, contentType, cache)
		if compress {
			headers(w, cnst.Brotli)
		}
		http.FileServer(http.Dir(runDir)).ServeHTTP(w, r)
	}))
}

func isDir(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//If url is a directory
		if strings.HasSuffix(r.URL.Path, "/") {
			//Then return a 404 to prevent displaying all files in the directory
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

// TODO security add Access-Control-Allow-Origin //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders ...string) {
	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set(hdrs.XFrameOptions, cnst.Deny)
	for _, lookup := range setHeaders {
		switch lookup {
		case cache:
			w.Header().Set(hdrs.CacheControl, "public, max-age=31622400")
			w.Header().Set("Vary", hdrs.AcceptEncoding)
		case nocache:
			w.Header().Set(hdrs.CacheControl, "no-cache, no-store, must-revalidate")
			w.Header().Set(hdrs.Expires, "0")
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
		if r.Method != http.MethodGet {
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

	render(w, page{
		Title:  "Error",
		Status: http.StatusNotFound,
		Data: map[string]interface{}{
			"Type": errorType,
		},
	})
}

/*func serveImg(dir, mimeType string, fileSystem http.FileSystem) {
	http.HandleFunc(dir, isDir(func(w http.ResponseWriter, r *http.Request) {
		//If client accepts Webp images
		if strings.Contains(r.Header.Get(cnst.Accept), cnst.WEBP) {
			r.URL.Path += ".webp"
			w.Header().Set(cnst.ContentType, cnst.WEBP)
		} else {
			w.Header().Set(cnst.ContentType, mimeType)
		}
		http.FileServer(fileSystem).ServeHTTP(w, r)
	}))
}

func serveImages(dir, mimeType string, fileSystem http.FileSystem) {
	http.HandleFunc(
		dir,
		isDir(
			func(w http.ResponseWriter, r *http.Request) {
				//If client accepts Webp images
				if strings.Contains(r.Header.Get(cnst.Accept), cnst.WEBP) {
					r.URL.Path += ".webp"
					w.Header().Set(cnst.ContentType, cnst.WEBP)
				} else {
					w.Header().Set(cnst.ContentType, mimeType)
				}
				http.FileServer(fileSystem).ServeHTTP(w, r)
	}))
}*/
