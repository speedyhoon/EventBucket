package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
//	"fmt"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is a test."))
}

func main() {
//	http.ListenAndServe(":1113", makeGzipHandler(handler))
	chttp.Handle("/", http.FileServer(http.Dir("./httpsvr")))
//	http.HandleFunc("/", HomeHandler) // homepage
	http.HandleFunc("/", makeGzipHandler(HomeHandler)) // homepage
	http.ListenAndServe(":80", nil)
}

var chttp = http.NewServeMux()

func HomeHandler(w http.ResponseWriter, r *http.Request) {
//
//	if (strings.Contains(r.URL.Path, ".cjv")) {
////		w.Header().Set("Content Type", "text/css,application/x-javascript")
//		w.Header().Set("Content-Type", "text/css")
//		w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
//		w.Header().Set("Vary", "Accept-Encoding")
////		chttp.ServeHTTP(w, r)
////		fmt.Fprintf(w, "BBBBBBBBBBBB")
//	}else
//	}else if (strings.Contains(r.URL.Path, ".")) {
////		chttp.ServeHTTP(w, r)
//	} else {
////		fmt.Fprintf(w, "HomeHandler")
//	}

	if (strings.Contains(r.URL.Path, ".js")){
		if strings.Contains(r.Header.Get("Accept"), "text/css"){
			w.Header().Set("Content-Type", "text/css")
		}else{
			w.Header().Set("Content-Type", "application/javascript")
		}
//		w.Header().Set("Content-Type", "*/*")
	}
	if (strings.Contains(r.URL.Path, ".jpg")){
		w.Header().Set("Cache-Control", "public")
	}

	w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
	w.Header().Set("Vary", "Accept-Encoding")
	chttp.ServeHTTP(w, r)
}

//				text/css
//text/css,*/*;q=0.1
//application/javascript
                                      //application/javascript;text/css
