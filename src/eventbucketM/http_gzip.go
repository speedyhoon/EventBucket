package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
	"regexp"
)

const (
//GET
	URL_home           	= "/"
	URL_about           	= "/about"
	URL_clubs            = "/clubs"
	URL_licence         	= "/licence"
	//URL_licence_summary	= "/licence-summary"
	URL_archive			  	= "/archive"
	URL_organisers       = "/organisers"
	URL_event            = "/event/"							//event Id special type characters only allowed
//GET with PARAMETERS
	URL_club                 = "/club/"
	//URL_events               = "/events/"
	URL_eventSettings        = "/eventSettings/"			//event id
	URL_scoreboard           = "/scoreboard/"				//event id/range_id
	URL_totalScores          = "/totalScores/"			//event id/range_id
	URL_totalScoresAll       = "/totalScoresAll/"		//event id/range_id
	URL_startShooting        = "/startShooting/"			//event id/range_id
	URL_startShootingAll     = "/startShootingAll/"		//event id/range_id
	URL_queryShooterList 	 = "/queryShooterList"
//POST
	URL_clubInsert           = "/clubInsert"
	//	URL_champInsert          = "/champInsert"
	URL_eventInsert          = "/eventInsert"
	//	URL_eventInsert2         = "/eventInsert2"
	URL_eventRangeInsert     = "/rangeInsert"
	URL_eventAggInsert       = "/aggInsert"
	URL_shooterInsert        = "/shooterInsert"
	URL_shooterListInsert    = "/shooterListInsert"
	URL_updateSortScoreBoard = "/updateSortScoreBoard"
	URL_updateTotalScores    = "/updateTotalScores"
	URL_updateShotScores     = "/updateShotScores"
	URL_updateEventGrades    = "/updateEventGrades"
	URL_updateEventName      = "/updateEventName/"
	URL_updateRange          = "/updateRange"
	URL_updateIsPrizeMeet    = "/updateIsPrizeMeet"
	//	URL_dateUpdate           = "/dateUpdate/"
	URL_club_mound_update    = "/clubMoundUpdate/"
	URL_clubMoundInsert      = "/clubMoundInsert/"
	URL_clubDetailsUpsert      = "/clubDetailsUpsert/"
	URL_updateShooterList    = "/updateShooterList"
	URL_eventShotsNSighters  = "/eventShotsNSighters/"
	//	URL_rangeReport          = "/rangeReport/"
	URL_randomData           = "/random-data/"
)

func start() {
	go DB()
	serveDir(DIR_JS)
	serveDir(DIR_CSS)
	serveDir(DIR_PNG)
	serveDir(DIR_JPEG)
	serveDir(DIR_SVG)

	GetRedirectPermanent(URL_about, about)
	GetRedirectPermanent(URL_archive, archive)
	GetRedirectPermanent(URL_licence, licence)
	GetRedirectPermanent(URL_organisers, organisers)
	GetRedirectPermanent(URL_clubs, clubs)
	//	GetRedirectPermanent(URL_events, events)

	GetParameters(URL_event, event)
	GetParameters(URL_club, club)
	GetParameters(URL_eventSettings, eventSettings)
	GetParameters(URL_scoreboard, scoreboard)
	GetParameters(URL_eventShotsNSighters, eventShotsNSighters)
	//	GetParameters(URL_startShooting, startShooting)
	//	GetParameters(URL_startShootingAll, startShootingAll)
	//	GetParameters(URL_totalScores, totalScores)
	//	GetParameters(URL_totalScoresAll, totalScoresAll)
	//	GetParameters(URL_rangeReport, range_report)

	Post(URL_eventInsert, eventInsert)
	Post(URL_queryShooterList, queryShooterList)	//Search for a shooter by first, surname & club
	Post(URL_updateShooterList, PostVia(updateShooterList, URL_organisers))
	Post(URL_clubInsert, PostVia(clubInsert, URL_organisers))
	Post(URL_updateRange, rangeUpdate2)
	//	Post(URL_dateUpdate, dateUpdate)
	Post(URL_eventRangeInsert, rangeInsert)
	Post(URL_eventAggInsert, aggInsert)
	Post(URL_shooterInsert, shooterInsert)
	Post(URL_shooterListInsert, shooterListInsert)
	//	Post(URL_updateTotalScores, updateTotalScores)
	Post(URL_updateShotScores, updateShotScores)
	Post(URL_updateSortScoreBoard, updateSortScoreBoard)
	//	Post(URL_updateEventName, updateEventName)
	Post(URL_updateEventGrades, updateEventGrades)
	Post(URL_updateIsPrizeMeet, updateIsPrizeMeet)
	//	Post(URL_champInsert, PostVia(champInsert, URL_organisers))
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

func Gzip(h http.Handler, w http.ResponseWriter, r *http.Request){
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
	//Setup url as a subdirectory path. e.g. if url = "foobar" then "http://localhost/foobar" is available
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) {templator(runner(), w, r)}))
}
func GetRedirectPermanent(url string, runner func()Page){
	Get(url, runner)
	//Setup redirect back to subdirectory "url". Needed when url parameters are not wanted/needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url + "/", http.RedirectHandler(url, http.StatusMovedPermanently))
}
func GetParameters(url string, runner func(string)Page) {
	//TODO add a way to get the event id, club id and range id from the url
	h := func(w http.ResponseWriter, r *http.Request) {templator(runner(getIdFromUrl(r, url)), w, r)}
	http.Handle(url, serveHtml(h))
}
//TODO Post - Post to endpoint. If valid return to X page, else return to previous page with form filled out and wrong values highlighted with error message
//TODO Add ajax detection to so ajax requests are not redirected back to the referrer page.
func Post(url string, runner http.HandlerFunc){
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

func httpHeaders(w http.ResponseWriter, set_headers []string) {
	//TODO Only set CSP when not in debug mode
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self'; img-src 'self' data:")
	headers := map[string][2]string{
		"expire":  [2]string{"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //TODO should it return GMT time?  //Expiry date is in 1 year, 0 months & 0 days in the future
		"cache":   [2]string{"Vary", "Accept-Encoding"},
		"public":  [2]string{"Cache-Control", "public"},
		"gzip":    [2]string{"Content-Encoding", "gzip"},
		"html":    [2]string{"Content-Type", "text/html; charset=utf-8"},
		"noCache": [2]string{"Cache-Control", "no-cache, no-store, must-revalidate"},
		"expireNow":[2]string{"Expires", "0"},
		"pragma":  [2]string{"Pragma", "no-cache"},
		DIR_CSS:   [2]string{"Content-Type", "text/css; charset=utf-8"},
		DIR_JS:    [2]string{"Content-Type", "text/javascript"},
		DIR_PNG:   [2]string{"Content-Type", "image/png"},
		DIR_JPEG:  [2]string{"Content-Type", "image/jpeg"},
//		DIR_WEBP:  [2]string{"Content-Type", "image/webp"},
		DIR_SVG:   [2]string{"Content-Type", "image/svg+xml"},
	}
	for _, lookup := range set_headers {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func getIdFromUrl(r *http.Request, page_url string) string {
	//TODO add validation checking for id using regex pattens
	//TODO add a http layer function between p_page functions and main.go so that the event_id or club_id can be validated and the p_page functions don't have to interact with http at all
	return r.URL.Path[len(page_url):]
}

type Page struct{
	TemplateFile, Title string
	Theme string
	Data M
	v8Url *regexp.Regexp
}
