package main

import "net/http"

func eventReport(w http.ResponseWriter, r *http.Request, event Event) {
	templater(w, page{
		Title:   "Event Report",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event": event,
		},
	})
}
