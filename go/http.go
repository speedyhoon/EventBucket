package main

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	dirRoot     = "./"
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

	contentType    = "Content-Type"
	cacheControl   = "Cache-Control"
	expires        = "Expires"
	cache          = "cache"
	nocache        = "nocache"
	gzip           = "gzip"
	acceptEncoding = "Accept-Encoding"
)

func serveFile(fileName string) {
	http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
		// Check if the request contains accept gzip encoding header & return the appropriate resource
		// Unfortunately uncompressed responses may still be required even though all modern browsers support gzip
		//webmasters.stackexchange.com/questions/22217/which-browsers-handle-content-encoding-gzip-and-which-of-them-has-any-special
		//www.stevesouders.com/blog/2009/11/11/whos-not-getting-gzip/
		//BUG gzip serving isn't working
		/*if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
			headers(w, []string{cache, gzip})
			warn.Println("Gzipper", dirGzip+fileName)
			http.ServeFile(w, r, dirGzip+fileName)
		} else {*/
		headers(w, []string{cache})
		//		warn.Println("no Gzip", dirRoot+fileName)
		http.ServeFile(w, r, dirRoot+fileName)
		//		warn.Print("The request didn't contain gzip")
		//		}
	})
}

func serveDir(contentType string, allowGzip bool) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing.
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			if allowGzip && strings.Contains(r.Header.Get(acceptEncoding), gzip) {
				headers(w, []string{contentType, gzip, cache})
				http.StripPrefix(contentType, http.FileServer(http.Dir(dirGzip))).ServeHTTP(w, r)
			} else {
				headers(w, []string{contentType, cache})
				http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
				//				warn.Print("The request didn't contain gzip")
			}
		}))
}

var headerOptions = map[string][2]string{
	gzip:   {"Content-Encoding", "gzip"},
	"html": {contentType, "text/html; charset=utf-8"},
	dirCSS: {contentType, "text/css; charset=utf-8"},
	dirJS:  {contentType, "text/javascript"},
	dirPNG: {contentType, "image/png"},
	//dirSVG:    {contentType, "image/svg+xml"},
	//dirWOF2:   {contentType, "application/font-woff2"},
	//dirJPEG:   {contentType, "image/jpeg"},
}

//research //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders []string) {
	//w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'")

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

func getParameters(url string, pageFunc func(http.ResponseWriter, *http.Request, string), regex *regexp.Regexp) {
	var parameters, lowerParams string
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			parameters = strings.TrimPrefix(r.URL.Path, url)
			lowerParams = strings.ToLower(parameters)
			if regex.MatchString(lowerParams) {
				//normal event id specified
				pageFunc(w, r, lowerParams)
				return
			}

			if parameters == "" || !strings.HasPrefix(r.URL.Path, url) {
				//no prefix - redirect to + 's':", strings.TrimSuffix(url, "/")+"s")
				http.Redirect(w, r, strings.TrimSuffix(url, "/")+"s", http.StatusNotFound)
				return
			}
			if parameters != lowerParams {
				//redirect to lowercase event page
				http.Redirect(w, r, url+lowerParams, http.StatusSeeOther)
				return
			}

			//parameters don't match the regex string - 404 enent id not found
			errorHandler(w, r, http.StatusNotFound)
		})
}
