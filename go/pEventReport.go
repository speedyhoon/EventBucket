package main

import (
	"net/http"
	"strconv"
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
		errorHandler(w, r, "event")
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

func shooterReport(w http.ResponseWriter, r *http.Request, eventID, shooterID string) {
	event, err := getEvent(eventID)
	//If event not found in the database, return error event not found (404).
	if err != nil {
		errorHandler(w, r, "event")
		return
	}

	uShooterID, err := strconv.Atoi(shooterID)
	//If shooter not available in the event, return error shooter not found (404).
	if err != nil || uShooterID >= len(event.Shooters) {
		errorHandler(w, r, "shooter")
		return
	}

	templater(w, page{
		Title:   "Shooter Report",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":   event,
			"Shooter": event.Shooters[uShooterID],
		},
	})
}
