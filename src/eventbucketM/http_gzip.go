package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func file_headers_n_gzip(h http.Handler, content_type string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
		http_headers(w, []string{"expire", "cache", content_type})
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(gzr, r)
		} else {
			h.ServeHTTP(w, r)
			fmt.Println("This request does not support gzip")
			//			Info.Println("This request does not support gzip")
		}
	}
}

func pageTemplate(url string)string{
	switch url{
	case "/", "home", URL_about, URL_licence, /*URL_licence_summary,*/ URL_archive:
		return TEMPLATE_HOME
	case URL_scoreboard:
		return TEMPLATE_EMPTY
	}
	return TEMPLATE_ADMIN
}

func gzipper(url string, runner func()M) http.Handler {
	if dev_mode_DEBUG{
		return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			/*http_headers(w, []string{"html", "nocache0", "nocache1", "nocache2"})
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				defer gz.Close()
				gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
				h.ServeHTTP(gzr, r)
			} else {
				h.ServeHTTP(w, r)
				fmt.Println("This request does not support gzip")
				//			Info.Println("This request does not support gzip")
			}*/
			templator(pageTemplate(url), url, runner(), w)
		}))
	}else{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			/*http_headers(w, []string{"html", "nocache0", "nocache1", "nocache2"})
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				defer gz.Close()
				gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
				h.ServeHTTP(gzr, r)
			} else {
				h.ServeHTTP(w, r)
				fmt.Println("This request does not support gzip")
				//			Info.Println("This request does not support gzip")
			}*/
			templator(pageTemplate(url), url, runner(), w)
		})
	}
}




//func get(url string, runner func()Page){
//	http.Handle(url, gzipPage(url, runner))
//}

//func gzipPage(url string, runner func()Page) http.Handler {
func get(url string, runner func()Page) {
	pageName := url
	if url == "/"{
		pageName = "home"
	}
	if dev_mode_DEBUG{
		http.Handle(url,
			http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer dev_mode_timeTrack(time.Now(), r.RequestURI)
				templatePage(url, runner(), w)
			})))
	}else{
		http.Handle(url,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer dev_mode_timeTrack(time.Now(), r.RequestURI)
				templatePage(pageName, runner(), w)
			}))
	}
}

func gziParameters(url string, runner func(string)M) http.Handler {
	if dev_mode_DEBUG{
		return http.HandlerFunc(agent.WrapHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			//TODO add redirect to differnet page is no parameters  OR add error message on page???
			//		dump("r.URL")
			//		export(r.URL)
			//		dump("url")
			//		export(url)
			//		dump("output")
			//		export(get_id_from_url(r, url))
			templator(pageTemplate(url), url, runner(get_id_from_url(r, url)), w)
		}))
	}else{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer dev_mode_timeTrack(time.Now(), r.RequestURI)
			//TODO add redirect to differnet page is no parameters  OR add error message on page???
			//		dump("r.URL")
			//		export(r.URL)
			//		dump("url")
			//		export(url)
			//		dump("output")
			//		export(get_id_from_url(r, url))
			templator(pageTemplate(url), url, runner(get_id_from_url(r, url)), w)
		})
	}
}

func html_headers_n_gzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer dev_mode_timeTrack(time.Now(), r.RequestURI)
//		http_headers(w, []string{"html", "nocache0", "nocache1", "nocache2"})
		http_headers(w, []string{"html", "nocache"})
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(gzr, r)
		} else {
			h.ServeHTTP(w, r)
			fmt.Println("This request does not support gzip")
			//			Info.Println("This request does not support gzip")
		}
	}
}

func http_headers(w http.ResponseWriter, set_headers []string) {
	//	mins_in_year := 525949	//	hours_in_year := 8765.81
	//	one_year := time.Duration(mins_in_year)*time.Minute
	headers := map[string][2]string{
		"expire": [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO it should return GMT time I think
		//		"expire":[2]string{"Expires", time.Now().UTC().Add(one_year).Format(time.RFC1123)},//TODO it should return GMT time I think
		//		"0cache":[2]string{"Expires", time.Now().UTC().Format(time.RFC1123)},//TODO it should return GMT time I think
//		"nocache0": [2]string{"Cache-Control", "no-cache, no-store, must-revalidate"},
//		"nocache1": [2]string{"Expires", "0"},
//		"nocache2": [2]string{"Pragma", "no-cache"},
		"cache":    [2]string{"Vary", "Accept-Encoding"},
//		"csp":      [2]string{"Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';"}, //content-security-policy.com
		"gzip":     [2]string{"Content-Encoding", "gzip"},
		"html":     [2]string{"Content-Type", "text/html; charset=utf-8"},
		"css":      [2]string{"Content-Type", "text/css; charset=utf-8"},
		//TODO which mime type is best for javascript?
		//"js":		[2]string{"Content-Type", "application/javascript"},
		"js":  [2]string{"Content-Type", "text/javascript"},
		"png": [2]string{"Content-Type", "image/png"},
		"jpg": [2]string{"Content-Type", "image/jpeg"},
		"gif": [2]string{"Content-Type", "image/gif"},
		"webp": [2]string{"Content-Type", "image/webp"},
		"svg":  [2]string{"Content-Type", "image/svg+xml"},
	}
	w.Header().Set("Content-Security-Policy", "Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';")
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

func redirectPermanent(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently) //Search engine Optimisation
	}
}

