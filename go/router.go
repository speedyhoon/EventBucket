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

func post(url string, formID uint8, page func(f form) (string, error)) {
	http.HandleFunc(
		url,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				/*405 Method Not Allowed
