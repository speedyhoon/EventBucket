package main

import (
	"net/http"
//	"os/exec"
)

const (
//GET
	URL_home           	= "/"
	URL_about           	= "/about"
	URL_club             = "/club"
	URL_licence         	= "/licence"
	//URL_licence_summary	= "/licence-summary"
	URL_archive			  	= "/archive"
	URL_organisers       = "/organisers"
	URL_event            = "/event/"							//event Id special type characters only allowed
//GET with PARAMETERS
	URL_clubs                = "/clubs/"
	//URL_events               = "/events/"
	URL_eventSettings        = "/eventSettings/"			//event id
	URL_scoreboard           = "/scoreboard/"				//event id/range_id
	URL_totalScores          = "/totalScores/"			//event id/range_id
	URL_totalScoresAll       = "/totalScoresAll/"		//event id/range_id
	URL_startShooting        = "/startShooting/"			//event id/range_id
	URL_startShootingAll     = "/startShootingAll/"		//event id/range_id
	URL_queryShooterList 	 = "/queryShooterList"
//POST
	URL_clubInsert           = "/clubInsert"
	//	URL_champInsert          = "/champInsert"
	URL_eventInsert          = "/eventInsert"
	//	URL_eventInsert2         = "/eventInsert2"
	URL_eventRangeInsert     = "/rangeInsert"
	URL_eventAggInsert       = "/aggInsert"
	URL_shooterInsert        = "/shooterInsert"
	URL_shooterListInsert    = "/shooterListInsert"
	URL_updateSortScoreBoard = "/updateSortScoreBoard"
	URL_updateTotalScores    = "/updateTotalScores"
	URL_updateShotScores     = "/updateShotScores"
	URL_updateEventGrades    = "/updateEventGrades"
	URL_updateEventName      = "/updateEventName/"
	URL_updateRange          = "/updateRange"
	URL_updateIsPrizeMeet    = "/updateIsPrizeMeet"
	//	URL_dateUpdate           = "/dateUpdate/"
	//	URL_club                 = "/club/"
	//	URL_clubs                = "/clubs/"
	//	URL_club_settings        = "/clubSettings/"
	URL_club_mound_update    = "/clubMoundUpdate/"
	URL_club_mound_insert    = "/clubMoundInsert/"
	URL_updateShooterList    = "/updateShooterList"
	URL_eventShotsNSighters  = "/eventShotsNSighters/"
	//	URL_rangeReport          = "/rangeReport/"
	URL_randomData           = "/random-data/"
)

func main() {
	dev_mode_NewRelicDebugging()
	go DB()
	serveDir(DIR_JS)
	serveDir(DIR_CSS)
	serveDir(DIR_PNG)
	serveDir(DIR_JPEG)
	serveDir(DIR_SVG)
	//TODO remove these two directories
	serveDir("/logos/")
	serveDir("/html/")
	//TODO remove all custom functions for each set of pages. Just make it overall easier & more flexible to setup new & change existing pages
	GetRedirectPermanent(URL_about, about)
	GetRedirectPermanent(URL_archive, archive)
	GetRedirectPermanent(URL_licence, licence)
	GetRedirectPermanent(URL_organisers, organisers)
	GetRedirectPermanent(URL_clubs, clubs)
//	GetRedirectPermanent(URL_events, events)
	GetParameters(URL_event, event)
	GetParameters(URL_club, club)
	GetParameters(URL_eventSettings, eventSettings)
	GetParameters(URL_scoreboard, scoreboard)
	GetParameters(URL_eventShotsNSighters, eventShotsNSighters)
//	GetParameters(URL_startShooting, startShooting)
//	GetParameters(URL_startShootingAll, startShootingAll)
//	GetParameters(URL_totalScores, totalScores)
//	GetParameters(URL_totalScoresAll, totalScoresAll)
//	GetParameters(URL_club, club)
//	GetParameters(URL_club_settings, club_settings)
//	GetParameters(URL_rangeReport, range_report)
	Post(URL_eventInsert, eventInsert)
	Post(URL_queryShooterList, queryShooterList)	//Search for a shooter by first, surname & club
	Post(URL_updateShooterList, PostVia(updateShooterList, URL_organisers))
	Post(URL_clubInsert, PostVia(clubInsert, URL_organisers))
	Post(URL_updateRange, rangeUpdate2)
//	Post(URL_dateUpdate, dateUpdate)
	Post(URL_eventRangeInsert, rangeInsert)
	Post(URL_eventAggInsert, aggInsert)
	Post(URL_shooterInsert, shooterInsert)
	Post(URL_shooterListInsert, shooterListInsert)
//	Post(URL_updateTotalScores, updateTotalScores)
	Post(URL_updateShotScores, updateShotScores)
	Post(URL_updateSortScoreBoard, updateSortScoreBoard)
//	Post(URL_updateEventName, updateEventName)
	Post(URL_updateEventGrades, updateEventGrades)
	Post(URL_updateIsPrizeMeet, updateIsPrizeMeet)
	Post(URL_randomData, dev_mode_random_data)
//	Post(URL_champInsert, PostVia(champInsert, URL_organisers))
//	Post(URL_clubMoundInsert, clubMoundInsert)
	Get(URL_home, home)
//	url := "http://localhost"+URL_about
//	if exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start() != nil{
//		Warning.Printf("Unable to open a web browser for "+url)
//	}
	Info.Println("ready to go")
	Warning.Println("ListenAndServe: %v", http.ListenAndServe(":80", nil))
}
