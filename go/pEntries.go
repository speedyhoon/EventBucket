package main

import (
	"net/http"

	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/session"
)

func entries(w http.ResponseWriter, r *http.Request, event Event) {
	fs, action := session.Get(w, r, frmEventShooterNew, frmEventShooterExisting)
	if action == frmEventShooterExisting {
		// Existing shooter select box.
		fs[action].Fields[3].Err = fs[frmEventShooterExisting].Fields[0].Err
		// Grade.
		fs[action].Fields[6].Err = fs[frmEventShooterExisting].Fields[1].Err
		fs[action].Fields[6].Value = fs[frmEventShooterExisting].Fields[1].Value
		// Age Group.
		fs[action].Fields[4].Err = fs[frmEventShooterExisting].Fields[2].Err
		fs[action].Fields[4].Value = fs[frmEventShooterExisting].Fields[2].Value
		// Existing Shooter button.
		fs[action].Fields[7].Err = fs[frmEventShooterExisting].Fields[3].Err
	}
	fs[frmEventShooterNew].Fields[2].Options = clubsDataList()

	grades := eventGrades(event.Grades)
	fs[frmEventShooterNew].Fields[6].Options = grades
	fs[frmEventShooterNew].Fields[6].Value = event.ID
	fs[frmEventShooterNew].Fields[7].Value = event.ID
	// TODO what was the below ignored field used for?
	//fs[frmEventShooterNew].Fields = append(fs[frmEventShooterNew].Fields, frm.Field{Value: event.ID})

	fs[frmEventShooterNew].Fields[3].Options = searchShootersOptions("", "", event.Club.Name)

	// Provide event ID for link to change available grades in event settings.
	fs[frmEventShooterNew].Fields[6].Value = event.ID

	render(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":           event,
			"eventShooterNew": fs[frmEventShooterNew],
			"AvailableGrades": grades,
			"AgeGroups":       dataListAgeGroup(),
		},
	})
}

func eventInsert(f frm.Form) (string, error) {
	// Try to find an existing club and insert and insert one if it doesn't exist.
	clubID, err := clubInsertIfNone(f.Fields[0].Str())
	if err != nil {
		return "", err
	}

	// Insert new event into database.
	ID, err := Event{
		ClubID:   clubID,
		Name:     f.Fields[1].Str(),
		DateTime: f.Fields[2].Time(),
		Closed:   false,
		AutoInc:  AutoInc{Range: 1}, // The next incremental range id to use.
	}.insert()
	return urlEventSettings + ID, err
}

func eventAvailableGradesUpsert(f frm.Form) (string, error) {
	eventID := f.Fields[1].Str()
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &f.Fields[0].Value, &Event{}, updateEventGrades)
}

func eventShooterInsert(f frm.Form) (string, error) {
	// Populate club name if it is empty.
	if f.Fields[2].Value == "" {
		f.Fields[2].Value = defaultClub().Name
	} else if _, err := clubInsertIfNone(f.Fields[2].Str()); err != nil {
		return "", err
	}

	shooter := Shooter{
		FirstName: f.Fields[0].Str(),
		NickName:  f.Fields[0].Str(),
		Surname:   f.Fields[1].Str(),
		Club:      f.Fields[2].Str(),
		AgeGroup:  f.Fields[4].Uint(),
		Sex:       f.Fields[5].Checked(),
		Grades:    f.Fields[6].Uints(),
	}
	// Insert shooter into Shooter Bucket.
	shooterID, err := shooter.insert()
	if err != nil {
		return "", err
	}

	// Insert shooter into event.
	eventID := f.Fields[7].Str()
	shooter.ID = shooterID
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &shooter, &Event{}, eventShooterInsertDB)
}

func eventShooterExistingInsert(f frm.Form) (string, error) {
	eventID := f.Fields[3].Str()
	shooter, err := getShooter(f.Fields[0].Str())
	if err != nil {
		return "", err
	}
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &Shooter{
			ID:        shooter.ID,
			FirstName: shooter.NickName,
			Surname:   shooter.Surname,
			Club:      shooter.Club,
			Grades:    f.Fields[1].Uints(),
			AgeGroup:  f.Fields[2].Uint(),
			Sex:       shooter.Sex,
		}, &Event{}, eventShooterInsertDB)
}

func eventShooterUpdate(f frm.Form) (string, error) {
	eventID := f.Fields[1].Str()
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &EventShooter{
			ID:        f.Fields[0].Uint(),
			FirstName: f.Fields[2].Str(),
			Surname:   f.Fields[3].Str(),
			Club:      f.Fields[4].Str(),
			Grade:     f.Fields[5].Uint(),
			AgeGroup:  f.Fields[6].Uint(),
			Sex:       f.Fields[7].Checked(),
			Disabled:  f.Fields[8].Checked(),
		}, &Event{}, eventShooterUpdater)
}
