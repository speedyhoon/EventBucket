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
	dirCSS      = "/c/"
	dirJS       = "/j/"
	dirWEBP     = "/w/"
	urlEvents   = "/"
	urlAbout    = "/about"
	urlArchive  = "/archive"
	urlClubs    = "/clubs"
	urlLicence  = "/license"
	urlShooters = "/shooters"
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
	getParameters("/q/", barcodeQR, regexBarcode)
	getParameters("/x/", barcodeDM, regexBarcode)
	getRedirectPermanent(urlAbout, about)
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

	//TODO BUG any url breaks when appending "&*((&*%"
	get404(events)
}

func post(url string, formID uint8, page func(f form) (string, error)) {
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				/*405 Method Not Allowed
				A request was made of a resource using a request method not supported by that resource; for example,
				using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
				//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
				http.Redirect(w, r, r.Referer(), http.StatusMethodNotAllowed)
				return
			}
			f, ok := validBoth(r, formID)
			if !ok {
				setSession(w, f)
				http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
				return
			}
			redirect, err := page(f)
			//Display any insert errors onscreen.
			if err != nil {
				formError(w, r, f, err)
				return
			}
			if redirect == "" {
				redirect = r.Referer()
			}
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		},
	)
}

func get(url string, formID uint8, page func(http.ResponseWriter, *http.Request, form)) {
	http.HandleFunc(
		url,
		isGetMethod(func(w http.ResponseWriter, r *http.Request) {
			f, _ := validBoth(r, formID)
			page(w, r, f)
		}))
}

//Start listening to each websocket client that connects.
func processSocket(ws *websocket.Conn) {
	var msg string
	var formID uint8
	var err error
	send := func(str string) {
		err := websocket.Message.Send(ws, str)
		if err != nil {
			warn.Println(err)
		}
	}
	//Start a loop to listen to incoming websocket traffic from all clients.
	for {
		//Ignore any empty messages.
		if websocket.Message.Receive(ws, &msg) != nil || len(msg) < 1 {
			break
		}
		//The first character of the websocket message is used as a "router" to decide where to send the message.
		formID = msg[0]
		//Ignore any messages that do not have a case in this switch.
		switch formID {
		case eventTotalScores:
			var form url.Values
			err = json.Unmarshal([]byte(msg[1:]), &form)
			if err != nil {
				warn.Println(err)
				continue
			}
			if form, passed := isValid(form, formID); passed {
				send(eventTotalUpsert(form.Fields))
			} else {
				send(fmt.Sprintf("Unable to save %v.", msg))
			}
		case eventUpdateShotScore:
			var form url.Values
			err = json.Unmarshal([]byte(msg[1:]), &form)
			if err != nil {
				warn.Println(err)
				continue
			}

			if form, passed := isValid(form, formID); passed {
				send("!" + updateShotScores(form.Fields))
			} else {
				var response []byte
				response, err = json.Marshal(form)
				if err != nil {
					warn.Println(err)
					continue
				}
				send(fmt.Sprintf("!%U%s", msg[0], response))
			}
		case 126: //getDisciplines:
			var response []byte
			response, err = json.Marshal(globalDisciplines)
			if err != nil {
				warn.Println(err)
				continue
			}
			send(fmt.Sprintf("~%s", response))
		}
	}
}
