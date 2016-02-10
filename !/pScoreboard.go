package main

import "net/http"

func scoreboard(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	templater(w, page{
		Title:  "Scoreboard",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"Event": event,
		},
	})
}