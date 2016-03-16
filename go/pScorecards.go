package main

import (
	"net/http"
	"strconv"
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

	rangeID, err := strconv.ParseUint(ids[1], 10, 64)
	var hasRange bool
	var currentRange Range
	for i, r := range event.Ranges {
		if r.ID == rangeID {
			currentRange = event.Ranges[i]
			continue
		}
	}
	if !hasRange {
		errorHandler(w, r, http.StatusNotFound, "range")
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
