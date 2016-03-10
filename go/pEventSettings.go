package main

import "net/http"

//TODO fix event settings details form
//TODO fix event settings add range form
//TODO fix event settings add agg form
//TODO event details isn't getting default data

func eventSettings(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}
	ranges := dataListRanges(event.Ranges)

	action, forms := sessionForms2(w, r, eventDetails, eventRangeNew, eventAggNew)
	if action == nil || action != nil && *action != eventDetails {
		forms[0].Fields[0].Value = event.Name
		forms[0].Fields[1].Value = event.Club
		forms[0].Fields[2].Value = event.Date
		forms[0].Fields[3].Value = event.Time
		forms[0].Fields[4].Checked = event.Closed
		forms[0].Fields[5].Value = event.ID
	}

	templater(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":       event.Ranges,
			"EventDetails": forms[0],
			"RangesQty":    len(ranges),
			"AddRange":     forms[1],
			"AddAgg":       forms[2],
		},
	})
}

func dataListRanges(ranges []Range) []option {
	var options []option
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, option{Label: r.Name})
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
	err := eventAddRange(eventID, Range{Name: submittedForm.Fields[0].Value})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventAggInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[2].Value
	err := eventAddRange(eventID, Range{Name: submittedForm.Fields[0].Value, IsAgg: true})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

/*


func rangeAggInsert(validatedValues map[string]string, isAgg bool) {
	newRange := Range{Name: validatedValues["name"]}
	if isAgg {
		newRange.IsAgg = true
		newRange.Aggregate = validatedValues["agg"]
	}
	eventID := validatedValues["eventid"]
	rangeID, eventData := eventAddRange(eventID, newRange)
	go calcNewAggRangeScores(eventID, rangeID, eventData)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsAddAggForm("", []Option{}).inputs, r)
	rangeAggInsert(validatedValues, true)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}
*/
func eventShooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[6].Value
	shooter := EventShooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[4].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[5].internalValue.(uint64),
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
		Grade:     submittedForm.Fields[1].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[2].internalValue.(uint64),
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
