package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//TODO maybe add page redirects for pages. replace uppercase urls with lowercase
//Is this needed or not??????????????????????????
/*func redirectToUppercase() {
	//listOfPages := map[string]*func(){
	//	urlEventSettings: eventSettings,
	//}
	listOfPages := map[string]string{
		urlEventSettings: strings.ToLower(urlEventSettings),
	}
	for from, to := range listOfPages {
		GetRedirectPermanentTo(from, to)
	}
}
func GetRedirectPermanentTo(from, to string) {
	http.Handle(from, http.RedirectHandler(to, http.StatusMovedPermanently))
}*/

func start() {
	go startDatabase()
	serveDir(dirJS)
	serveDir(dirCSS)
	serveDir(dirPNG)
	serveDir(dirJPEG)
	serveDir(dirSVG)
	serveDir(dirWOF)
	serveDir(dirWOF2)

	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, archive)
	getRedirectPermanent(urlLicence, licence)
	getRedirectPermanent(urlShooters, shooters)
	getRedirectPermanent(urlClubs, clubs)

	getParameters(urlEvent, event)
	getParameters(urlClub, club)
	getParameters(urlEventSettings, eventSettings)
	getParameters(urlScoreboard, scoreboard)
	getParameters(urlEventShotsNSighters, eventShotsNSighters)
	getParameters(urlStartShooting, startShooting)
	getParameters(urlStartShootingAll, startShootingAll)
	getParameters(urlTotalScores, totalScores)
	getParameters(urlTotalScoresAll, totalScoresAll)
	//	GetParameters(urlRangeReport, rangeReport)

	//Event
	post(urlEventInsert, eventInsert)
	post(urlEventUpdateShooter, eventUpdateShooter)
	post(urlQueryShooterList, searchShooter)
	post(urlQueryShooterGrade, searchShooterGrade)
	//Nraa
	post(urlUpdateShooterList, postVia(nraaStartUpdateShooterList, urlShooters))
	post(urlClubInsert, postVia(clubInsert, urlClubs)) //TODO redirect to actual club created
	post(urlUpdateRange, rangeUpdate)
	post(urlEventRangeInsert, rangeInsert)
	post(urlEventAggInsert, aggInsert)
	post(urlShooterInsert, shooterInsert)
	post(urlShooterListInsert, shooterListInsert)
	//Total Scores
	post(urlUpdateTotalScores, updateTotalScores)
	//Startshooting
	post(urlUpdateShotScores, updateShotScores2)
	post(urlUpdateSortScoreBoard, updateSortScoreBoard)
	post(urlUpdateEventGrades, updateEventGrades)
	post(urlUpdateIsPrizeMeet, updateIsPrizeMeet)
	//	Post(urlChampInsert, PostVia(champInsert, urlChampionships))	//TODO redirect to actual club created
	//Club insert/update
	post(urlClubMoundInsert, clubMoundInsert)
	post(urlClubDetailsUpsert, clubDetailsUpsert)

	get(urlHome, home)
}

var (
	// VURLHome should be unexported
	VURLHome = regexp.MustCompile("^/$")
	// VURLEventShotsNSighters should be unexported
	VURLEventShotsNSighters = regexp.MustCompile("^" + urlEventShotsNSighters + "([" + idCharsetRegex + "]+)$")
	//VURLEvent               = regexp.MustCompile("^" + urlEvent + "([" + idCharsetRegex + "]+)$")
	//VURLClub                = regexp.MustCompile("^" + urlClub + "([" + idCharsetRegex + "]+)$")
	//VURLEventSettings       = regexp.MustCompile("^" + urlEventSettings + "([" + idCharsetRegex + "]+)$")
	//VURLScoreboard          = regexp.MustCompile("^" + urlScoreboard + "([" + idCharsetRegex + "]+)$")
)

// Page is exported
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

func gzipper(h http.Handler, w http.ResponseWriter, r *http.Request) {
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

func get(url string, runner func() Page) {
	//Setup url as a subdirectory path. e.g. if url = "foobar" then "http://localhost/foobar" is available
	http.Handle(url, serveHtml(func(w http.ResponseWriter, r *http.Request) { templator(runner(), w, r) }))
}
func getRedirectPermanent(url string, runner func() Page) {
	get(url, runner)
	//Redirects back to subdirectory "url". Needed when url parameters are not wanted/needed.
	//e.g. if url = "foobar" then "http://localhost/foobar/fdsa" will redirect to "http://localhost/foobar"
	http.Handle(url+"/", http.RedirectHandler(url, http.StatusMovedPermanently))
}
func getParameters(url string, runner func(string) Page) {
	//TODO add a way to get the event id, club id and range id from the url
	h := func(w http.ResponseWriter, r *http.Request) { templator(runner(getIDFromURL(r, url)), w, r) }
	http.Handle(url, serveHtml(h))
}

func post(url string, runner http.HandlerFunc) {
	//TODO Post - Post to endpoint. If valid return to X page, else return to previous page with form filled out and wrong values highlighted with error message
	//TODO Add ajax detection to so ajax requests are not redirected back to the referrer page.
	http.HandleFunc(url, serveHtml(runner))
}

func postVia(runThisFirst func(http.ResponseWriter, *http.Request), url string) func(http.ResponseWriter, *http.Request) {
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
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; font-src 'self'") //TODO remove unsafe inline when start shooting gets its settings a different way
	w.Header().Set("X-Frame-Options", "DENY")                                                                                                                                       //developer.mozilla.org/en-US/docs/Web/HTTP/X-Frame-Options
	headers := map[string][2]string{
		"expire":    {"Expires", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC1123)}, //RESEARCH should it return GMT time?  //Expiry date is in 1 year, 0 months & 0 days in the future
		"cache":     {"Vary", "Accept-Encoding"},
		"public":    {"Cache-Control", "public"},
		"gzip":      {"Content-Encoding", "gzip"},
		"html":      {"Content-Type", "text/html; charset=utf-8"},
		"noCache":   {"Cache-Control", "no-cache, no-store, must-revalidate"},
		"expireNow": {"Expires", "0"},
		"pragma":    {"Pragma", "no-cache"},
		dirCSS:      {"Content-Type", "text/css; charset=utf-8"},
		dirJS:       {"Content-Type", "text/javascript"},
		dirPNG:      {"Content-Type", "image/png"},
		dirJPEG:     {"Content-Type", "image/jpeg"},
		dirSVG:      {"Content-Type", "image/svg+xml"},
		dirWOF:      {"Content-Type", "application/font-woff"},
		dirWOF2:     {"Content-Type", "application/font-woff2"},
	}
	for _, lookup := range setHeaders {
		w.Header().Set(headers[lookup][0], headers[lookup][1])
	}
}

func getIDFromURL(r *http.Request, pageURL string) string {
	//TODO add automatic validation checking for id using regex pattens
	return r.URL.Path[len(pageURL):]
}
