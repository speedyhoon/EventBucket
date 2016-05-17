package main

import (
	"fmt"
	"net/http"
	"regexp"
)

const (
	get         = "GET"
	pst         = "POST"
	dirCSS      = "dirCSS"
	dirJS       = "dirJS"
	dirPNG      = "dirPNG"
	dirGIF      = "dirGIF"
	dirSVG      = "dirSVG"
	urlEvents   = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlLicence  = "/licence"
	urlShooters = "/shooters"
	//GET with PARAMETERS
	urlClub            = "/club/"             //clubID
	urlEntries         = "/entries/"          //eventID
	urlEntryList       = "/print-entry-list/" //eventID
	urlEventSettings   = "/event-settings/"   //eventID
	urlEventReport     = "/event-report/"     //eventID
	urlScoreboard      = "/scoreboard/"       //eventID
	urlEnterShots      = "/enter-shots/"      //eventID
	urlEnterShotsAll   = "/enter-shots-all/"  //eventID
	urlPrintScorecards = "/print-cards/"      //eventID/shooterID
	urlEnterTotals     = "/enter-totals/"     //eventID
	urlEnterTotalsAll  = "/enter-totals-all/" //eventID
)

var (
	//URL validation matching
	regexID      = regexp.MustCompile(`^[a-z0-9]+$`)
	regexPath    = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+$`)
	regexBarcode = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+#[a-z0-9]+$`)
)

func pages() {
	serveFile("favicon.ico")
	serveFile("robots.txt")
	serveDir(dirCSS, "./cz/")
	serveDir(dirGIF, "")
	serveDir(dirJS, "./jz/")
	serveDir(dirPNG, "")
	serveDir(dirSVG, "./vz/")
	getParameters("/b/", barcode2D, regexBarcode)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlClubs, clubs)
	getRedirectPermanent(urlLicence, licence)
	gt(urlShooters, shooterSearch, shooters)
	getParameters(urlClub, club, regexID)
	getParameters(urlEntries, entries, regexID)
	getParameters(urlEntryList, entryList, regexID)
	getParameters(urlEventSettings, eventSettings, regexID)
	getParameters(urlEventReport, eventReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexPath)
	getParameters(urlEnterShots, scorecardsIncomplete, regexPath)
	getParameters(urlEnterShotsAll, scorecardsAll, regexPath)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlEnterTotals, totalScoresIncomplete, regexPath)
	getParameters(urlEnterTotalsAll, totalScoresAll, regexPath)
	post(pst, clubNew, clubInsert)
	post(pst, clubDetails, clubDetailsUpsert)
	post(pst, clubMoundNew, clubMoundInsert)
	post(pst, eventNew, eventInsert)
	post(pst, eventDetails, eventDetailsUpsert)
	post(pst, eventRangeNew, eventRangeInsert)
	post(pst, eventAggNew, eventAggInsert)
	post(pst, eventShooterNew, eventShooterInsert)
	post(pst, eventShooterExisting, eventShooterExistingInsert)
	post(get, eventShooterSearch, eventSearchShooters)
	post(pst, shooterNew, shooterInsert)
	post(pst, shooterDetails, shooterUpdate)
	post(pst, eventTotalScores, eventTotalUpsert)
	post(pst, eventAvailableGrades, eventAvailableGradesUpsert)
	post(pst, eventUpdateShotScore, updateShotScores)
	post(pst, importShooter, importShooters)
	post(get, mapResults, mapClubs)
	post(pst, clubMoundEdit, editClubMound)
	post(pst, eventUpdateRange, updateRange)
	post(pst, eventUpdateAgg, updateAgg)
	post(pst, eventEditShooter, eventShooterUpdate)

	//BUG any url breaks when appending "&*((&*%"
	get404(urlEvents, events)
}

func post(method string, formID uint8, runner func(http.ResponseWriter, *http.Request, form, func())) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			//TODO maybe don't redirect user?
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		submittedFields, isValid := validPost(r, getForm(formID))
		redirect := func() { http.Redirect(w, r, r.Referer(), http.StatusSeeOther) }
		newForm := form{
			action: formID,
			Fields: submittedFields,
		}
		if !isValid && method != get {
			setSession(w, newForm)
			redirect()
			return
		}
		runner(w, r, newForm, redirect)
	}
	http.HandleFunc(fmt.Sprintf("/%d", formID), h)
}

//func eventShooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {

func gt(url string, formID uint8, runner func(http.ResponseWriter, *http.Request, form, bool)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			//TODO maybe don't redirect user?
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		submittedFields, isValid := validGet(r, getForm(formID))
		newForm := form{
			action: formID,
			Fields: submittedFields,
		}
		runner(w, r, newForm, isValid)
	})
}
