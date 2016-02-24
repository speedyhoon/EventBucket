package main

import "net/http"

func eventReport(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	templater(w, page{
		Title:   "Event Report",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: M{
			"Event": event,
		},
	})
}
