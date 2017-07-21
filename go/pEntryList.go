package main

import "net/http"

func entryList(w http.ResponseWriter, r *http.Request, event Event) {
	templater(w, page{
		Title:   "Entry List",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event": event,
		},
	})
}
