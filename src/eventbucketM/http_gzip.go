package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
)

func serveDir(contentType string){
	http.Handle(contentType,
//PROD		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			//If url is a directory return a 404 to prevent displaying a directory listing
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			Gzip(http.FileServer(http.Dir(DIR_ROOT)), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc{
//PROD	return func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
		httpHeaders(w, []string{"html", "nocache1", "nocache2", "nocache3"})
		Gzip(h, w, r)
	}))
}

func Gzip(h http.Handler, w http.ResponseWriter, r *http.Request){
	//Return a gzip compressed response if appropriate
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		h.ServeHTTP(gzr, r)
		return
	}
	//else return a normal response
	h.ServeHTTP(w, r)
}

func Get(url string, runner func()Page){
	//Setup url as a subdirectory path. e.g. if url = "foobar" then "http://localhost/foobar" is available
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) {templator(runner(), w, r)}))
}
func GetRedirectPermanent(url string, runner func()Page){
	Get(url, runner)
	//Setup redirect back to subdirectory "url". Needed when url parameters are not wanted/needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url + "/", http.RedirectHandler(url, http.StatusMovedPermanently))
}
func GetParameters(url string, runner func(string)Page) {
	//TODO add a way to get the event id, club id and range id from the url
	h := func(w http.ResponseWriter, r *http.Request) {templator(runner(getIdFromUrl(r, url)), w, r)}
	http.Handle(url, serveHtml(h))
}
//TODO Post - Post to endpoint. If valid return to X page, else return to previous page with form filled out and wrong values highlighted with error message
//TODO Add ajax detection to so ajax requests are not redirected back to the referrer page.
func Post(url string, runner http.HandlerFunc){
	http.HandleFunc(url, serveHtml(runner))
}
func PostVia(runThisFirst func(http.ResponseWriter, *http.Request), url string) func(http.ResponseWriter, *http.Request) {
	//Always redirect after a successful Post to "url". Otherwise redirect back to referrer page.
	//When Ajax is not in use this stops the server responding to Post requests and causes the user to request page "url"
	//This stops the browser from displaying the Post url in the address bar
	//TODO redirect to referrer page on Post failure
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, url, http.StatusSeeOther) //303 mandating the change of request type to GET
	}
}

func httpHeaders(w http.ResponseWriter, set_headers []string) {
	headers := map[string][2]string{
		//Expiry date is in 1 year, 0 months & 0 days
		"expire":  [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO it should return GMT time I think?
		"cache":   [2]string{"Vary", "Accept-Encoding"},
		"public":  [2]string{"Cache-Control", "public"},
		"gzip":    [2]string{"Content-Encoding", "gzip"},
		"html":    [2]string{"Content-Type", "text/html; charset=utf-8"},
		"nocache1":[2]string{"Cache-Control", "no-cache, no-store, must-revalidate"},
		"nocache2":[2]string{"Expires", "0"},
		"nocache3":[2]string{"Pragma", "no-cache"},
		DIR_CSS:   [2]string{"Content-Type", "text/css; charset=utf-8"},
		DIR_JS:    [2]string{"Content-Type", "text/javascript"},
		DIR_PNG:   [2]string{"Content-Type", "image/png"},
		DIR_JPEG:  [2]string{"Content-Type", "image/jpeg"},
//		DIR_GIF:   [2]string{"Content-Type", "image/gif"},
//		DIR_WEBP:  [2]string{"Content-Type", "image/webp"},
		DIR_SVG:   [2]string{"Content-Type", "image/svg+xml"},
	}
	//TODO re-enable security polocy
//	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:;")
	for _, lookup := range set_headers {
//		if lookup != "nocache" {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
		/*}else{
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Expires", "0")
			w.Header().Set("Pragma", "no-cache")
		}*/
	}
}

/*func httpHeaders(w http.ResponseWriter, headers []string) {
	headerSettings := map[string]map[string]string{
		//Expiry date is in 1 year, 0 months & 0 days
		"expire": map[string]string{"Expires": time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO it should return GMT time I think?
		"cache":  map[string]string{"Vary": "Accept-Encoding"},
		"nocache":map[string]string{"Cache-Control": "no-cache, no-store, must-revalidate", "Expires": "0", "Pragma": "no-cache"},
		"gzip":   map[string]string{"Content-Encoding": "gzip"},
		"html":   map[string]string{"Content-Type": "text/html; charset=utf-8"},
		DIR_CSS:  map[string]string{"Content-Type": "text/css; charset=utf-8", "Cache-Control": "public"},
		DIR_JS:   map[string]string{"Content-Type": "text/javascript", "Cache-Control": "public"},
		DIR_PNG:  map[string]string{"Content-Type": "image/png", "Cache-Control": "public"},
		DIR_JPEG: map[string]string{"Content-Type": "image/jpeg", "Cache-Control": "public"},
//		DIR_GIF:  map[string]string{"Content-Type": "image/gif", "Cache-Control": "public"},
//		DIR_WEBP: map[string]string{"Content-Type": "image/webp", "Cache-Control": "public"},
		DIR_SVG:  map[string]string{"Content-Type": "image/svg+xml", "Cache-Control": "public"},
		"logos":  map[string]string{},
	}
	//TODO re-enable security polocy
//	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:;")
	for _, lookup := range headers {
		for field, value := range headerSettings[lookup] {
			w.Header().Set(field, value)
		}
	}
}*/

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func getIdFromUrl(r *http.Request, page_url string) string {
	//TODO add validation checking for id using regex pattens
	//TODO add a http layer function between p_page functions and main.go so that the event_id or club_id can be validated and the p_page functions don't have to interact with http at all
	return r.URL.Path[len(page_url):]
}
