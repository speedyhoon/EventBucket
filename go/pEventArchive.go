package main

import "net/http"

func eventArchive(w http.ResponseWriter, r *http.Request) {
	events, err := getEvents(onlyClosed)
	render(w, page{
		Title: "Archive",
		Error: err,
		Data: map[string]interface{}{
			"Events": events,
		},
	})
}

func onlyClosed(event Event) bool {
	return event.Closed
}
