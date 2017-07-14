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
			"eventShooterNew": shooterEntry,
			"GradesAvailable": pageForms[2],
			"AvailableGrades": grades,
			"AgeGroups":       dataListAgeGroup(),
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, f form) {
	//Try to find an existing club and insert and insert one if it doesn't exist.
	clubID, err := clubInsertIfMissing(f.Fields[0].Value)
	if err != nil {
		formError(w, r, f, err)
		return
	}

	//Insert new event into database.
	ID, err := Event{
		Club:    f.Fields[0].Value,
		ClubID:  clubID,
		Name:    f.Fields[1].Value,
		Date:    f.Fields[2].Value,
		Time:    f.Fields[3].Value,
		Closed:  false,
		AutoInc: AutoInc{Range: 1}, //The next incremental range id to use.
	}.insert()

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlEventSettings+ID, http.StatusSeeOther)
}

func eventAvailableGradesUpsert(w http.ResponseWriter, r *http.Request, f form) {
	eventID := f.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &f.Fields[0].valueUintSlice, &Event{}, updateEventGrades)

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}

func eventShooterInsert(w http.ResponseWriter, r *http.Request, f form) {
	clubID, err := clubInsertIfMissing(f.Fields[2].Value)
	if err != nil {
		formError(w, r, f, err)
		return
	}

	shooter := Shooter{
		FirstName: f.Fields[0].Value,
		NickName:  f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      f.Fields[2].Value,
		ClubID:    clubID,
		AgeGroup:  f.Fields[4].valueUint,
		Sex:       f.Fields[5].Checked,
		Grade:     f.Fields[6].valueUintSlice,
	}
	//Insert shooter into Shooter Bucket
	shooterID, err := shooter.insert()
	if err != nil {
		formError(w, r, f, err)
		return
	}

	//Insert shooter into event
	eventID := f.Fields[7].Value
	shooter.ID = shooterID
	err = updateDocument(tblEvent, eventID, &shooter, &Event{}, eventShooterInsertDB)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}

func eventShooterExistingInsert(w http.ResponseWriter, r *http.Request, f form) {
	eventID := f.Fields[3].Value
	shooter, err := getShooter(f.Fields[0].Value)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	err = updateDocument(tblEvent, eventID, &Shooter{
		ID:        shooter.ID,
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
		Grade:     f.Fields[1].valueUintSlice,
		AgeGroup:  f.Fields[2].valueUint,
		Sex:       shooter.Sex,
	}, &Event{}, eventShooterInsertDB)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}

func eventShooterUpdate(w http.ResponseWriter, r *http.Request, f form) {
	eventID := f.Fields[1].Value
	err := updateDocument(tblEvent, eventID, &EventShooter{
		ID:        f.Fields[0].valueUint,
		FirstName: f.Fields[2].Value,
		Surname:   f.Fields[3].Value,
		Club:      f.Fields[4].Value,
		Grade:     f.Fields[5].valueUint,
		AgeGroup:  f.Fields[6].valueUint,
		Sex:       f.Fields[7].Checked,
		Disabled:  f.Fields[8].Checked,
	}, &Event{}, eventShooterUpdater)

	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlEntries+eventID, http.StatusSeeOther)
}
