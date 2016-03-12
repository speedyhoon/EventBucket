package main

import (
	"fmt"
	"net/http"
	"regexp"
)

const (
	get         = "GET"
	pst         = "POST"
	urlHome     = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlEvents   = "/events"
	urlLicence  = "/licence"
	urlShooters = "/shooters"
	//GET with PARAMETERS
	urlClub            = "/club/"           //clubID
	urlClubSettings    = "/club-settings/"  //clubID
	urlEntries         = "/entries/"        //eventID
	urlEventSettings   = "/event-settings/" //eventID
	urlEventReport     = "/event-report/"   //eventID
	urlScoreboard      = "/scoreboard/"     //eventID
	urlScorecards      = "/scorecards/"     //eventID
	urlPrintScorecards = "/print-cards/"    //eventID/shooterID
	urlTotalScores     = "/total-scores/"   //eventID
)

var (
	//URL validation matching
	regexID      = regexp.MustCompile(`^[a-z0-9]+$`)
	regexPath    = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+$`)
	regexBarcode = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+/[a-z0-9]+$`)
)

func pages() {
	serveFile(favicon)
	serveFile(robots)
	serveDir(dirCSS, "./cz/")
	serveDir(dirJS, "./jz/")
	serveDir(dirPNG, "")
	serveDir(dirGIF, "")
	getParameters("/b/", base64QrH, regexBarcode)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlClubs, clubs)
	getRedirectPermanent(urlEvents, events)
	getRedirectPermanent(urlLicence, licence)
	getRedirectPermanent(urlShooters, shooters)
	getRedirectPermanent("/all", all)
	getRedirectPermanent("/report", report)
	getParameters(urlClub, club, regexID)
	getParameters(urlClubSettings, clubSettings, regexID)
	getParameters(urlEntries, entries, regexID)
	getParameters(urlEventSettings, eventSettings, regexID)
	getParameters(urlEventReport, eventReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexID)
	getParameters(urlScorecards, scorecards, regexID)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlTotalScores, totalScores, regexID)
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
	post(pst, shooterSearch, searchShooters)

	//BUG any url breaks when appending "&*((&*%"
	get404(urlHome, home)
}

func post(method string, formID uint8, runner func(http.ResponseWriter, *http.Request, form, func())) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		submittedFields, isValid := isValid(r, getForm(formID))
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

var (
	mainMenu = []menu{{
		Name: "Home",
		Link: urlHome,
	}, {
		Name: "Events",
		Link: urlEvents,
		SubMenu: []menu{{
			Name: "Entries",
			Link: urlEntries,
		}, {
			Name: "Settings",
			Link: urlEventSettings,
		}, {
			Name: "Scoreboard",
			Link: urlScoreboard,
		}, {
			Name: "Scorecards",
			Link: urlScorecards,
		}, {
			Name: "Total Scores",
			Link: urlTotalScores,
		}, {
			Name: "Report",
			Link: urlEventReport,
		}},
	}, {
		Name: "Clubs",
		Link: urlClubs,
		SubMenu: []menu{{
			Name: "Club",
			Link: urlClub,
		}, {
			Name: "Settings",
			Link: urlClubSettings,
		}},
	}, {
		Name: "Shooters",
		Link: urlShooters,
	}, {
		Name: "Archive",
		Link: urlArchive,
	}, {
		Name: "About",
		Link: urlAbout,
	}, {
		Name: "Licence",
		Link: urlLicence,
	}}
)
