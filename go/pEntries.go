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
	shooterEntry.Fields[2].Options = clubsDataList()

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
		JS:      []string{"main", "shooterSearch"},
		Data: map[string]interface{}{
			"Event":           event,
			"ShooterEntry":    shooterEntry,
			"GradesAvailable": pageForms[2],
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

	if len(searchShootersOptions(submittedForm.Fields[0].Value, submittedForm.Fields[1].Value, submittedForm.Fields[2].Value)) <= 1 { //search always returns a blank option for html rendering so the select box isn't mandatory.
		shooterInsert(w, r, form{Fields: []field{
			submittedForm.Fields[0],
			submittedForm.Fields[1],
			submittedForm.Fields[2],
			submittedForm.Fields[4],
			submittedForm.Fields[5],
		}}, redirect)
	} else {
		http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
	}
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
