package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/v8"
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
	urlEventSettings   = "/event-settings/"   //eventID
	urlEventReport     = "/event-report/"     //eventID
	urlScoreboard      = "/scoreboard/"       //eventID
	urlEnterShots      = "/enter-shots/"      //eventID
	urlEnterShotsAll   = "/enter-shots-all/"  //eventID
	urlPrintScorecards = "/print-scorecards/" //eventID/shooterID
	urlEnterTotals     = "/enter-totals/"     //eventID
	urlEnterTotalsAll  = "/enter-totals-all/" //eventID
)

var (
	mainMenu = []menu{{
		Name: "Events",
		Link: urlEvents,
		SubMenu: []menu{{
			Name: "Entries",
			Link: urlEntries,
		}, {
			Name: "Event Settings",
			Link: urlEventSettings,
		}, {
			Name:      "Scoreboard",
			Link:      urlScoreboard,
			RangeMenu: true,
		}, {
			Name:      "Enter Shots",
			Link:      urlEnterShots,
			RangeMenu: true,
		}, {
			Name:      "Enter Totals",
			Link:      urlEnterTotals,
			RangeMenu: true,
		}, {
			Name: "Event Report",
			Link: urlEventReport,
		}},
	}, {
		Name: "Clubs",
		Link: urlClubs,
		SubMenu: []menu{{
			Name: "Club",
			Link: urlClub,
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
	getParameters(urlEventSettings, eventSettings, regexID)
	getParameters(urlEventReport, eventReport, regexID)
	getParameters(urlScoreboard, scoreboard, regexPath)
	getParameters(urlEnterShots, enterShotsIncomplete, regexPath)
	getParameters(urlEnterShotsAll, enterShotsAll, regexPath)
	getParameters(urlPrintScorecards, printScorecards, regexPath)
	getParameters(urlEnterTotals, enterTotalsIncomplete, regexPath)
	getParameters(urlEnterTotalsAll, enterTotalsAll, regexPath)

	//TODO BUG any url breaks when appending "&*((&*%"
	get404(events)
}

func post(url string, formID uint8, p func(f frm.Form) (string, error)) {
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				/*405 Method Not Allowed
				A request was made of a resource using a request method not supported by that resource; for example,
				using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
				//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
				//http.Redirect(w, r, r.Referer(), http.StatusMethodNotAllowed)

				render(w, page{
					Title:  "Error",
					Status: http.StatusMethodNotAllowed,
					Data: map[string]interface{}{
						"Type": "incorrect form submission",
					},
				})

				return
			}
			f, ok := v8.IsValidRequest(r, getFields(formID))
			form := frm.Form{Action: formID, Fields: f}
			if !ok {
				redirectError(w, r, form)
				return
			}

			var redirect string
			redirect, form.Err = p(form)
			//Display any insert errors onscreen.
			if form.Err != nil {
				form.Fields[0].Focus = true
				redirectError(w, r, form)
				return
			}

			//Display default acknowledgement message using a session
			session.Set(w, frm.Form{Action: formID}) //Don't store validated fields in session

			if redirect == "" {
				redirect = r.Referer()
			}
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		},
	)
}

func redirectError(w http.ResponseWriter, r *http.Request, f frm.Form) {
	session.Set(w, f)
	http.Redirect(w, r, fmt.Sprintf("%v#%v", r.Referer(), f.Action), http.StatusSeeOther)
}

func get(url string, formID uint8, page func(http.ResponseWriter, *http.Request, []frm.Field)) {
	http.HandleFunc(
		url,
		isGetMethod(func(w http.ResponseWriter, r *http.Request) {
			f, _ := v8.IsValidRequest(r, getFields(formID))
			page(w, r, f)
		}))
}

//Start listening to each websocket client that connects.
func processSocket(ws *websocket.Conn) {
	var msg string
	var formID uint8
	var err error
	send := func(str string) {
		if err = websocket.Message.Send(ws, str); err != nil {
			log.Println(err)
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
			var urlValues url.Values
			err = json.Unmarshal([]byte(msg[1:]), &urlValues)
			if err != nil {
				log.Println(err)
				continue
			}
			if fields, passed := v8.IsValid(urlValues, getFields(formID)); passed {
				send(eventTotalUpsert(fields))
			} else {
				send(fmt.Sprintf("Unable to save %v.", msg))
			}
		case eventUpdateShotScore:
			var urlValues url.Values
			err = json.Unmarshal([]byte(msg[1:]), &urlValues)
			if err != nil {
				log.Println(err)
				continue
			}

			if fields, passed := v8.IsValid(urlValues, getFields(formID)); passed {
				send("!" + updateShotScores(fields))
			} else {
				var response []byte
				response, err = json.Marshal(fields)
				if err != nil {
					log.Println(err)
					continue
				}
				send(fmt.Sprintf("!%U%s", msg[0], response))
			}
		case 126: //getDisciplines:
			var response []byte
			response, err = json.Marshal(globalDisciplines)
			if err != nil {
				log.Println(err)
				continue
			}
			send(fmt.Sprintf("~%s", response))
		}
	}
}
