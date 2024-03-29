package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/speedyhoon/frm"
)

func enterTotalsAll(w http.ResponseWriter, r *http.Request, event Event, rangeID rID) {
	enterTotals(w, r, true, event, rangeID)
}

func enterTotalsIncomplete(w http.ResponseWriter, r *http.Request, event Event, rangeID rID) {
	enterTotals(w, r, false, event, rangeID)
}

func enterTotals(w http.ResponseWriter, r *http.Request, showAll bool, event Event, rangeID rID) {
	currentRange, err := eventRange(event.Ranges, rangeID, w, r)
	if err != nil {
		return
	}

	//_, forms := sessionForms(w, r, eventTotalScores)
	//t.Printf("%+v\n", forms[0])

	var hidden int
	if !showAll && !event.Closed {
		rangeID := currentRange.StrID()
		for _, shooter := range event.Shooters {
			if shooter.Disabled || shooter.Scores[rangeID].Total >= 1 {
				hidden++
			}
		}
	}

	render(w, page{
		Title:   "Enter Totals",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: currentRange.Name,
		Error:   err,
		Data: map[string]interface{}{
			"Range":       currentRange,
			"Event":       event,
			"URL":         "enter-totals",
			"ShowAll":     showAll,
			"Hidden":      hidden,
			"Plural":      plural(hidden, " is", "s have"),
			"Disciplines": globalDisciplines,
		},
	})
}

func eventTotalUpsert(fields []frm.Field) string {
	// Save score to event in database.
	err := updateDocument(tblEvent, fields[2].Str(), &shooterScore{
		rangeID: fields[3].Str(),
		id:      fields[4].Uint(),
		score: Score{
			Total:    fields[0].Uint(),
			Centers:  fields[1].Uint(),
			ShootOff: fields[5].Uint(),
		}}, &Event{}, upsertScore)

	// Return any upsert errors onscreen.
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Saved %d.%d for shooter %d on range %v in event %v.", fields[0].Value, fields[1].Value, fields[4].Value, fields[3].Value, fields[2].Value)
}

func eventRange(ranges []Range, rangeID rID, w http.ResponseWriter, r *http.Request) (Range, error) {
	for _, r := range ranges {
		if r.ID == uint(rangeID) {
			// If range is an aggregate return an error message.
			if r.IsAgg {
				return Range{}, errors.New("range is an aggregate and scores can't be entered directly")
			}
			// Return valid range.
			return r, nil
		}
	}
	// Otherwise event doesn't contain a range with that id and return 404.
	errorHandler(w, r, "range")
	return Range{}, errors.New("a range with that ID doesn't exists in this event")
}
