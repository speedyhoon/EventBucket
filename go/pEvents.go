package main

import "net/http"

func events(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents()
	templater(w, page{
		Title: "Events",
		Error: err,
		Data: M{
			"NewEvent":   getFormSession(w, r, eventNew),
			"ListEvents": listEvents,
		},
	})
}
