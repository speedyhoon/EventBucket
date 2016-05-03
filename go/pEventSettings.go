package main

import (
	"fmt"
	"net/http"
)

func eventSettings(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	var club Club
	if !event.Closed && event.Club != "" {
		club, err = getClubByName(event.Club)
	}

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	action, forms := sessionForms(w, r, eventDetails, eventRangeNew, eventAggNew)
	if action != eventDetails {
		forms[0].Fields[0].Value = event.Name
		forms[0].Fields[1].Value = event.Club
		forms[0].Fields[2].Value = event.Date
		forms[0].Fields[3].Value = event.Time
		forms[0].Fields[4].Checked = event.Closed
		forms[0].Fields[5].Checked = event.AverTwin
		forms[0].Fields[6].Value = event.ID
	}
	forms[1].Fields[1].Value = eventID

	forms[2].Fields[1].Options = dataListRanges(event.Ranges)
	forms[2].Fields[2].Value = eventID

	templater(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Closed":        event.Closed,
			"Ranges":        event.Ranges,
			"Event":         event,
			"EventDetails":  forms[0],
			"AddRange":      forms[1],
			"AddAgg":        forms[2],
			"RangeDataList": club.Mounds,
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
	eventID := submittedForm.Fields[6].Value
	err := updateDocument(tblEvent, eventID, &Event{
		Name:     submittedForm.Fields[0].Value,
		Club:     submittedForm.Fields[1].Value,
		Date:     submittedForm.Fields[2].Value,
		Time:     submittedForm.Fields[3].Value,
		Closed:   submittedForm.Fields[4].Checked,
		AverTwin: submittedForm.Fields[5].Checked,
	}, &Event{}, updateEventDetails)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventRangeInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &Range{Name: submittedForm.Fields[0].Value}, &Event{}, eventAddRange)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventAggInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[2].Value
	err := updateDocument(tblEvent, eventID, &Range{
		Name:  submittedForm.Fields[0].Value,
		Aggs:  submittedForm.Fields[1].valueUintSlice,
		IsAgg: true,
	}, &Event{}, eventAddAgg)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventShooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[7].Value
	err := updateDocument(tblEvent, eventID, &EventShooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[4].valueUint,
		AgeGroup:  submittedForm.Fields[5].valueUint,
		Ladies:    submittedForm.Fields[6].Checked,
	}, &Event{}, eventShooterInsertDB)
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
	err = updateDocument(tblEvent, eventID, &EventShooter{
		Grade:     submittedForm.Fields[1].valueUint,
		AgeGroup:  submittedForm.Fields[2].valueUint,
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
	}, &Event{}, eventShooterInsertDB)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
