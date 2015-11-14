package main

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	dirGzip     = "dirGzip"
	urlHome     = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlEvents   = "/events"
	urlLicence  = "/licence"
	urlShooters = "/shooters"
	//GET with PARAMETERS
	urlEvent = "/event/" //eventID
)

func serveFile(fileName string) {
	http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
		headers(w, []string{cache})
		serveGzip(w, r,
			func() {
				http.ServeFile(w, r, dirGzip+fileName)
			},
			func() {
				http.ServeFile(w, r, dirRoot+fileName)
			})
	})
}

func serveDir(contentType string) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing.
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			headers(w, []string{cache})
			serveGzip(w, r,
				func() {
					http.StripPrefix(contentType, http.FileServer(http.Dir(dirGzip))).ServeHTTP(w, r)
				},
				func() {
					http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
				})
		}))
}

//Check if the request contains accept gzip encoding header & execute the appropriate function
func serveGzip(w http.ResponseWriter, r *http.Request, ungzipped, gzipped func()) {
	if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
		headers(w, []string{gzip})
		gzipped()
	} else {
		ungzipped()
		warn.Print("The current browser does not support gzip")
	}
}

const (
	contentType    = "Content-Type"
	cacheControl   = "Cache-Control"
	expires        = "Expires"
	cache          = "cache"
	nocache        = "nocache"
	gzip           = "gzip"
	acceptEncoding = "Accept-Encoding"
)

var headerOptions = map[string][2]string{
	gzip:   {"Content-Encoding", "gzip"},
	"html": {contentType, "text/html; charset=utf-8"},
	dirJS:  {contentType, "text/javascript"},
	//dirCSS:    {contentType, "text/css; charset=utf-8"},
	//dirSVG:    {contentType, "image/svg+xml"},
	//dirWOF2:   {contentType, "application/font-woff2"},
	//dirPNG:    {contentType, "image/png"},
	//dirJPEG:   {contentType, "image/jpeg"},
	//dirWOF:    {contentType, "application/font-woff"},
}

//research //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders []string) {
	//w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")

	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set("X-Frame-Options", "DENY")
	for _, lookup := range setHeaders {
		switch lookup {
		case cache:
			w.Header().Set(cacheControl, "public")
			w.Header().Set(expires, expiresTime)
			w.Header().Set("Vary", acceptEncoding)
			break
		case nocache:
			w.Header().Set(cacheControl, "no-cache, no-store, must-revalidate")
			w.Header().Set(expires, "0")
			w.Header().Set("Pragma", "no-cache")
			break
		default:
			w.Header().Set(headerOptions[lookup][0], headerOptions[lookup][1])
		}
	}
}

func get404(url string, pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url {
				errorHandler(w, r, http.StatusNotFound)
				return
			}
			pageFunc(w, r)
		})
}

func getRedirectPermanent(url string, pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url, pageFunc)
	//Redirects back to subdirectory "url". Needed when url parameters are not wanted or needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url+"/", http.RedirectHandler(url, http.StatusMovedPermanently))
}

func getParameters(url string, pageFunc func(http.ResponseWriter, *http.Request, string), regex, regexWeak *regexp.Regexp) {
	var parameters string
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			parameters = strings.TrimPrefix(r.URL.Path, url)
			info.Println("getParameters", parameters)
			if !regex.MatchString(parameters) {
				info.Println("failed regex", `"`, r.URL.Path, `"`, "url=", url)
				if regexWeak.MatchString(parameters) {
					info.Println("passed Weak")
					http.Redirect(w, r, strings.ToLower(r.URL.Path), http.StatusSeeOther)
				} else {
					info.Println("failed ALL")
					errorHandler(w, r, http.StatusNotFound)
				}
				return
			}
			pageFunc(w, r, parameters)
		})
}
