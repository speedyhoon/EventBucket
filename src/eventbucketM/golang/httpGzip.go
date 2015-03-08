package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func start() {
	go startDatabase()
	serveDir(DIR_JS)
	serveDir(DIR_CSS)
	serveDir(DIR_PNG)
	serveDir(DIR_JPEG)
	serveDir(DIR_SVG)
	serveDir(DIR_WOF)
	serveDir(DIR_WOF2)

	GetRedirectPermanent(URL_about, about)
	GetRedirectPermanent(URL_archive, archive)
	GetRedirectPermanent(URL_licence, licence)
	GetRedirectPermanent(URL_shooters, shooters)
	GetRedirectPermanent(URL_clubs, clubs)

	GetParameters(URL_event, event)
	GetParameters(URL_club, club)
	GetParameters(URL_eventSettings, eventSettings)
	GetParameters(URL_scoreboard, scoreboard)
	GetParameters(URL_eventShotsNSighters, eventShotsNSighters)
	GetParameters(URL_startShooting, startShooting)
	GetParameters(URL_startShootingAll, startShootingAll)
	GetParameters(URL_totalScores, totalScores)
	GetParameters(URL_totalScoresAll, totalScoresAll)
	//	GetParameters(URL_rangeReport, range_report)

	Post(URL_eventInsert, eventInsert)
	Post(URL_queryShooterList, searchShooter)
	Post(URL_queryShooterGrade, searchShooterGrade)
	Post(URL_updateShooterList, PostVia(nraaStartUpdateShooterList, URL_shooters))
	Post(URL_clubInsert, PostVia(clubInsert, URL_clubs)) //TODO redirect to actual club created
	Post(URL_updateRange, rangeUpdate)
	Post(URL_eventRangeInsert, rangeInsert)
	Post(URL_eventAggInsert, aggInsert)
	Post(URL_shooterInsert, shooterInsert)
	Post(URL_shooterListInsert, shooterListInsert)
	//	Post(URL_updateTotalScores, updateTotalScores)
	Post(URL_updateShotScores, updateShotScores)
	Post(URL_updateSortScoreBoard, updateSortScoreBoard)
	Post(URL_updateEventGrades, updateEventGrades)
	Post(URL_updateIsPrizeMeet, updateIsPrizeMeet)
	//	Post(URL_champInsert, PostVia(champInsert, URL_championships))	//TODO redirect to actual club created
	//Club insert/update
	Post(URL_clubMoundInsert, clubMoundInsert)
	Post(URL_clubDetailsUpsert, clubDetailsUpsert)

	Get(URL_home, home)
}

var (
	VURL_home                = regexp.MustCompile("^/$")
	VURL_event               = regexp.MustCompile("^" + URL_event + "([" + ID_CHARSET_REGEX + "]+)$")
	VURL_eventShotsNSighters = regexp.MustCompile("^" + URL_eventShotsNSighters + "([" + ID_CHARSET_REGEX + "]+)$")
	VURL_club                = regexp.MustCompile("^" + URL_club + "([" + ID_CHARSET_REGEX + "]+)$")
	VURL_eventSettings       = regexp.MustCompile("^" + URL_eventSettings + "([" + ID_CHARSET_REGEX + "]+)$")
	VURL_scoreboard          = regexp.MustCompile("^" + URL_scoreboard + "([" + ID_CHARSET_REGEX + "]+)$")
)

type Page struct {
	TemplateFile, Title string
	Theme               string
	Data                M
	v8Url               *regexp.Regexp
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(h http.Handler, w http.ResponseWriter, r *http.Request) {
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

func Get(url string, runner func() Page) {
	//Setup url as a subdirectory path. e.g. if url = "foobar" then "http://localhost/foobar" is available
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) { templator(runner(), w, r) }))
}
func GetRedirectPermanent(url string, runner func() Page) {
	Get(url, runner)
	//Redirects back to subdirectory "url". Needed when url parameters are not wanted/needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url+"/", http.RedirectHandler(url, http.StatusMovedPermanently))
}
func GetParameters(url string, runner func(string) Page) {
	//TODO add a way to get the event id, club id and range id from the url
	h := func(w http.ResponseWriter, r *http.Request) { templator(runner(getIdFromUrl(r, url)), w, r) }
	http.Handle(url, serveHtml(h))
}

func Post(url string, runner http.HandlerFunc) {
	//TODO Post - Post to endpoint. If valid return to X page, else return to previous page with form filled out and wrong values highlighted with error message
	//TODO Add ajax detection to so ajax requests are not redirected back to the referrer page.
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

func httpHeaders(w http.ResponseWriter, setHeaders []string) {
	//TODO Only set CSP when not in debug mode
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'") //TODO remove unsafe inline when start shooting gets its settings a different way
	headers := map[string][2]string{
		"expire":    {"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //RESEARCH should it return GMT time?  //Expiry date is in 1 year, 0 months & 0 days in the future
		"cache":     {"Vary", "Accept-Encoding"},
		"public":    {"Cache-Control", "public"},
		"gzip":      {"Content-Encoding", "gzip"},
		"html":      {"Content-Type", "text/html; charset=utf-8"},
		"noCache":   {"Cache-Control", "no-cache, no-store, must-revalidate"},
		"expireNow": {"Expires", "0"},
		"pragma":    {"Pragma", "no-cache"},
		DIR_CSS:     {"Content-Type", "text/css; charset=utf-8"},
		DIR_JS:      {"Content-Type", "text/javascript"},
		DIR_PNG:     {"Content-Type", "image/png"},
		DIR_JPEG:    {"Content-Type", "image/jpeg"},
		DIR_SVG:     {"Content-Type", "image/svg+xml"},
		DIR_WOF:     {"Content-Type", "application/font-woff"},
		DIR_WOF2:    {"Content-Type", "application/font-woff2"},
	}
	for _, lookup := range setHeaders {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}

func getIdFromUrl(r *http.Request, pageUrl string) string {
	//TODO add automatic validation checking for id using regex pattens
	return r.URL.Path[len(pageUrl):]
}
