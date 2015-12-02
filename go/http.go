package main

import (
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
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
		// Check if the request contains accept gzip encoding header & return the appropriate resource
		// Unfortunately uncompressed responses may still be required even though all modern browsers support gzip
		//webmasters.stackexchange.com/questions/22217/which-browsers-handle-content-encoding-gzip-and-which-of-them-has-any-special
		//www.stevesouders.com/blog/2009/11/11/whos-not-getting-gzip/
		if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
			headers(w, []string{cache, gzip})
			warn.Println("Gzipper", dirGzip+fileName)
			http.ServeFile(w, r, dirGzip+fileName)
		} else {
			headers(w, []string{cache})
			warn.Println("no Gzip", dirRoot+fileName)
			http.ServeFile(w, r, dirRoot+fileName)
			warn.Print("The request didn't contain gzip")
		}
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
			if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
				headers(w, []string{gzip, cache})
				http.StripPrefix(contentType, http.FileServer(http.Dir(dirGzip))).ServeHTTP(w, r)
			} else {
				headers(w, []string{cache})
				http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
				warn.Print("The request didn't contain gzip")
			}
		}))
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
	//dirJS:  {contentType, "text/javascript"},
	//dirCSS:    {contentType, "text/css; charset=utf-8"},
	//dirSVG:    {contentType, "image/svg+xml"},
	//dirWOF2:   {contentType, "application/font-woff2"},
	dirPNG: {contentType, "image/png"},
	//dirJPEG:   {contentType, "image/jpeg"},
	//dirWOF:    {contentType, "application/font-woff"},
}

//research //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders []string) {
	//w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
	//	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'")

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

//SessionID's can't have space or semicolon
//should be 16 characters
//stackoverflow.com/questions/1969232/allowed-characters-in-cookies
func newSessionID() string {
	var newSessionId string
	var i, randInt int
	for i < 24 {
		randInt = 33 + rand.Intn(93)
		if randInt != 59 { //ignore semicolons ;
			i++
			newSessionId += string(randInt)
		}
	}
	return "z=" + newSessionId
}

func sessionError(error string) string {
	sessionID := newSessionID()
	globalSessions[sessionID] = error
	return sessionID + "; Expires=" + time.Now().Add(1*time.Minute).Format("Mon, Jan 2 2006 15:04:05 GMT")
}

var globalSessions map[string]string
