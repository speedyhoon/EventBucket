package main

import (
	"net/http"
	"strings"
)

func serveDir(contentType string) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing.
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				warning.Print("The current browser does not support gzip")
			}
			w.Header().Set("Content-Encoding", "gzip")

			//			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")	//Don't cache any response
			//			w.Header().Set("Pragma", "no-cache")
			//			w.Header().Set("Expires", "0")	//Expires now
			//			w.Header().Set("Content-Encoding", "gzip")
			//			Gzip(http.FileServer(http.Dir(DIR_ROOT)), w, r)
			//			http.FileServer(http.Dir("./gz/")).ServeHTTP(w, r)
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				http.StripPrefix(contentType, http.FileServer(http.Dir("./gz/"))).ServeHTTP(w, r)
			} else {
				http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
			}
		}))
}

func httpHeaders(w http.ResponseWriter, setHeaders []string) {
	//w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")

	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set("X-Frame-Options", "DENY")
	headers := map[string][2]string{
		//Cache
		"public":    {"Cache-Control", "public"},
		"expire":    {"Expires", expiresTime},
		"cache":     {"Vary", "Accept-Encoding"},

		//Don't cache
		"noCache":   {"Cache-Control", "no-cache, no-store, must-revalidate"},
		"expireNow": {"Expires", "0"},
		"pragma":    {"Pragma", "no-cache"},

		"html":      {"Content-Type", "text/html; charset=utf-8"},
		dirJS:       {"Content-Type", "text/javascript"},
		//dirCSS:      {"Content-Type", "text/css; charset=utf-8"},
		//dirSVG:      {"Content-Type", "image/svg+xml"},
		//dirWOF2:     {"Content-Type", "application/font-woff2"},
		//dirPNG:      {"Content-Type", "image/png"},
		//dirJPEG:     {"Content-Type", "image/jpeg"},
		//dirWOF:      {"Content-Type", "application/font-woff"},
	}
	for _, lookup := range setHeaders {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}
