package main

import (
	"net/http"
	"mgo"
)

var conn *mgo.Database

func main() {
	conn = DB()
	//GET
	http.HandleFunc("/", home)

	http.HandleFunc("/organisers", organisers)
	http.HandleFunc("/organisers/", redirectPermanent("/organisers"))
	//	http.HandleFunc("/", home)
	//	http.HandleFunc("/clubs", clubs)
	//	http.HandleFunc("/startShooting", startShooting)
	//	http.HandleFunc("/organiser", organisers)
	http.HandleFunc("/events/", redirectPermanent("/events"))
	http.HandleFunc("/events", events)


	http.HandleFunc("/event/", event)

	//	http.HandleFunc("/eventSetup", eventSetup)
	http.HandleFunc("/eventSettings/", eventSettings)
	http.HandleFunc("/try", tempTry)




	//POST
	http.HandleFunc("/clubInsert", redirectTo(clubInsert, "/organisers"))
	http.HandleFunc("/eventInsert", redirectTo(eventInsert, "/organisers"))
	http.HandleFunc("/champInsert", redirectTo(champInsert, "/organisers"))
	http.HandleFunc("/rangeInsert", redirectTo(rangeInsert, "/eventSettings/"))


	http.ListenAndServe(":80", nil)
}
