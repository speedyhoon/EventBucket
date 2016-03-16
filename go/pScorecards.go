package main

import (
	"net/http"
	"strings"
)

func scorecardsAll(w http.ResponseWriter, r *http.Request, parameters string) {
	scorecards(w, r, true, parameters)
}

func scorecardsIncomplete(w http.ResponseWriter, r *http.Request, parameters string) {
	scorecards(w, r, false, parameters)
}

func scorecards(w http.ResponseWriter, r *http.Request, showAll bool, parameters string) {
	//eventID/rangeID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var currentRange Range
	currentRange, err = eventRange(event.Ranges, ids[1], w, r)
	if err != nil {
		return
	}

	templater(w, page{
		Title:   "Scorecards",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"EventID":  event.ID,
			"Shooters": event.Shooters,
			"Range":    currentRange,
		},
	})
}
