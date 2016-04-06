package main

import (
	"errors"
	"net/http"
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

	templater(w, page{
		Title:   "Total Scores",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		JS:      "totalScores",
		Error:   err,
		Data: map[string]interface{}{
			"EventID":  event.ID,
			"Shooters": event.Shooters,
			"Range":    currentRange,
			"ShowAll":  showAll,
		},
	})
}

func eventTotalUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[2].Value
	rangeID := submittedForm.Fields[3].Value
	shooterID := submittedForm.Fields[4].valueUint

	//Insert new event into database.
	err := upsertScore(eventID, rangeID, shooterID, Score{
		Total:   submittedForm.Fields[0].valueUint,
		Centers: submittedForm.Fields[1].valueUint,
	})

	//Display any upsert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func eventRange(ranges []Range, rID string, w http.ResponseWriter, r *http.Request) (Range, error) {
	//If range id is not a number, return 404.
	rangeID, err := strToUint(rID)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "range")
		return Range{}, err
	}

	for _, r := range ranges {
		if r.ID == rangeID {
			//If range is an aggregate return an error message.
			if r.IsAgg {
				return Range{}, errors.New("Range is an aggregate and scores can't be entered directly.")
			}
			//Return valid range.
			return r, nil
		}
	}
	//Otherwise event doesn't contain a range with that id and return 404.
	errorHandler(w, r, http.StatusNotFound, "range")
	return Range{}, errors.New("Range with that ID doesn't exists in this event")
}
