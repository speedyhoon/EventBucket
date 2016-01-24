package main

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	dirRoot         = "./"
	dirGzip         = "dirGzip"
	urlHome         = "/"
	urlAbout        = "/about"
	urlArchive      = "/archive"
	urlClubs        = "/clubs"
	urlClubSettings = "/club-settings/"
	urlEvents       = "/events"
	urlLicence      = "/licence"
	urlShooters     = "/shooters"
	//GET with PARAMETERS
	urlEvent         = "/event/"          //eventID
	urlEventSettings = "/event-settings/" //eventID
	urlClub          = "/club/"           //clubID

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

/*TODO if no parameters provided, keep user on the same page but didplay when they need to provide in order for the page to work.
not doing this may frustrate some users who want to get to the club settings page but can't remember the club id.
then display a list of clubs and status code 404
*/
func getParameters(url string, pageFunc func(http.ResponseWriter, *http.Request, string), regex *regexp.Regexp) {
	var parameters, lowerParams string
	http.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
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
			whoops(w, r, url)
		})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	//All EventBucket page urls and ids are lowercase
	lowerURL := strings.ToLower(strings.TrimSuffix(r.URL.Path, "/"))

	//prevents a redirect loop if url is already in lowercase letters.
	if r.URL.Path != lowerURL {

		//check if the request matches any of the pages that don't require parameters
		if strings.Count(lowerURL, "/") >= 2 {
			for _, page := range []string{urlAbout, urlArchive, urlClubs /*urlEvent,*/, urlEvents, urlLicence, urlShooters} {
				if strings.HasPrefix(lowerURL, page) {
					//redirect to page without parameters
					http.Redirect(w, r, page, http.StatusSeeOther)
					return
				}
			}
		}
		http.Redirect(w, r, lowerURL, http.StatusSeeOther)
		return
	}
	w.WriteHeader(status)
	templater(w, page{
		Title: "Error",
		Data: M{
			//		"Status": "404 Page Not Found",
			"Status": status,
		},
	})
}