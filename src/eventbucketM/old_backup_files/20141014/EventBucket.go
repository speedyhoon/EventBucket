package main

import (
	"net/http"
	"fmt"
	"os/exec"
)

func main() {
	DB_connection()
	file_server := http.FileServer(http.Dir(DIR_ROOT))
	//TODO make sure ALL resources don't have a . extension to save network bandwidth
	http.Handle(DIR_JS, file_headers_n_gzip(file_server, "js"))
	http.Handle(DIR_CSS, file_headers_n_gzip(file_server, "css"))
	http.Handle(DIR_PNG, file_headers_n_gzip(file_server, "png"))
	http.Handle(DIR_JPEG, file_headers_n_gzip(file_server, "jpg"))
	http.Handle(DIR_ICON, file_headers_n_gzip(file_server, "png"))
	http.Handle(DIR_SVG, file_headers_n_gzip(file_server, "svg"))

	//GET
	http.HandleFunc("/", html_headers_n_gzip(home))
	http.HandleFunc(URL_about, html_headers_n_gzip(about))
	http.HandleFunc(URL_about+"/", html_headers_n_gzip(redirectPermanent(URL_about)))
	http.HandleFunc(URL_licence, html_headers_n_gzip(licence))
	http.HandleFunc(URL_licence+"/", html_headers_n_gzip(redirectPermanent(URL_licence)))
	http.HandleFunc(URL_licence_summary, html_headers_n_gzip(licence_summary))
	http.HandleFunc(URL_licence_summary+"/", html_headers_n_gzip(redirectPermanent(URL_licence_summary)))

	http.HandleFunc(URL_organisers, html_headers_n_gzip(organisers))
	http.HandleFunc(URL_organisers+"/", html_headers_n_gzip(redirectPermanent(URL_organisers)))

	http.HandleFunc(URL_event, html_headers_n_gzip(event))
	http.HandleFunc(URL_eventSettings, html_headers_n_gzip(eventSettings))

	http.HandleFunc(URL_startShooting, html_headers_n_gzip(startShooting))
	http.HandleFunc(URL_startShootingAll, html_headers_n_gzip(startShootingAll))

	http.HandleFunc(URL_totalScores, html_headers_n_gzip(totalScores))
	http.HandleFunc(URL_totalScoresAll, html_headers_n_gzip(totalScoresAll))

	http.HandleFunc(URL_scoreboard, html_headers_n_gzip(scoreboard))
	http.HandleFunc(URL_club, html_headers_n_gzip(club))
//	http.HandleFunc(URL_clubs, html_headers_n_gzip(clubs))
	http.HandleFunc(URL_club_settings, html_headers_n_gzip(club_settings))
	//	//	http.HandleFunc("/clubs", clubs)
	//	http.HandleFunc("/events/", redirectPermanent("/events"))
	//	http.HandleFunc("/club/", html_headers_n_gzip(club))

	//Search for a shooter by first, surname & club
	http.HandleFunc(URL_queryShooterList, html_headers_n_gzip(queryShooterList))

	//POST
	http.HandleFunc(URL_clubInsert, html_headers_n_gzip(redirectVia(clubInsert, URL_organisers)))
	http.HandleFunc(URL_eventInsert, html_headers_n_gzip(redirectVia(eventInsert, URL_organisers)))

	http.HandleFunc(URL_eventInsert2, html_headers_n_gzip(eventInsert2))

	http.HandleFunc(URL_updateRange, html_headers_n_gzip(rangeUpdate2))
	http.HandleFunc(URL_dateUpdate, html_headers_n_gzip(dateUpdate))
	http.HandleFunc(URL_eventRangeInsert, html_headers_n_gzip(rangeInsert))
	http.HandleFunc(URL_eventAggInsert, html_headers_n_gzip(aggInsert))

	http.HandleFunc(URL_shooterInsert, html_headers_n_gzip(shooterInsert))
	http.HandleFunc(URL_shooterListInsert, html_headers_n_gzip(shooterListInsert))

	http.HandleFunc(URL_updateTotalScores, html_headers_n_gzip(updateTotalScores))
	http.HandleFunc(URL_updateShotScores, html_headers_n_gzip(updateShotScores))
	http.HandleFunc(URL_updateSortScoreBoard, html_headers_n_gzip(updateSortScoreBoard))

	http.HandleFunc(URL_updateEventName, html_headers_n_gzip(updateEventName))

	//Add shooters to event
	http.HandleFunc(URL_updateEventGrades, html_headers_n_gzip(updateEventGrades))
	http.HandleFunc(URL_updateShooterList, html_headers_n_gzip(redirectVia(updateShooterList, URL_organisers)))

	http.HandleFunc(URL_updateIsPrizeMeet, html_headers_n_gzip(updateIsPrizeMeet))

	http.HandleFunc(URL_eventShotsNSighters, html_headers_n_gzip(eventShotsNSighters))

	//	http.HandleFunc("/champInsert", redirectVia(champInsert, "/organisers"))
	//	http.HandleFunc("/clubMoundInsert", html_headers_n_gzip(clubMoundInsert))



	http.HandleFunc("/random-data/", html_headers_n_gzip(random_data))
	err := http.ListenAndServe(":80", nil)
	if err != nil{
//		log.Fatal("ListenAndServe: ", err)
		fmt.Printf("ListenAndServe: %v", err)
	}
	if PRODUCTION {
		url := "http://localhost/"
		err = exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
		if err != nil{
			fmt.Printf("Unable to open a web browser for %v", url)
		}
	}
}
