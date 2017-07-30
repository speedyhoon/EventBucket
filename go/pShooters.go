package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func shooters(w http.ResponseWriter, r *http.Request, f form) {
	_, forms := sessionForms(w, r, shooterNew, shootersImport)
	shooters, err := getSearchShooters(f.Fields[0].Value, f.Fields[1].Value, f.Fields[2].Value)

	//Search for shooters in the default club if EventBucket was not started in debug mode & all values are empty.
	if f.Fields[0].Value == "" && f.Fields[1].Value == "" && f.Fields[2].Value == "" {
		f.Fields[2].Value = defaultClubName()
		f.Fields[2].Placeholder = f.Fields[2].Value
	}

	templater(w, page{
		Title: "Shooters",
		Error: err,
		Data: map[string]interface{}{
			"shooterNew":     forms[0],
			"shootersImport": forms[1],
			"shooterList":    shooters,
			"ShooterSearch":  f,
			"Grades":         globalGradesDataList,
			"AgeGroups":      dataListAgeGroup(),
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, f form) {
	err := updateDocument(tblShooter, f.Fields[6].Value, &Shooter{
		FirstName: f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      f.Fields[2].Value,
		Grades:    f.Fields[3].valueUintSlice,
		AgeGroup:  f.Fields[4].valueUint,
		Sex:       f.Fields[5].Checked,
	}, &Shooter{}, updateShooterDetails)
	//Display any insert errors onscreen.
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func eventSearchShooters(w http.ResponseWriter, r *http.Request, f form) {
	templater(w, page{
		template: "shooterSearch",
		Data: map[string]interface{}{
			"shooterList": searchShootersOptions(f.Fields[0].Value, f.Fields[1].Value, f.Fields[2].Value),
		},
	})
}

func shooterInsert(w http.ResponseWriter, r *http.Request, f form) {
	//Add new club if there isn't already a club with that name
	clubID, err := clubInsertIfMissing(f.Fields[2].Value)
	if err != nil {
		formError(w, r, f, err)
		return
	}

	//Insert new shooter
	_, err = Shooter{
		FirstName: f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      clubID,
		Grades:    f.Fields[3].valueUintSlice,
		AgeGroup:  f.Fields[4].valueUint,
		Sex:       f.Fields[5].Checked,
	}.insert()
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func importShooters(w http.ResponseWriter, r *http.Request, f form) {
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
	_, err = buf.ReadFrom(file)
	if err != nil{
		warn.Println(err)
	}

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
				shooter.Club = clubID
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
	club, ok := getClubByName(clubName)
	if ok{
		//return existing club
		return club.ID, nil
	}	
	//Club doesn't exist so try to insert it.
		return Club{Name: clubName}.insert()
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
