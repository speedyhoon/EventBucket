package main

func start() {
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
//	Post(URL_champInsert, PostVia(champInsert, URL_organisers))
//	Post(URL_clubMoundInsert, clubMoundInsert)
	Get(URL_home, home)
}
