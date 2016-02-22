package main

import "net/http"

func eventSettings(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}
	ranges := filteredRanges(event.Ranges)

	templater(w, page{
		Title:  "Event Settings",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"Event":     event,
			"RangesQty": len(ranges),
			"AddRange": form{Fields: []field{
				{},
				{Value: toB36(event.ID)},
			},
			},
			"AddAgg": form{Fields: []field{
				{},
				{Options: ranges},
				{Value: toB36(event.ID)},
			},
			},
		},
	})
}

func filteredRanges(ranges []Range) []option {
	var options []option
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, option{Label: r.Name})
		}
	}
	return options
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
	eventID := submittedForm.Fields[5].Value
	shooter := EventShooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[4].internalValue.(uint64),
	}
	err := eventShooterInsertDB(eventID, shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
	}
	http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)
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
	http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)
}
