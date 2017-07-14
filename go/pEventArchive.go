package main

import "net/http"

func eventArchive(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents(onlyClosed)
	templater(w, page{
		Title: "Archive",
		Error: err,
		Data: map[string]interface{}{
			"ListEvents": listEvents,
		},
	})
}

func onlyClosed(event Event) bool {
	return event.Closed
}
