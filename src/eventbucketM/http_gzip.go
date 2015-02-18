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
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
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
	//TODO Only set CSP when not in debug mode
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:;")
	headers := map[string][2]string{
		"expire":  [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO should it return GMT time?  //Expiry date is in 1 year, 0 months & 0 days in the future
		"cache":   [2]string{"Vary", "Accept-Encoding"},
		"public":  [2]string{"Cache-Control", "public"},
		"gzip":    [2]string{"Content-Encoding", "gzip"},
		"html":    [2]string{"Content-Type", "text/html; charset=utf-8"},
		"noCache": [2]string{"Cache-Control", "no-cache, no-store, must-revalidate"},
		"expireNow":[2]string{"Expires", "0"},
		"pragma":  [2]string{"Pragma", "no-cache"},
		DIR_CSS:   [2]string{"Content-Type", "text/css; charset=utf-8"},
		DIR_JS:    [2]string{"Content-Type", "text/javascript"},
		DIR_PNG:   [2]string{"Content-Type", "image/png"},
		DIR_JPEG:  [2]string{"Content-Type", "image/jpeg"},
//		DIR_GIF:   [2]string{"Content-Type", "image/gif"},
//		DIR_WEBP:  [2]string{"Content-Type", "image/webp"},
		DIR_SVG:   [2]string{"Content-Type", "image/svg+xml"},
	}
	for _, lookup := range set_headers {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}

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
