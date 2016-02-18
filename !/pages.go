package main

import (
	"fmt"
	"net/http"
)

const (
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
	urlEvent           = "/event/"          //eventID
	urlEventSettings   = "/event-settings/" //eventID
	urlEventReport     = "/event-report/"   //eventID
	urlScoreboard      = "/scoreboard/"     //eventID
	urlScorecards      = "/scorecards/"     //eventID
	urlPrintScorecards = "/print-cards/"    //eventID/shooterID
	urlTotalScores     = "/total-scores/"   //eventID
)

func pages() {
	serveFile(favicon)
	serveFile(robots)
	serveDir(dirCSS, false)
	serveDir(dirJS, false)
	serveDir(dirPNG, false)
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
	getParameters(urlEvent, event, regexID)
	getParameters(urlEventSettings, eventSettings, regexID)
	getParameters(urlEventReport, eventReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexID)
	getParameters(urlScorecards, scorecards, regexID)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlTotalScores, totalScores, regexID)
	post(clubNew, clubInsert)
	post(clubDetails, clubDetailsUpsert)
	post(clubMoundNew, clubMoundInsert)
	post(eventNew, eventInsert)
	post(eventRangeNew, eventRangeInsert)
	post(eventAggNew, eventAggInsert)
	post(eventShooterExisting, eventShooterExistingInsert)
	post(eventShooterNew, eventShooterInsert)

	//BUG any url breaks when appending "&*((&*%"
	get404(urlHome, home)
}

func post(formID uint8, runner func(http.ResponseWriter, *http.Request, form, func())) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		/*for z, input := range GlobalForms[formID] {
			info.Println("pages:", z, input.Options, len(input.Options))
		}*/
		submittedFields, isValid := isValid(r, getForm(formID))
		/*for z, input := range submittedFields {
			info.Println("submittedFields:", z, input.Options, len(input.Options))
		}*/
		redirect := func() { http.Redirect(w, r, r.Referer(), http.StatusSeeOther) }
		newForm := form{
			action: formID,
			Fields: submittedFields,
		}
		//trace.Println("form isValid=", isValid)
		if !isValid {
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
	}, {
		Name: "Clubs",
		Link: urlClubs,
	}, {
		Name: "Shooters",
		Link: urlShooters,
	}, {
		Name: "About",
		Link: urlAbout,
	}, {
		Name: "Licence",
		Link: urlLicence,
	}, {
		Name: "Archive",
		Link: urlArchive,
	}, {
		Name: "Report",
		Link: "/report/",
	}}
	subMenus = map[string][]menu{
		urlEvent: {{
			Name: "Event",
			Link: urlEvent,
		}, {
			Name: "Event Settings",
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
			Name: "Event Report",
			Link: urlEventReport,
		}, {
			Name: "Print Scorecards",
			Link: urlPrintScorecards,
		},
		},
		urlClub: {{
			Name: "Club",
			Link: urlClub,
		}, {
			Name: "Club Settings",
			Link: urlClubSettings,
		}},
	}
)
