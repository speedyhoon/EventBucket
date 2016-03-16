package main

import (
	"errors"
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

	var currentRange Range
	currentRange, err = eventRange(event.Ranges, ids[1], w, r)
	if err != nil {
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
			"Range":    currentRange,
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

func eventRange(ranges []Range, rID string, w http.ResponseWriter, r *http.Request) (Range, error) {
	var currentRange Range
	rangeID, err := strconv.ParseUint(rID, 10, 64)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "range")
		return currentRange, err
	}
	var hasRange bool
	for _, r := range ranges {
		if r.ID == rangeID {
			currentRange = r
			hasRange = true
			continue
		}
	}
	if !hasRange {
		errorHandler(w, r, http.StatusNotFound, "range")
		return currentRange, errors.New("Range with that ID doesn't exists in this event")
	}
	return currentRange, nil
}
