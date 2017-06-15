package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func shooters(w http.ResponseWriter, r *http.Request, submittedForm form, isValid bool) {
	_, pageForms := sessionForms(w, r, shooterNew, importShooter)
	shooters, shooterQty, err := getSearchShooters(submittedForm.Fields[0].Value, submittedForm.Fields[1].Value, submittedForm.Fields[2].Value)

	//Search for shooters in the default club if EventBucket was not started in debug mode & all values are empty.
	if submittedForm.Fields[0].Value == "" && submittedForm.Fields[1].Value == "" && submittedForm.Fields[2].Value == "" {
		defaultClub := defaultClubName()
		submittedForm.Fields[2].Value = defaultClub
		submittedForm.Fields[2].Placeholder = defaultClub
	}

	templater(w, page{
		Title: "Shooters",
		Error: err,
		JS:    []string{"main"},
		Data: map[string]interface{}{
			"NewShooter":     pageForms[0],
			"ImportShooters": pageForms[1],
			"ListShooters":   shooters,
			"ShooterSearch":  submittedForm,
			"QtyShooters":    shooterQty,
			"Grades":         globalGradesDataList,
			"AgeGroups":      dataListAgeGroup(),
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form) {
	err := updateDocument(tblShooter, submittedForm.Fields[6].Value, &Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].valueUintSlice,
		AgeGroup:  submittedForm.Fields[4].valueUint,
		Sex:       submittedForm.Fields[5].Checked,
	}, &Shooter{}, updateShooterDetails)
	//Display any insert errors onscreen.
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func eventSearchShooters(w http.ResponseWriter, r *http.Request, submittedForm form) {
	templater(w, page{
		Title:    "Shooter Search",
		template: templateNone,
		Data: map[string]interface{}{
			"ListShooters": searchShootersOptions(submittedForm.Fields[0].Value, submittedForm.Fields[1].Value, submittedForm.Fields[2].Value),
		},
	})
}

func shooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	//Add new club if there isn't already a club with that name
	clubID, err := clubInsertIfMissing(submittedForm.Fields[2].Value)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}

	//Insert new shooter
	_, err = Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		ClubID:    clubID,
		Grade:     submittedForm.Fields[3].valueUintSlice,
		AgeGroup:  submittedForm.Fields[4].valueUint,
		Sex:       submittedForm.Fields[5].Checked,
	}.insert()
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func importShooters(w http.ResponseWriter, r *http.Request) {
	//Form validation doesn't yet have a
	file, _, err := r.FormFile("f")
	if err != nil {
		warn.Println(err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	defer file.Close()

	//Read file contents into bytes buffer.
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	//Convert file source into structs.
	var shooters []Shooter
	err = json.Unmarshal(buf.Bytes(), &shooters)
	if err != nil {
		warn.Println(err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	var clubID string
	//Insert each shooter into database. //TODO look into using a batch write to update the database.
	for _, shooter := range shooters {
		if shooter.Club != "" {
			clubID, err = clubInsertIfMissing(shooter.Club)
			if err != nil {
				warn.Println(err)
			} else {
				shooter.ClubID = clubID
			}
		}

		if _, err = shooter.insert(); err != nil {
			warn.Println(err)
		}
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

//Add new club if there isn't already a club with that name
func clubInsertIfMissing(clubName string) (string, error) {
	club, err := getClubByName(clubName)
	//Club doesn't exist so try to insert it.
	if err != nil {
		return Club{Name: clubName}.insert()
	}
	//return existing club
	return club.ID, err
}

//TODO move into a config file or database?
func dataListAgeGroup() []option {
	//TODO would changing option.Value to an interface reduce the amount of code to convert types?
	return []option{
		{Value: "0", Label: "None"},
		{Value: "1", Label: "U21"},
		{Value: "2", Label: "U25"},
		{Value: "3", Label: "Veteran"},
		{Value: "4", Label: "Super Veteran"},
	}
}
