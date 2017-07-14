package main

import (
	"fmt"
	"net/http"
)

func eventSettings(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If the event isn't found in the database, return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var club Club
	if !event.Closed && event.Club != "" {
		club, err = getClubByName(event.Club)
	}

	//Retrieve any submitted form that failed to save.
	action, forms := sessionForms(w, r, eventDetails, eventRangeNew, eventAggNew, eventRangeUpdate, eventAggUpdate)
	if action != eventDetails {
		forms[0].Fields[0].Value = event.Club
		forms[0].Fields[1].Value = event.Name
		forms[0].Fields[2].Value = event.Date
		forms[0].Fields[3].Value = event.Time
		forms[0].Fields[4].Checked = event.Closed
		forms[0].Fields[5].Value = event.ID
	}
	forms[1].Fields[1].Value = eventID

	forms[2].Fields[1].Options = dataListRanges(event.Ranges, true)
	forms[2].Fields[2].Value = eventID

	templater(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":        dataListRanges(event.Ranges, false),
			"Event":         event,
			"EventDetails":  forms[0],
			"AddRange":      forms[1],
			"AddAgg":        forms[2],
			"RangeDataList": club.Mounds,
			"eventRangeUpdate": forms[3],
			"eventAggUpdate": forms[4],
		},
	})
}

func dataListRanges(ranges []Range, selected bool) (options []option) {
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, option{Label: r.Name, Value: fmt.Sprintf("%d", r.ID), Selected: selected})
		}
	}
	return options
}

func eventDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	eventID := submittedForm.Fields[5].Value
	err := updateDocument(tblEvent, eventID, &Event{
		Club:   submittedForm.Fields[0].Value,
		Name:   submittedForm.Fields[1].Value,
		Date:   submittedForm.Fields[2].Value,
		Time:   submittedForm.Fields[3].Value,
		Closed: submittedForm.Fields[4].Checked,
	}, &Event{}, updateEventDetails)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventRangeInsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	eventID := submittedForm.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &Range{Name: submittedForm.Fields[0].Value}, &Event{}, eventAddRange)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func updateRange(w http.ResponseWriter, r *http.Request, submittedForm form) {
	eventID := submittedForm.Fields[0].Value
	err := updateDocument(tblEvent, eventID, &Range{
		ID:     submittedForm.Fields[1].valueUint,
		Name:   submittedForm.Fields[2].Value,
		Locked: submittedForm.Fields[3].Checked,
		Order:  submittedForm.Fields[4].valueUint,
	}, &Event{}, editRange)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func updateAgg(w http.ResponseWriter, r *http.Request, submittedForm form) {
	eventID := submittedForm.Fields[0].Value
	err := updateDocument(tblEvent, eventID, &Range{
		ID:    submittedForm.Fields[1].valueUint,
		Name:  submittedForm.Fields[2].Value,
		Aggs:  submittedForm.Fields[3].valueUintSlice,
		Order: submittedForm.Fields[4].valueUint,
	}, &Event{}, editRange)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventAggInsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	eventID := submittedForm.Fields[2].Value
	err := updateDocument(tblEvent, eventID, &Range{
		Name:  submittedForm.Fields[0].Value,
		Aggs:  submittedForm.Fields[1].valueUintSlice,
		IsAgg: true,
	}, &Event{}, eventAddAgg)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}
