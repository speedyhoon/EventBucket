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
			httpHeaders(w, []string{"expire", "cache", contentType})
			NewGzipper(http.FileServer(http.Dir(DIR_ROOT)), w, r)
		})))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc{
//PROD	return func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
		httpHeaders(w, []string{"html", "nocache"})
		NewGzipper(h, w, r)
	}))
}

func NewGzipper(h http.Handler, w http.ResponseWriter, r *http.Request){
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
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) {templatorNew(runner(), w, r)}))
}
func GetRedirectPermanent(url string, runner func()Page){
	Get(url, runner)
	http.Handle(url + "/", http.RedirectHandler(url, http.StatusMovedPermanently))
}
func GetParameters(url string, runner func(string)Page) {
	h := func(w http.ResponseWriter, r *http.Request) {templatorNew(runner(getIdFromUrl(r, url)), w, r)}
	http.Handle(url, serveHtml(h))
}
func Post(url string, runner http.HandlerFunc){
	http.HandleFunc(url, serveHtml(runner))
}

func httpHeaders(w http.ResponseWriter, set_headers []string) {
	headers := map[string][2]string{
		//Expiry date is in 1 year, 0 months & 0 days
		"expire": [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO it should return GMT time I think
		"cache":    [2]string{"Vary", "Accept-Encoding"},
		"gzip":     [2]string{"Content-Encoding", "gzip"},
		"html":     [2]string{"Content-Type", "text/html; charset=utf-8"},
		DIR_CSS:    [2]string{"Content-Type", "text/css; charset=utf-8"},
		DIR_JS:     [2]string{"Content-Type", "text/javascript"},
		DIR_PNG:    [2]string{"Content-Type", "image/png"},
		DIR_JPEG:   [2]string{"Content-Type", "image/jpeg"},
//		DIR_GIF:    [2]string{"Content-Type", "image/gif"},
//		DIR_WEBP:   [2]string{"Content-Type", "image/webp"},
		DIR_SVG:    [2]string{"Content-Type", "image/svg+xml"},
	}
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';")
	for _, lookup := range set_headers {
		if lookup != "nocache" {
			w.Header().Set(headers[lookup][0], headers[lookup][1])
		}else{
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Expires", "0")
			w.Header().Set("Pragma", "no-cache")
		}
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

/*func redirectPermanent(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently) //Search engine Optimisation
	}
}

func redirectVia(runThisFirst func(http.ResponseWriter, *http.Request), path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, path, http.StatusSeeOther) //303 mandating the change of request type to GET
	}
}*/

func redirecter(path string, w http.ResponseWriter, r *http.Request) {
	//	return func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, path, http.StatusSeeOther)
	//	}
}

func getIdFromUrl(r *http.Request, page_url string) string {
	//TODO add validation checking for id using regex pattens
	//TODO add a http layer function between p_page functions and main.go so that the event_id or club_id can be validated and the p_page functions don't have to interact with http at all
	return r.URL.Path[len(page_url):]
}
