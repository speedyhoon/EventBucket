//package main
//
//import (
//	"compress/gzip"
//	"io"
//	"net/http"
//	"strings"
//)
//
//type gzipResponseWriter struct {
//	io.Writer
//	http.ResponseWriter
//}
//
//func (w gzipResponseWriter) Write(b []byte) (int, error) {
//	return w.Writer.Write(b)
//}
//
//func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
//			fn(w, r)
//			return
//		}
//		w.Header().Set("Content-Encoding", "gzip")
//		gz := gzip.NewWriter(w)
//		defer gz.Close()
//		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
//		fn(gzr, r)
//	}
//}
//
//func handler(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Content-Type", "text/plain")
//	w.Write([]byte("This is a test."))
//}
//
//func main() {
//	http.ListenAndServe(":1113", makeGzipHandler(handler))
//}




package main

import (
"fmt"
"net/http"
"strings"
)

var chttp = http.NewServeMux()

func main() {

	chttp.Handle("/", http.FileServer(http.Dir("./httpsvr")))

	http.HandleFunc("/", HomeHandler) // homepage
	http.ListenAndServe(":1113", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if (strings.Contains(r.URL.Path, ".cjv")) {
//		w.Header().Set("Content Type", "text/css,application/x-javascript")
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
		w.Header().Set("Vary", "Accept-Encoding")
		chttp.ServeHTTP(w, r)
//		fmt.Fprintf(w, "BBBBBBBBBBBB")
	}else if (strings.Contains(r.URL.Path, ".js") || strings.Contains(r.URL.Path, ".png")) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
		w.Header().Set("Vary", "Accept-Encoding")
		chttp.ServeHTTP(w, r)
	}else if (strings.Contains(r.URL.Path, ".")) {
		chttp.ServeHTTP(w, r)
	} else {
		fmt.Fprintf(w, "HomeHandler")
	}
}

//				text/css
//text/css,*/*;q=0.1
//application/javascript
                                      //application/javascript;text/css
