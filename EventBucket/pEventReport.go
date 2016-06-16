package main

import (
	"net/http"
	"strconv"
	"strings"
)

func eventReport(w http.ResponseWriter, r *http.Request, eventID string) {
	templater(w, eventReportPage(w, r, eventID, "Event Report"))
}

func shootersReport(w http.ResponseWriter, r *http.Request, eventID string) {
	templater(w, eventReportPage(w, r, eventID, "Shooters Report"))
}

func eventReportPage(w http.ResponseWriter, r *http.Request, eventID, title string) page {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return page{}
	}

	return page{
		Title:   title,
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event": event,
		},
	}
}

func shooterReport(w http.ResponseWriter, r *http.Request, parameters string) {
	ids := strings.Split(parameters, "/")
	eventID := ids[0]

	event, err := getEvent(eventID)
	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	shooterID, err := strconv.Atoi(ids[1])
	if err != nil || shooterID >= len(event.Shooters) {
		errorHandler(w, r, http.StatusNotFound, "shooter")
		return
	}

	templater(w, page{
		Title:   "Shooter Report",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":   event,
			"Shooter": event.Shooters[shooterID],
		},
	})
}
