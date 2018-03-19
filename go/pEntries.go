package main

import (
	"net/http"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/forms"
)

func entries(w http.ResponseWriter, r *http.Request, event Event) {
	action, f := session.Forms(w, r, getForm, eventShooterNew, eventShooterExisting)
	shooterEntry := f[0]
	if action == eventShooterExisting {
		//Existing shooter select box
		shooterEntry.Fields[3].Error = f[1].Fields[0].Error
		//Grade
		shooterEntry.Fields[6].Error = f[1].Fields[1].Error
		shooterEntry.Fields[6].Value = f[1].Fields[1].Value
		//Age Group
		shooterEntry.Fields[4].Error = f[1].Fields[2].Error
		shooterEntry.Fields[4].Value = f[1].Fields[2].Value
		//Existing Shooter button
		shooterEntry.Fields[7].Error = f[1].Fields[3].Error
	}
	shooterEntry.Fields[2].Options = clubsDataList()

	grades := eventGrades(event.Grades)
	shooterEntry.Fields[6].Options = grades
	shooterEntry.Fields[6].Value = event.ID
	shooterEntry.Fields[7].Value = event.ID
	shooterEntry.Fields = append(shooterEntry.Fields, forms.Field{Value: event.ID})

	shooterEntry.Fields[3].Options = searchShootersOptions("", "", event.Club.Name)

	//Provide event ID for link to change available grades in event settings
	shooterEntry.Fields[6].Value = event.ID

	render(w, page{
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

func eventInsert(f forms.Form) (string, error) {
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

func eventAvailableGradesUpsert(f forms.Form) (string, error) {
	eventID := f.Fields[1].Value
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &f.Fields[0].ValueUintSlice, &Event{}, updateEventGrades)
}

func eventShooterInsert(f forms.Form) (string, error) {
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
		AgeGroup:  f.Fields[4].ValueUint,
		Sex:       f.Fields[5].Checked,
		Grades:    f.Fields[6].ValueUintSlice,
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

func eventShooterExistingInsert(f forms.Form) (string, error) {
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
			Grades:    f.Fields[1].ValueUintSlice,
			AgeGroup:  f.Fields[2].ValueUint,
			Sex:       shooter.Sex,
		}, &Event{}, eventShooterInsertDB)
}

func eventShooterUpdate(f forms.Form) (string, error) {
	eventID := f.Fields[1].Value
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &EventShooter{
			ID:        f.Fields[0].ValueUint,
			FirstName: f.Fields[2].Value,
			Surname:   f.Fields[3].Value,
			Club:      f.Fields[4].Value,
			Grade:     f.Fields[5].ValueUint,
			AgeGroup:  f.Fields[6].ValueUint,
			Sex:       f.Fields[7].Checked,
			Disabled:  f.Fields[8].Checked,
		}, &Event{}, eventShooterUpdater)
}
