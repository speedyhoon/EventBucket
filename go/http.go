package main

import (
	"net/http"
	"strings"
)

func serveFile(fileName string) {
	http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
		headers(w, []string{cache})
		serveGzip(w, r,
			func() {
				http.ServeFile(w, r, "./gz/"+fileName)
			},
			func() {
				http.ServeFile(w, r, "./"+fileName)
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
					http.StripPrefix(contentType, http.FileServer(http.Dir("./gz/"))).ServeHTTP(w, r)
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
	//dirCSS:      {contentType, "text/css; charset=utf-8"},
	//dirSVG:      {contentType, "image/svg+xml"},
	//dirWOF2:     {contentType, "application/font-woff2"},
	//dirPNG:      {contentType, "image/png"},
	//dirJPEG:     {contentType, "image/jpeg"},
	//dirWOF:      {contentType, "application/font-woff"},
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

/*

func get(url string, runner func() Page) {
	//Setup url as a subdirectory path. e.g. if url = "foobar" then "http://localhost/foobar" is available
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) { templater(runner(), w, r) }))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	//	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer devModeTimeTrack(time.Now(), r.RequestURI)
		headers(w, []string{"html", nocache})
		//		gzipper(h, w, r)
	})
}
*/
