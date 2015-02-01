package main

import (
	"net/http"
	"fmt"
	"os/exec"
)

const (
	VERSION = 0.99
//	PRODUCTION = false //False = output dev warnings, E.g. Template errors
	//TEST_MODE = false //display links to add n shooters or fillout all scores for a given range
	//Known issue - turning off minify breaks the startshooting page. moving to the next sibling in a table row return the textnode of whitespace instead of the next <td> tag
//	MINIFY     = true  //turn on minify html
	//HTML Templates:
	//location "folder path/%v(filename).extension"
	PATH_HTML_MINIFIED = "htm/%v.htm"
	PATH_HTML_SOURCE   = "html/%v.html"
	//Main template html files
	TEMPLATE_HOME  = "_template_home"
	TEMPLATE_ADMIN = "_template_admin"
	TEMPLATE_EMPTY = "_template_empty"
	//folder structure
	DIR_ROOT = "./root/"
	DIR_CSS  = "/c/"
	DIR_JPEG = "/e/"
//	DIR_GIF  = "/g/"
	DIR_JS   = "/j/"
	DIR_PNG  = "/p/"
	DIR_SVG  = "/v/"
//	DIR_WEBP = "/w/"
	FAVICON = "a"

//GET
	URL_home           	= "/"
	URL_about           	= "/about"
	URL_licence         	= "/licence"
	//	URL_licence_summary 	= "/licence-summary"
	URL_archive			  	= "/archive"
	URL_organisers       = "/organisers"
	URL_event            = "/event/"							//event Id special type characters only allowed
	//	URL_events           = "/events/"
//GET with PARAMETERS
	URL_eventSettings    = "/eventSettings/"				//event id
	URL_scoreboard       = "/scoreboard/"					//event id/range_id
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
	URL_eventShotsNSighters  = "/eventShotsNSighters"
	//	URL_rangeReport          = "/rangeReport/"
	URL_randomData           = "/random-data/"
)

func main() {
	if dev_mode_DEBUG{
		agent.Verbose = true
		agent.CollectHTTPStat = true
		agent.NewrelicLicense = "abf730f5454a9a1e78af7a75bfe04565e9e0d3f1"
		agent.Run()
	}
	go DB()
	serveDir(DIR_JS)
	serveDir(DIR_CSS)
	serveDir(DIR_PNG)
	serveDir(DIR_JPEG)
	serveDir(DIR_SVG)

	GetRedirectPermanent(URL_about, about)
	GetRedirectPermanent(URL_archive, archive)
	GetRedirectPermanent(URL_licence, licence)
	GetRedirectPermanent(URL_organisers, organisers)
	GetParameters(URL_event, event)
	GetParameters(URL_eventSettings, eventSettings)
	GetParameters(URL_scoreboard, scoreboard)
//POST
	Post(URL_eventInsert, eventInsert)
/*
	http.HandleFunc(URL_startShooting, html_headers_n_gzip(startShooting))
	http.HandleFunc(URL_startShootingAll, html_headers_n_gzip(startShootingAll))
	http.HandleFunc(URL_totalScores, html_headers_n_gzip(totalScores))
	http.HandleFunc(URL_totalScoresAll, html_headers_n_gzip(totalScoresAll))
//	http.HandleFunc(URL_club, html_headers_n_gzip(club))
//	http.HandleFunc(URL_clubs, html_headers_n_gzip(clubs))
//	http.HandleFunc(URL_club_settings, html_headers_n_gzip(club_settings))
//	http.HandleFunc("/clubs", clubs)
//	http.HandleFunc("/events/", redirectPermanent("/events"))
//	http.HandleFunc("/club/", html_headers_n_gzip(club))
//	http.HandleFunc(URL_rangeReport, html_headers_n_gzip(range_report))
	//Search for a shooter by first, surname & club
	http.HandleFunc(URL_queryShooterList, html_headers_n_gzip(queryShooterList))
//	http.HandleFunc(URL_updateShooterList, html_headers_n_gzip(redirectVia(updateShooterList, URL_organisers)))
	http.HandleFunc(URL_clubInsert, html_headers_n_gzip(redirectVia(clubInsert, URL_organisers)))
	Post(URL_updateRange, rangeUpdate2)
//	http.HandleFunc(URL_dateUpdate, html_headers_n_gzip(dateUpdate))
	http.HandleFunc(URL_eventRangeInsert, html_headers_n_gzip(rangeInsert))
	http.HandleFunc(URL_eventAggInsert, html_headers_n_gzip(aggInsert))
	http.HandleFunc(URL_shooterInsert, html_headers_n_gzip(shooterInsert))
	http.HandleFunc(URL_shooterListInsert, html_headers_n_gzip(shooterListInsert))
	Post(URL_updateTotalScores, updateTotalScores)
	Post(URL_updateShotScores, updateShotScores)
	Post(URL_updateSortScoreBoard, updateSortScoreBoard)
//	http.HandleFunc(URL_updateEventName, html_headers_n_gzip(updateEventName))
	//Add shooters to event
	Post(URL_updateEventGrades, updateEventGrades)
	Post(URL_updateIsPrizeMeet, updateIsPrizeMeet)
	http.HandleFunc(URL_eventShotsNSighters, html_headers_n_gzip(eventShotsNSighters))
	http.HandleFunc(URL_randomData, html_headers_n_gzip(dev_mode_random_data))
	//	http.HandleFunc("/champInsert", redirectVia(champInsert, "/organisers"))
	//	http.HandleFunc("/clubMoundInsert", html_headers_n_gzip(clubMoundInsert))
*/
	Get(URL_home, home)
	url := "http://localhost"
	if exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start() != nil{
		fmt.Println("Unable to open a web browser for "+url)
	}
	fmt.Printf("ListenAndServe: %v", http.ListenAndServe(":80", nil))
}

//TODO remove all custom functions for each set of pages. Just make it overall easier & more flexible to setup new & change existing pages
/*Pages with no data and redirects to the proper page without any parameters
	home
	about
	archive
	licence
	organisers

Pages with event Id parameter
	event
	event settings
	scoreboard

Pages with event Id and Range Id
	start shooting
	start shooting all
	total scores
	total scores all

Post data to Update
	query shooter list
	club insert
	event insert
	update range
	event range insert
	event agg insert
	shoorter insert
	shooter list insert
	update total scores
	update shot scores
	update sort scoreboard
	update event grades
	update is prize meeting
	event shots & sighters

Debug tools
	random data
*/
