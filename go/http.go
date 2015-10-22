package main

import (
	"net/http"
	"strings"
)

func serveFile(fileName string) {
	http.HandleFunc("/" + fileName, func(w http.ResponseWriter, r *http.Request) {

			httpHeaders(w, []string{cache})
			gzipContents(w, r,
				func(){
					http.ServeFile(w, r, "./gz/"+fileName)
				},
				func(){
					http.ServeFile(w, r, "./"+fileName)
				})

//			if strings.Contains(r.Header.Get(acceptEncoding), "gzip") {
//				http.ServeFile(w, r, "./gz/"+fileName)
//			}else {
//				http.ServeFile(w, r, "./"+fileName)
//				trace.Print("served " + fileName + " with NO compression")
//			}
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

			//			w.Header().Set(cacheControl, "no-cache, no-store, must-revalidate")	//Don't cache any response
			//			w.Header().Set("Pragma", "no-cache")
			//			w.Header().Set("Expires", "0")	//Expires now
			//			w.Header().Set("Content-Encoding", "gzip")
			//			Gzip(http.FileServer(http.Dir(DIR_ROOT)), w, r)
			//			http.FileServer(http.Dir("./gz/")).ServeHTTP(w, r)
			httpHeaders(w, []string{cache})

			gzipContents(w, r,
				func(){
					http.StripPrefix(contentType, http.FileServer(http.Dir("./gz/"))).ServeHTTP(w, r)
				},
				func(){
					http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
				})


//			if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
//				httpHeaders(w, []string{gzip})
//				http.StripPrefix(contentType, http.FileServer(http.Dir("./gz/"))).ServeHTTP(w, r)
//			} else {
//				http.FileServer(http.Dir(dirRoot)).ServeHTTP(w, r)
//				warning.Print("The current browser does not support gzip")
//			}
		}))
}

func gzipContents(w http.ResponseWriter, r *http.Request, serveNormal, serveGzip func()){
	if strings.Contains(r.Header.Get(acceptEncoding), gzip) {
		httpHeaders(w, []string{gzip})
		serveGzip()
	} else {
		serveNormal()
		warning.Print("The current browser does not support gzip")
	}
}



const (
	contentType = "Content-Type"
	cacheControl = "Cache-Control"
	expires = "Expires"
	cache   = "cache"
	nocache = "nocache"
	gzip = "gzip"
	acceptEncoding = "Accept-Encoding"
)

func httpHeaders(w http.ResponseWriter, setHeaders []string) {
	//w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; font-src 'self'")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")

	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set("X-Frame-Options", "DENY")
	headers := map[string][2]string{
		//Cache
//		"public":    {cacheControl, "public"},
//		"expire":    {expires, expiresTime},
//		"cache":     {"Vary", acceptEncoding},

		//Don't cache
//		"noCache":   {cacheControl, "no-cache, no-store, must-revalidate"},
//		"expireNow": {expires, "0"},
//		"pragma":    {"Pragma", "no-cache"},

		gzip:      {"Content-Encoding", "gzip"},
		"html":      {contentType, "text/html; charset=utf-8"},
		dirJS:       {contentType, "text/javascript"},
		//dirCSS:      {contentType, "text/css; charset=utf-8"},
		//dirSVG:      {contentType, "image/svg+xml"},
		//dirWOF2:     {contentType, "application/font-woff2"},
		//dirPNG:      {contentType, "image/png"},
		//dirJPEG:     {contentType, "image/jpeg"},
		//dirWOF:      {contentType, "application/font-woff"},
	}
	for _, lookup := range setHeaders {
		switch lookup{
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
			w.Header().Set(headers[lookup][0], headers[lookup][1])
		}
	}
}
