package main

import (
	"net/http"
	"mgo"
	"strings"
	"compress/gzip"
	"io"
//	"github-spdy"
)

var conn *mgo.Database

const (
	MINIFY = false
	DEV = true	//turn on dev warnings, E.g. Template errors
)

func main() {
	conn = DB()

	changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Set some header.
//			w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';")
			w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
			w.Header().Set("Vary", "Accept-Encoding")

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				dump("request does not want gzip")
				// Serve with the actual handler.
				h.ServeHTTP(w, r)
			}else{
				dump("request wants GZIP")

				w.Header().Set("Content-Encoding", "gzip")
//				w.Header().Set("Content-Type", "text/css")
				gz := gzip.NewWriter(w)
				defer gz.Close()
				gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
				//		pageController(gzr, r)
				h.ServeHTTP(gzr, r)

//				comp, _ := gzip.NewWriter(c)
//				defer comp.Close()
//
//				c.SetHeader("Content-Encoding", "gzip")
//
//				buf := make([]byte, 1048576)
//				for {
//					nbytes, err := file.Read(buf)
//					if err != nil {
//						return
//					}
//					comp.Write(buf[0:nbytes])
//				}
			}
		}
	}

	http.Handle("/i/", changeHeaderThenServe(http.FileServer(http.Dir("./"))))
	http.Handle("/css/", changeHeaderThenServe(http.FileServer(http.Dir("./"))))

	//File Server
	//TODO override http headers to set type for each folder
	//TODO make sure all resources don't have a . extension to save network bandwidth
					//	http.Handle("/i/", httpHeadersSet_2(http.FileServer(http.Dir("./"))))
					//	http.Handle("/css/", httpHeadersSet_2(http.FileServer(http.Dir("./"))))

	//GET
	http.HandleFunc("/", httpHeadersSet(home))
//	http.HandleFunc("/organisers", httpHeadersSet(organisers))
//	http.HandleFunc("/organisers/", redirectPermanent("/organisers"))
//	//	http.HandleFunc("/clubs", clubs)
//	//	http.HandleFunc("/startShooting", startShooting)
//	http.HandleFunc("/events", httpHeadersSet(events))
//	http.HandleFunc("/events/", redirectPermanent("/events"))
//	http.HandleFunc("/event/", httpHeadersSet(event))
//	http.HandleFunc("/club/", httpHeadersSet(club))
//	//	http.HandleFunc("/eventSetup", eventSetup)
//	http.HandleFunc("/eventSettings/", httpHeadersSet(eventSettings))
//	http.HandleFunc("/startShooting/", httpHeadersSet(startShooting))
//	http.HandleFunc("/totalScores/", httpHeadersSet(totalScores))
//	http.HandleFunc("/scoreboard/", httpHeadersSet(scoreboard))
//
//	//POST
//	http.HandleFunc("/clubInsert", redirectVia(clubInsert, "/organisers"))
//	http.HandleFunc("/eventInsert", redirectVia(eventInsert, "/organisers"))
//	http.HandleFunc("/champInsert", redirectVia(champInsert, "/organisers"))
//	http.HandleFunc("/rangeInsert", httpHeadersSet(rangeInsert))
//	http.HandleFunc("/shooterInsert", httpHeadersSet(shooterInsert))
//	//	http.HandleFunc("/rangeInsert", redirectTo("eventSettings/5317198d81a27b3006872b0f"))
//	http.HandleFunc("/clubMoundInsert", httpHeadersSet(clubMoundInsert))
//	http.HandleFunc("/updateTotalScores", httpHeadersSet(updateTotalScores))
//	http.HandleFunc("/dummyTest", testNewCreationStructSaving)
//	http.HandleFunc("/customSchema", customSchema)
//	http.HandleFunc("/try", httpHeadersSet(try))
	http.ListenAndServe(":80", nil)
}


func httpHeadersSet(pageController http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
//		http://content-security-policy.com/
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';")
		w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
		w.Header().Set("Vary", "Accept-Encoding")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			pageController(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		pageController(gzr, r)
	}
}

func httpHeadersSet_2(pageController http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//	http://content-security-policy.com/
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self';")
		w.Header().Set("Expires", "Thu, 16 Apr 2015 20:00:00 GMT")
		w.Header().Set("Vary", "Accept-Encoding")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			pageController.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		pageController.ServeHTTP(gzr, r)
	}
}

//changeHeaderThenServe := func(h http.Handler) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		// Set some header.
//		w.Header().Add("Keep-Alive", "300")
//		// Serve with the actual handler.
//		h.ServeHTTP(w, r)
//	}
//}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
