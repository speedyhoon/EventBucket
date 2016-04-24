package main

import "net/http"

func entries(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)
	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	action, pageForms := sessionForms(w, r, eventShooterNew, eventShooterExisting, eventAvailableGrades)
	shooterEntry := pageForms[0]
	listClubs, err := getClubs()
	if err != nil {
		shooterEntry.Error = err
	}
	if action == eventShooterExisting {
		shooterEntry.Fields[3].Error = pageForms[1].Fields[0].Error
		//Grade
		shooterEntry.Fields[4].Error = pageForms[1].Fields[1].Error
		shooterEntry.Fields[4].Value = pageForms[1].Fields[1].Value
		//Age Group
		shooterEntry.Fields[5].Error = pageForms[1].Fields[2].Error
		shooterEntry.Fields[5].Value = pageForms[1].Fields[2].Value
		//Add Existing Shooter button
		shooterEntry.Fields[6].Error = pageForms[1].Fields[3].Error
	}
	shooterEntry.Fields[2].Options = dataListClubs(listClubs)

	shooterEntry.Fields[4].Options = eventGrades(event.Grades)
	shooterEntry.Fields[6].Value = eventID
	shooterEntry.Fields[7].Value = eventID

	//AvailableGrades
	pageForms[2].Fields[0].Options = availableGrades(event.Grades)
	pageForms[2].Fields[1].Value = eventID

	templater(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":           event,
			"ShooterEntry":    shooterEntry,
			"AvailableGrades": pageForms[2],
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	//Try to find an existing club
	club, err := getClubByName(submittedForm.Fields[0].Value)
	clubID := club.ID
	if err != nil {
		//Insert a new club
		clubID, err = insertClub(Club{Name: submittedForm.Fields[0].Value})
	}

	//Insert new event into database.
	ID, err := insertEvent(Event{
		Club:    submittedForm.Fields[0].Value,
		ClubID:  clubID,
		Name:    submittedForm.Fields[1].Value,
		Date:    submittedForm.Fields[2].Value,
		Time:    submittedForm.Fields[3].Value,
		Closed:  false,
		AutoInc: AutoInc{Range: 1}, //The next incremental range id to use.
	})

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+ID, http.StatusSeeOther)
}

func eventAvailableGradesUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &submittedForm.Fields[0].valueUintSlice, &Event{}, updateEventGrades)

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
