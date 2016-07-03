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
		//Existing shooter select box
		shooterEntry.Fields[3].Error = pageForms[1].Fields[0].Error
		//Grade
		shooterEntry.Fields[6].Error = pageForms[1].Fields[1].Error
		shooterEntry.Fields[6].Value = pageForms[1].Fields[1].Value
		//Age Group
		shooterEntry.Fields[4].Error = pageForms[1].Fields[2].Error
		shooterEntry.Fields[4].Value = pageForms[1].Fields[2].Value
		//Existing Shooter button
		shooterEntry.Fields[7].Error = pageForms[1].Fields[3].Error
	}
	shooterEntry.Fields[2].Options = clubsDataList()

	grades := eventGrades(event.Grades)
	shooterEntry.Fields[6].Options = grades
	shooterEntry.Fields[6].Value = eventID
	shooterEntry.Fields[7].Value = eventID
	shooterEntry.Fields = append(shooterEntry.Fields, field{Value: eventID})

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
			"GradesAvailable": pageForms[2],
			"AvailableGrades": grades,
			"AgeGroups":       dataListAgeGroup(),
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	//Try to find an existing club and insert and insert one if it doesn't exist.
	clubID, err := clubInsertIfMissing(submittedForm.Fields[0].Value)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
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
	clubID, err := clubInsertIfMissing(submittedForm.Fields[2].Value)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}

	shooter := Shooter{
		FirstName: submittedForm.Fields[0].Value,
		NickName:  submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		ClubID:    clubID,
		AgeGroup:  submittedForm.Fields[4].valueUint,
		Ladies:    submittedForm.Fields[5].Checked,
		Grade:     submittedForm.Fields[6].valueUintSlice,
	}
	//Insert shooter into Shooter Bucket
	shooterID, err := insertShooter(shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}

	//Insert shooter into event
	eventID := submittedForm.Fields[7].Value
	shooter.ID = shooterID
	err = updateDocument(tblEvent, eventID, &shooter, &Event{}, eventShooterInsertDB)
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
	err = updateDocument(tblEvent, eventID, &Shooter{
		ID:        shooter.ID,
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
		Grade:     submittedForm.Fields[1].valueUintSlice,
		AgeGroup:  submittedForm.Fields[2].valueUint,
		Ladies:    shooter.Ladies,
	}, &Event{}, eventShooterInsertDB)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}

func eventShooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	eventID := submittedForm.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &EventShooter{
		ID:        submittedForm.Fields[0].valueUint,
		FirstName: submittedForm.Fields[2].Value,
		Surname:   submittedForm.Fields[3].Value,
		Club:      submittedForm.Fields[4].Value,
		Grade:     submittedForm.Fields[5].valueUint,
		AgeGroup:  submittedForm.Fields[6].valueUint,
		Ladies:    submittedForm.Fields[7].Checked,
		Disabled:  submittedForm.Fields[8].Checked,
	}, &Event{}, eventShooterUpdater)

	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
