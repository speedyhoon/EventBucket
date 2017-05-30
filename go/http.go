package main

import (
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	contentType    = "Content-Type"
	cacheControl   = "Cache-Control"
	expires        = "Expires"
	cache          = "cache"
	nocache        = "nocache"
	cGzip          = "gzip"
	br             = "br"
	acceptEncoding = "Accept-Encoding"
	csp            = "Content-Security-Policy"
	formatGMT      = "Mon, 02 Jan 2006 15:04:05 GMT" //Date format
)

func serveFile(fileName string) {
	http.HandleFunc(fileName, func(w http.ResponseWriter, r *http.Request) {
		//Check if the request contains accept gzip encoding header & return the appropriate resource
		//Unfortunately uncompressed responses may still be required even though all modern browsers support gzip
		//webmasters.stackexchange.com/questions/22217/which-browsers-handle-content-encoding-gzip-and-which-of-them-has-any-special
		//www.stevesouders.com/blog/2009/11/11/whos-not-getting-gzip/
		headers(w, cache)
		http.ServeFile(w, r, filepath.Join(runDir, fileName))
	})
}

func serveDir(contentType string, compress bool) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing.
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			headers(w, contentType, cache)
			if compress {
				headers(w, br)
			}
			http.FileServer(http.Dir(runDir)).ServeHTTP(w, r)
		}))
}

var (
	headerOptions = map[string][2]string{
		cGzip:   {"Content-Encoding", "gzip"},
		br:      {"Content-Encoding", "br"},
		"html":  {contentType, "text/html; charset=utf-8"},
		dirCSS:  {contentType, "text/css; charset=utf-8"},
		dirJS:   {contentType, "text/javascript"},
		dirSVG:  {contentType, "image/svg+xml"},
		dirWEBP: {contentType, "image/webp"},
	}
	//Used for every HTTP request with cache headers set.
	cacheExpires = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
)

//security add Access-Control-Allow-Origin //net.tutsplus.com/tutorials/client-side-security-best-practices/
func headers(w http.ResponseWriter, setHeaders ...string) {
	//The page cannot be displayed in a frame, regardless of the site attempting to do so. //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	w.Header().Set("X-Frame-Options", "DENY")
	for _, lookup := range setHeaders {
		switch lookup {
		case cache:
			w.Header().Set(cacheControl, "public")
			w.Header().Set(expires, cacheExpires)
			w.Header().Set("Vary", acceptEncoding)
			break
		case nocache:
			w.Header().Set(cacheControl, "no-cache, no-store, must-revalidate")
			w.Header().Set(expires, "0")
			w.Header().Set("Pragma", "no-cache")
			break
		default:
			//Set resource content type header or set content encoding gzip header
			if lookup == cGzip || lookup == br || headerOptions[lookup][0] == "Content-Type" {
				w.Header().Set(headerOptions[lookup][0], headerOptions[lookup][1])
			}
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

func get404(url string, pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url {
				errorHandler(w, r, http.StatusNotFound, "")
				return
			}
			pageFunc(w, r)
		})
}

func getRedirectPermanent(url string, pageFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			//Don't accept post or put requests
			if r.Method != get {
				http.Redirect(w, r, url, http.StatusSeeOther)
			}

			pageFunc(w, r)
		})
	//Redirects back to subdirectory "url". Needed when url parameters are not wanted or needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url+"/", http.RedirectHandler(url, http.StatusMovedPermanently))
}

func getParameter(url string, pageFunc func(http.ResponseWriter, *http.Request, string), regex *regexp.Regexp) {
	var parameters, lowerParams string
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			//Don't accept post or put requests
			if r.Method != get {
				http.Redirect(w, r, url, http.StatusSeeOther)
			}

			parameters = strings.TrimPrefix(r.URL.Path, url)
			lowerParams = strings.ToLower(parameters)

			if parameters != lowerParams {
				//Redirect to page with lowercase parameters.
				http.Redirect(w, r, url+lowerParams, http.StatusSeeOther)
				return
			}

			if regex.MatchString(lowerParams) {
				pageFunc(w, r, lowerParams)
				return
			}
			errorType := "event"
			if url == urlClub {
				errorType = "club"
			}
			errorHandler(w, r, http.StatusNotFound, errorType)
		})
}

func getParameters(url string, pageFunc func(http.ResponseWriter, *http.Request, string, string), regex *regexp.Regexp) {
	var parameters, lowerParams string
	var ids []string
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			//Don't accept post or put requests
			if r.Method != get {
				http.Redirect(w, r, url, http.StatusSeeOther)
			}

			parameters = strings.TrimPrefix(r.URL.Path, url)
			lowerParams = strings.ToLower(parameters)

			if parameters != lowerParams {
				//Redirect to page with lowercase parameters.
				http.Redirect(w, r, url+lowerParams, http.StatusSeeOther)
				return
			}

			if regex.MatchString(lowerParams) {
				ids = strings.Split(lowerParams, "/")
				pageFunc(w, r, ids[0], ids[1])
				return
			}
			errorType := "event"
			if url == urlClub {
				errorType = "club"
			}
			errorHandler(w, r, http.StatusNotFound, errorType)
		})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, errorType string) {
	//func errorHandler(gzw gzipResponseWriter, r *http.Request, status int, errorType string) {
	//All EventBucket page urls and ids are lowercase
	lowerURL := strings.ToLower(r.URL.Path)

	//Redirect if url contains any uppercase letters.
	if r.URL.Path != lowerURL {
		http.Redirect(w, r, lowerURL, http.StatusSeeOther)
		return
	}
	lowerURL = strings.TrimSuffix(r.URL.Path, "/")

	//check if the request matches any of the pages that don't require parameters
	if strings.Count(lowerURL, "/") >= 2 {
		for _, page := range []string{urlAbout, urlArchive, urlClubs, urlLicence, urlShooters} {
			if strings.HasPrefix(lowerURL, page) {
				//redirect to page without parameters
				http.Redirect(w, r, page, http.StatusSeeOther)
				return
			}
		}
	}
	w.WriteHeader(status)
	templater(w, page{
		Title: "Error",
		Data: map[string]interface{}{
			"Type": errorType,
		},
	})
}

func formError(w http.ResponseWriter, r *http.Request, f form, err error) {
	f.Error = err
	setSession(w, f)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

/*//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainExpiresTime() {
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		setExpiresTime()
	}
}

//Set expiry date 1 year, 0 months & 0 days in the future.
func setExpiresTime() {
	//Date format is the same as Go`s time.RFC1123 but uses "GMT" timezone instead of "UTC" time standard.
	cacheExpires = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
	//w3.org: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
}*/
