package main

import (
	"net/http"

	"github.com/speedyhoon/forms"
	"github.com/speedyhoon/session"
)

func entries(w http.ResponseWriter, r *http.Request, event Event) {
	fs, action := session.Get(w, r, getFields, eventShooterNew, eventShooterExisting)
	if action == eventShooterExisting {
		//Existing shooter select box
		fs[action].Fields[3].Err = fs[eventShooterExisting].Fields[0].Err
		//Grade
		fs[action].Fields[6].Err = fs[eventShooterExisting].Fields[1].Err
		fs[action].Fields[6].Value = fs[eventShooterExisting].Fields[1].Value
		//Age Group
		fs[action].Fields[4].Err = fs[eventShooterExisting].Fields[2].Err
		fs[action].Fields[4].Value = fs[eventShooterExisting].Fields[2].Value
		//Existing Shooter button
		fs[action].Fields[7].Err = fs[eventShooterExisting].Fields[3].Err
	}
	fs[eventShooterNew].Fields[2].Options = clubsDataList()

	grades := eventGrades(event.Grades)
	fs[eventShooterNew].Fields[6].Options = grades
	fs[eventShooterNew].Fields[6].Value = event.ID
	fs[eventShooterNew].Fields[7].Value = event.ID
	//TODO what was the below ignored field used for?
	//fs[eventShooterNew].Fields = append(fs[eventShooterNew].Fields, forms.Field{Value: event.ID})

	fs[eventShooterNew].Fields[3].Options = searchShootersOptions("", "", event.Club.Name)

	//Provide event ID for link to change available grades in event settings
	fs[eventShooterNew].Fields[6].Value = event.ID

	render(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":           event,
			"eventShooterNew": fs[eventShooterNew],
			"AvailableGrades": grades,
			"AgeGroups":       dataListAgeGroup(),
		},
	})
}

func eventInsert(f forms.Form) (string, error) {
	//Try to find an existing club and insert and insert one if it doesn't exist.
	clubID, err := clubInsertIfMissing(f.Fields[0].Str())
	if err != nil {
		return "", err
	}

	//Insert new event into database.
	ID, err := Event{
		ClubID:   clubID,
		Name:     f.Fields[1].Str(),
		DateTime: f.Fields[2].Time(),
		Closed:   false,
		AutoInc:  AutoInc{Range: 1}, //The next incremental range id to use.
	}.insert()
	return urlEventSettings + ID, err
}

func eventAvailableGradesUpsert(f forms.Form) (string, error) {
	eventID := f.Fields[1].Str()
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &f.Fields[0].Value, &Event{}, updateEventGrades)
}

func eventShooterInsert(f forms.Form) (string, error) {
	//Populate club name if it is empty
	if f.Fields[2].Value == "" {
		f.Fields[2].Value = defaultClub().Name
	} else if _, err := clubInsertIfMissing(f.Fields[2].Str()); err != nil {
		return "", err
	}

	shooter := Shooter{
		FirstName: f.Fields[0].Str(),
		NickName:  f.Fields[0].Str(),
		Surname:   f.Fields[1].Str(),
		Club:      f.Fields[2].Str(),
		AgeGroup:  f.Fields[4].Uint(),
		Sex:       f.Fields[5].Checked(),
		Grades:    f.Fields[6].UintSlice(),
	}
	//Insert shooter into Shooter Bucket
	shooterID, err := shooter.insert()
	if err != nil {
		return "", err
	}

	//Insert shooter into event
	eventID := f.Fields[7].Str()
	shooter.ID = shooterID
	return urlEntries + eventID,
		updateDocument(tblEvent, eventID, &shooter, &Event{}, eventShooterInsertDB)
}

func eventShooterExistingInsert(f forms.Form) (string, error) {
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
			Grades:    f.Fields[1].UintSlice(),
			AgeGroup:  f.Fields[2].Uint(),
			Sex:       shooter.Sex,
		}, &Event{}, eventShooterInsertDB)
}

func eventShooterUpdate(f forms.Form) (string, error) {
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
