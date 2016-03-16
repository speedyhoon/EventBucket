package main

import (
	"net/http"
	"strconv"
	"strings"
)

func totalScoresAll(w http.ResponseWriter, r *http.Request, parameters string) {
	totalScores(w, r, true, parameters)
}

func totalScoresIncomplete(w http.ResponseWriter, r *http.Request, parameters string) {
	totalScores(w, r, false, parameters)
}

func totalScores(w http.ResponseWriter, r *http.Request, showAll bool, parameters string) {
	//eventID/rangeID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	rangeID, err := strconv.ParseUint(ids[1], 10, 64)
	if rangeID >= uint64(len(event.Ranges)) {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

	templater(w, page{
		Title:   "Total Scores",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"EventID":  event.ID,
			"Shooters": event.Shooters,
			"Range":    event.Ranges[rangeID],
		},
	})
}

func eventTotalUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[2].Value
	rangeID := submittedForm.Fields[3].Value
	shooterID := submittedForm.Fields[4].internalValue.(uint64)

	//Insert new event into database.
	err := upsertScore(eventID, rangeID, shooterID, Score{
		Total:   submittedForm.Fields[0].internalValue.(uint64),
		Centers: submittedForm.Fields[1].internalValue.(uint64),
	})

	//Display any upsert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlTotalScores+eventID /*+"/"_rangeID*/, http.StatusSeeOther)
	//TODO trigger agg calculation immediatly. or maybe inline it within the same DB call?
}