func redirectVia(runThisFirst func(http.ResponseWriter, *http.Request), path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		runThisFirst(w, r)
		http.Redirect(w, r, path, http.StatusSeeOther) //303 mandating the change of request type to GET
	}
}

func redirecter(path string, w http.ResponseWriter, r *http.Request) {
	//	return func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, path, http.StatusSeeOther)
	//	}
}

func get_id_from_url(r *http.Request, page_url string) string {
	//TODO add validation checking for id using regex pattens
	/*TODO add a http layer function between p_page functions and main.go so that the event_id or club_id can
	be validated and the p_page functions don't have to interact with http at all*/
	//	var validID = regexp.MustCompile(`\A` + page_url + `[0-9a-f]{24}\z`)
//	url := fmt.Sprintf("%v", r.URL)
	//	if validID.MatchString(url) {
	//		templator("admin", eventSettings_HTML(), eventSettings_Data(url[len(page_url):]), w)
	//	}else {
	//		redirectPermanent("/events")
	//		fmt.Println("redirected user " + url)
	//	}


//	return url[len(page_url):]
	return r.URL.Path[len(page_url):]
}

//A better solution to gzip over http!
//package main
//
//import (
//"compress/gzip"
//"flag"
//"io"
//"http"
//"log"
//"os"
//"path"
//"strings"
//)
//
//var rootdir *string = flag.String("rootdir", "/home/pkf/intraday", "chroot to this directory.")
//var prefix *string = flag.String("prefix", "/intraday/", "prefix path in URLs")
//
//func checkencoding(req *http.Request) bool {
//	encoding := req.Header["Accept-Encoding"]
//	if encoding != "" {
//		for _, v := range strings.Split(encoding, ",", 10) {
//			if strings.TrimSpace(v) == "gzip" {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//// send a file, with optional compression
//func handler(c *http.Conn, req *http.Request) {
//
//	// we only support GET
//	if req.Method != "GET" {
//		log.Stderrf("req.Method is %s", req.Method)
//		return
//	}
//
//	// should we compress?
//	compress := checkencoding(req)
//	log.Stderrf("%s; compression requested: %v", req.URL.Path, compress)
//
//	// clean the path and make sure the file exists
//	// handles the case where people try to fetch /../../../../etc/passwd or something
//	cpath := path.Clean(req.URL.Path)
//	if !strings.HasPrefix(cpath, *prefix) {
//		http.NotFound(c, req)
//		return
//	}
//
//	cpath = cpath[len(*prefix):]
//
//	file, _ := os.Open(cpath, os.O_RDONLY, 0)
//	if file == nil {
//		http.NotFound(c, req)
//		return
//	}
//	defer file.Close()
//
//	// write it out
//	c.SetHeader("Content-Type", "text-plain; charset=us-ascii")
//
//	if compress {
//		comp, _ := gzip.NewWriter(c)
//		defer comp.Close()
//
//		c.SetHeader("Content-Encoding", "gzip")
//
//		buf := make([]byte, 1048576)
//		for {
//			nbytes, err := file.Read(buf)
//			if err != nil {
//				return
//			}
//			comp.Write(buf[0:nbytes])
//		}
//	} else {
//		io.Copy(c, file)
//	}
//}
//
//func main() {
//	flag.Parse()
//	err := os.Chdir(*rootdir)
//	if err != nil {
//		log.Stderr(err.String())
//		os.Exit(1)
//	}
//	log.Stderrf("starting in %s, expecting urls to begin with %s", *rootdir, *prefix)
//	http.ListenAndServe(":12345", http.HandlerFunc(handler))
//}
