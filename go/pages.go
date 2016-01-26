package main

import "net/http"

const (
	urlHome     = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlEvents   = "/events"
	urlLicence  = "/licence"
	urlShooters = "/shooters"
	//GET with PARAMETERS
	urlClub          = "/club/"           //clubID
	urlClubSettings  = "/club-settings/"  //clubID
	urlEvent         = "/event/"          //eventID
	urlEventSettings = "/event-settings/" //eventID
	urlEventReport   = "/event-report/"   //eventID
	urlScoreboard    = "/scoreboard/"     //eventID
	urlScorecards    = "/scorecards/"     //eventID
	urlTotalScores   = "/total-scores/"   //eventID
)

func pages() {
	serveFile(favicon)
	serveFile(robots)
	serveDir(dirCSS, false)
	serveDir(dirJS, false)
	serveDir(dirPNG, false)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlClubs, clubs)
	getRedirectPermanent(urlEvents, events)
	getRedirectPermanent(urlLicence, licence)
	getRedirectPermanent(urlShooters, shooters)
	getRedirectPermanent("/all", all)
	getRedirectPermanent("/report", report)
	getParameters(urlClub, club, regexId)
	getParameters(urlClubSettings, clubSettings, regexId)
	getParameters(urlEvent, event, regexId)
	getParameters(urlEventSettings, eventSettings, regexId)
	getParameters(urlEventReport, eventReport, regexId)
	getParameters(urlScoreboard, scoreboard, regexId)
	getParameters(urlScorecards, scorecards, regexId)
	getParameters(urlTotalScores, totalScores, regexId)
	post("/0", insertEvent)
	//	http.HandleFunc("/1", insertEvent)
	post("/2", insertClub)

	//BUG any url breaks when appending "&*((&*%"
	get404(urlHome, home)
}

func post(url string, runner http.HandlerFunc) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			//TODO redirect back to referer address
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		runner(w, r)
	}
	//	http.Handle(url, serveHtml(h))
	http.HandleFunc(url, h)
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
			Name: "Scoreboard",
			Link: urlScoreboard,
		}, {
			Name: "Scorecards",
			Link: urlScorecards,
		}, {
			Name: "Total Scores",
			Link: urlTotalScores,
		}, {
			Name: "Event Settings",
			Link: urlEventSettings,
		}, {
			Name: "Event Report",
			Link: urlEventReport,
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
