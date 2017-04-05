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
	pst         = "POST"
	dirCSS      = "/c/"
	dirJS       = "/j/"
	dirSVG      = "/v/"
	dirWEBP     = "/w/"
	urlEvents   = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlLicence  = "/license"	//TODO should this file be saved with uppercase file name and a .txt file name?
	urlShooters = "/shooters"
	urlSettings = "/settings"
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

func pages() {
	//TODO remove prefix "/" from url & urlLicense?
	serveFile("/favicon.ico")
	serveFile(urlLicence)
	serveDir(dirCSS, true)
	serveDir(dirJS, true)
	serveDir(dirSVG, true)
	serveDir(dirWEBP, false)
	http.Handle("/k/", websocket.Handler(processSocket))
	getParameters("/q/", barcodeQR, regexBarcode)
	getParameters("/x/", barcodeDM, regexBarcode)
	getRedirectPermanent(urlAbout, about)
	getRedirectPermanent(urlSettings, settings)
	getRedirectPermanent(urlArchive, eventArchive)
	getRedirectPermanent(urlClubs, clubs)
	getParameters(urlClub, club, regexID)
	getParameters(urlEntries, entries, regexID)
	getParameters(urlEntryList, entryList, regexID)
	getParameters(urlEventSettings, eventSettings, regexID)
	getParameters(urlEventReport, eventReport, regexID)
	getParameters(urlShooterReport, shooterReport, regexPath)
	getParameters(urlShootersReport, shootersReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexPath)
	getParameters(urlEnterShots, enterShotsIncomplete, regexPath)
	getParameters(urlEnterShotsAll, enterShotsAll, regexPath)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlEnterTotals, enterTotalsIncomplete, regexPath)
	getParameters(urlEnterTotalsAll, enterTotalsAll, regexPath)
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
	//post(pst, eventTotalScores, eventTotalUpsert)
	post(pst, eventAvailableGrades, eventAvailableGradesUpsert)
	//post(pst, eventUpdateShotScore, updateShotScores)
	http.HandleFunc("/17", importShooters) //Don't use normal form validation because a reusable file upload validation function hasn't been written yet.
	post(get, mapResults, mapClubs)
	post(pst, clubMoundEdit, editClubMound)
	post(pst, eventUpdateRange, updateRange)
	post(pst, eventUpdateAgg, updateAgg)
	post(pst, eventEditShooter, eventShooterUpdate)

	gt(urlShooters, shooterSearch, shooters)

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

func gt(url string, formID uint8, runner func(http.ResponseWriter, *http.Request, form, bool)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != get {
			/*405 Method Not Allowed
			A request was made of a resource using a request method not supported by that resource; for example,
			using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
			//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
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
				websocket.Message.Send(ws, eventTotalUpsert(form))
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
				websocket.Message.Send(ws, "!"+updateShotScores(form))
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
