package main

import "net/http"

func events(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents(onlyOpen)
	templater(w, page{
		Title: "Events",
		Error: err,
		Data: map[string]interface{}{
			"NewEvent":   getFormSession(w, r, eventNew),
			"ListEvents": listEvents,
		},
	})
}
