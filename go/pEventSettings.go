package main

import (
	"fmt"
	"net/http"
)

func eventSettings(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	action, forms := sessionForms2(w, r, eventDetails, eventRangeNew, eventAggNew)
	if action == nil || action != nil && *action != eventDetails {
		forms[0].Fields[0].Value = event.Name
		forms[0].Fields[1].Value = event.Club
		forms[0].Fields[2].Value = event.Date
		forms[0].Fields[3].Value = event.Time
		forms[0].Fields[4].Checked = event.Closed
		forms[0].Fields[5].Value = event.ID
	}
	forms[1].Fields[1].Value = eventID

	forms[2].Fields[1].Options = dataListRanges(event.Ranges)
	forms[2].Fields[2].Value = eventID

	templater(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  eventID,
		JS:      "tableSort",
		Heading: event.Name,
		Data: map[string]interface{}{
			"IsClosed":     event.Closed,
			"Ranges":       event.Ranges,
			"EventDetails": forms[0],
			"AddRange":     forms[1],
			"AddAgg":       forms[2],
		},
	})
}

func dataListRanges(ranges []Range) []option {
	var options []option
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, option{Label: r.Name, Value: fmt.Sprintf("%d", r.ID)})
		}
	}
	return options
}

func eventDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[5].Value
	err := updateEventDetails(Event{
		ID:     eventID,
		Name:   submittedForm.Fields[0].Value,
		Club:   submittedForm.Fields[1].Value,
		Date:   submittedForm.Fields[2].Value,
		Time:   submittedForm.Fields[3].Value,
		Closed: submittedForm.Fields[4].internalValue.(bool),
	})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventRangeInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[1].Value
	_, err := eventAddRange(eventID, Range{Name: submittedForm.Fields[0].Value})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventAggInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[2].Value
	rangeID, err := eventAddRange(eventID, Range{
		Name:  submittedForm.Fields[0].Value,
		Aggs:  submittedForm.Fields[1].valueUintSlice,
		IsAgg: true,
	})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
	go upsertAggScores(eventID, rangeID)
}

func eventShooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[6].Value
	shooter := EventShooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[4].valueUint,
		AgeGroup:  submittedForm.Fields[5].valueUint,
	}
	err := eventShooterInsertDB(eventID, shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
func eventShooterExistingInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[3].Value
	shooter, err := getShooter(submittedForm.Fields[0].Value)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	eventShooter := EventShooter{
		Grade:     submittedForm.Fields[1].valueUint,
		AgeGroup:  submittedForm.Fields[2].valueUint,
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
	}
	err = eventShooterInsertDB(eventID, eventShooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
