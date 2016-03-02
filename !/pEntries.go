package main

import "net/http"

func entries(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	action, pageForms := sessionForms(w, r, eventShooterNew, eventShooterExisting)
	shooterEntry := form{Fields: pageForms[0]}
	listClubs, err := getClubs()
	if err != nil {
		shooterEntry.Error = err.Error()
	}
	if action != nil && *action == eventShooterExisting {
		shooterEntry.Fields[3].Error = pageForms[1][0].Error
		//Grade
		shooterEntry.Fields[4].Error = pageForms[1][1].Error
		shooterEntry.Fields[4].Value = pageForms[1][1].Value
		//Age Group
		shooterEntry.Fields[5].Error = pageForms[1][2].Error
		shooterEntry.Fields[5].Value = pageForms[1][2].Value
		//Add Existing Shooter button
		shooterEntry.Fields[6].Error = pageForms[1][3].Error
	}
	shooterEntry.Fields[2].Options = dataListClubs(listClubs)

	shooterEntry.Fields[6].Value = eventID
	shooterEntry.Fields[7].Value = eventID

	templater(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: M{
			"Event":        event,
			"ShooterEntry": shooterEntry,
			"QtyEntries":   len(event.Shooters),
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	//Insert new event into database.
	ID, err := insertEvent(Event{
		Club:     submittedForm.Fields[0].Value,
		Name:     submittedForm.Fields[1].Value,
		DateTime: submittedForm.Fields[2].Value,
	})

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+ID, http.StatusSeeOther)
}
