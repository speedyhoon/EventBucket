package main

import "net/http"

func events(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents(onlyOpen)
	_, forms := sessionForms(w, r, eventNew)
	templater(w, page{
		Title: "Events",
		Error: err,
		Data: map[string]interface{}{
			"NewEvent":   forms[0],
			"ListEvents": listEvents,
		},
	})
}

func onlyOpen(event Event) bool {
	return !event.Closed
}
