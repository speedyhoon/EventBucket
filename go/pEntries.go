package main

import "net/http"

func entries(w http.ResponseWriter, r *http.Request, event Event) {
	action, forms := sessionForms(w, r, eventShooterNew, eventShooterExisting)
	shooterEntry := forms[0]
	if action == eventShooterExisting {
		//Existing shooter select box
		shooterEntry.Fields[3].Error = forms[1].Fields[0].Error
		//Grade
		shooterEntry.Fields[6].Error = forms[1].Fields[1].Error
		shooterEntry.Fields[6].Value = forms[1].Fields[1].Value
		//Age Group
		shooterEntry.Fields[4].Error = forms[1].Fields[2].Error
		shooterEntry.Fields[4].Value = forms[1].Fields[2].Value
		//Existing Shooter button
		shooterEntry.Fields[7].Error = forms[1].Fields[3].Error
	}
	shooterEntry.Fields[2].Options = clubsDataList()

	grades := eventGrades(event.Grades)
	shooterEntry.Fields[6].Options = grades
	shooterEntry.Fields[6].Value = event.ID
	shooterEntry.Fields[7].Value = event.ID
	shooterEntry.Fields = append(shooterEntry.Fields, field{Value: event.ID})

	shooterEntry.Fields[3].Options = searchShootersOptions("", "", event.Club.Name)

	//Provide event ID for link to change available grades in event settings
	shooterEntry.Fields[6].Value = event.ID

	templater(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":           event,
			"eventShooterNew": shooterEntry,
			"AvailableGrades": grades,
			"AgeGroups":       dataListAgeGroup(),
		},
	})
}

func eventInsert(f form) (string, error) {
	//Try to find an existing club and insert and insert one if it doesn't exist.
	clubID, err := clubInsertIfMissing(f.Fields[0].Value)
	if err != nil {
		return "", err
	}

	//Insert new event into database.
	ID, err := Event{
		ClubID:  clubID,
		Name:    f.Fields[1].Value,
		Date:    f.Fields[2].Value,
		Time:    f.Fields[3].Value,
		Closed:  false,
		AutoInc: AutoInc{Range: 1}, //The next incremental range id to use.
	}.insert()
	return urlEventSettings + ID, err
}

func eventAvailableGradesUpsert(f form) (string, error) {
	eventID := f.Fields[1].Value
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &f.Fields[0].valueUintSlice, &Event{}, updateEventGrades)
}

func eventShooterInsert(f form) (string, error) {
	//Populate club name if it is empty
	if f.Fields[2].Value == "" {
		f.Fields[2].Value = defaultClubName()
	} else if _, err := clubInsertIfMissing(f.Fields[2].Value); err != nil {
		return "", err
	}

	shooter := Shooter{
		FirstName: f.Fields[0].Value,
		NickName:  f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      f.Fields[2].Value,
		AgeGroup:  f.Fields[4].valueUint,
		Sex:       f.Fields[5].Checked,
		Grades:    f.Fields[6].valueUintSlice,
	}
	//Insert shooter into Shooter Bucket
	shooterID, err := shooter.insert()
	if err != nil {
		return "", err
	}

	//Insert shooter into event
	eventID := f.Fields[7].Value
	shooter.ID = shooterID
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &shooter, &Event{}, eventShooterInsertDB)
}

func eventShooterExistingInsert(f form) (string, error) {
	eventID := f.Fields[3].Value
	shooter, err := getShooter(f.Fields[0].Value)
	if err != nil {
		return "", err
	}
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &Shooter{
			ID:        shooter.ID,
			FirstName: shooter.NickName,
			Surname:   shooter.Surname,
			Club:      shooter.Club,
			Grades:    f.Fields[1].valueUintSlice,
			AgeGroup:  f.Fields[2].valueUint,
			Sex:       shooter.Sex,
		}, &Event{}, eventShooterInsertDB)
}

func eventShooterUpdate(f form) (string, error) {
	eventID := f.Fields[1].Value
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &EventShooter{
			ID:        f.Fields[0].valueUint,
			FirstName: f.Fields[2].Value,
			Surname:   f.Fields[3].Value,
			Club:      f.Fields[4].Value,
			Grade:     f.Fields[5].valueUint,
			AgeGroup:  f.Fields[6].valueUint,
			Sex:       f.Fields[7].Checked,
			Disabled:  f.Fields[8].Checked,
		}, &Event{}, eventShooterUpdater)
}
