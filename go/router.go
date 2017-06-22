package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"golang.org/x/net/websocket"
)

const (
	get         = "GET"
	post        = "POST"
	dirCSS      = "/c/"
	dirJS       = "/j/"
	dirWEBP     = "/w/"
	urlEvents   = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlLicence  = "/license"
	urlShooters = "/shooters"
	urlSettings = "/settings"
	urlSVG      = "/v"
	//GET with PARAMETERS
	urlClub            = "/club/"             //clubID
	urlEntries         = "/entries/"          //eventID
	urlEntryList       = "/print-entry-list/" //eventID
	urlEventSettings   = "/event-settings/"   //eventID
	urlEventReport     = "/event-report/"     //eventID
	urlShooterReport   = "/shooter-report/"   //eventID/shooterID
	urlShootersReport  = "/shooters-report/"  //eventID
	urlScoreboard      = "/scoreboard/"       //eventID
	urlEnterShots      = "/enter-shots/"      //eventID
	urlEnterShotsAll   = "/enter-shots-all/"  //eventID
	urlPrintScorecards = "/print-scorecards/" //eventID/shooterID
	urlEnterTotals     = "/enter-totals/"     //eventID
	urlEnterTotalsAll  = "/enter-totals-all/" //eventID
)

var (
	//URL validation matching
	regexID      = regexp.MustCompile(`^[a-z0-9]+$`)
	regexPath    = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+$`)
	regexBarcode = regexp.MustCompile(`^[a-z0-9]+/[a-z0-9]+#[a-z0-9]+$`)
)

func init() {
	serveFile("/favicon.ico", false)
	serveFile(urlLicence, false)
	serveFile(urlSVG, true)
	serveDir(dirCSS, true)
	serveDir(dirJS, true)
	serveDir(dirWEBP, false)
	http.Handle("/k/", websocket.Handler(processSocket))
	getParameter("/q/", barcodeQR, regexBarcode)
	getParameter("/x/", barcodeDM, regexBarcode)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlSettings, settings)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlClubs, clubs)
	getParameter(urlClub, club, regexID)
	getParameter(urlEntries, entries, regexID)
	getParameter(urlEntryList, entryList, regexID)
	getParameter(urlEventSettings, eventSettings, regexID)
	getParameter(urlEventReport, eventReport, regexID)
	getParameters(urlShooterReport, shooterReport, regexPath)
	getParameter(urlShootersReport, shootersReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexPath)
	getParameters(urlEnterShots, enterShotsIncomplete, regexPath)
	getParameters(urlEnterShotsAll, enterShotsAll, regexPath)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlEnterTotals, enterTotalsIncomplete, regexPath)
	getParameters(urlEnterTotalsAll, enterTotalsAll, regexPath)
	endpoint(post, clubNew, clubInsert)
	endpoint(post, clubDetails, clubDetailsUpsert)
	endpoint(post, clubMoundNew, clubMoundInsert)
	endpoint(post, eventNew, eventInsert)
	endpoint(post, eventDetails, eventDetailsUpsert)
	endpoint(post, eventRangeNew, eventRangeInsert)
	endpoint(post, eventAggNew, eventAggInsert)
	endpoint(post, eventShooterNew, eventShooterInsert)
	endpoint(post, eventShooterExisting, eventShooterExistingInsert)
	endpoint(get, eventShooterSearch, eventSearchShooters)
	endpoint(post, shooterNew, shooterInsert)
	endpoint(post, shooterDetails, shooterUpdate)
	//TODO re-enable
	//endpoint(post, eventTotalScores, eventTotalUpsert)
	endpoint(post, eventAvailableGrades, eventAvailableGradesUpsert)
	//TODO re-enable
	//endpoint(post, eventUpdateShotScore, updateShotScores)
	http.HandleFunc("/16", importShooters) //TODO file upload validation function hasn't been written yet.
	endpoint(get, mapResults, mapClubs)
	endpoint(post, clubMoundEdit, editClubMound)
	endpoint(post, eventUpdateRange, updateRange)
	endpoint(post, eventUpdateAgg, updateAgg)
	endpoint(post, eventEditShooter, eventShooterUpdate)

	gt(urlShooters, shooterSearch, shooters)

	//TODO BUG any url breaks when appending "&*((&*%"
	get404(urlEvents, events)
}

func endpoint(method string, formID uint8, runner func(http.ResponseWriter, *http.Request, form)) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			http.Redirect(w, r, r.Referer(), http.StatusMethodNotAllowed)
			return
		}
		form, ok := validPost(r, formID)
		if !ok && method != get {
			setSession(w, form)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
		runner(w, r, form)
	}
	http.HandleFunc(fmt.Sprintf("/%d", formID), h)
}

func gt(url string, formID uint8, runner func(http.ResponseWriter, *http.Request, form)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != get {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
			http.Redirect(w, r, r.Referer(), http.StatusMethodNotAllowed)
			return
		}
		form, _ := validGet(r, formID)
		runner(w, r, form)
	})
}

//Start listening to each websocket client that connects.
func processSocket(ws *websocket.Conn) {
	var msg string
	var command uint8
	var err error
	//Start a loop to listen to incoming websocket traffic from all clients.
	for {
		//Ignore any empty messages.
		if websocket.Message.Receive(ws, &msg) != nil || len(msg) < 1 {
			break
		}
		//The first character of the websocket message is used as a "router" to decide where to send the message.
		command = uint8(msg[0])
		//Ignore any messages that do not have a case in this switch.
		switch command {
		case eventTotalScores:
			var form url.Values
			err = json.Unmarshal([]byte(msg[1:]), &form)
			if err != nil {
				warn.Println(err)
				continue
			}
			if form, passed := isValid(form, getForm(command)); passed {
				websocket.Message.Send(ws, eventTotalUpsert(form.Fields))
			} else {
				websocket.Message.Send(ws, fmt.Sprintf("Unable to save %v.", msg))
			}
		case eventUpdateShotScore:
			var form url.Values
			err = json.Unmarshal([]byte(msg[1:]), &form)
			if err != nil {
				warn.Println(err)
				continue
			}

			if form, passed := isValid(form, getForm(command)); passed {
				websocket.Message.Send(ws, "!"+updateShotScores(form.Fields))
			} else {
				var response []byte
				response, err = json.Marshal(form)
				if err != nil {
					warn.Println(err)
					continue
				}
				websocket.Message.Send(ws, fmt.Sprintf("!%U%s", msg[0], response))
			}
		case 126: //getDisciplines:
			var response []byte
			response, err = json.Marshal(globalDisciplines)
			if err != nil {
				warn.Println(err)
				continue
			}
			websocket.Message.Send(ws, fmt.Sprintf("~%s", response))
		}
	}
}
