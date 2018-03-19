package main

import (
	"net/http"
	"github.com/speedyhoon/forms"
	"github.com/speedyhoon/session"
)

func shooters(w http.ResponseWriter, r *http.Request, form forms.Form) {
	_, f := session.Forms(w, r, getForm, shooterNew, shootersImport, shooterSearch)

	render(w, page{
		Title: "Shooters",
		Data: map[string]interface{}{
			"shooterNew":     f[0],
			"shootersImport": f[1],
			"shooterSearch":  form,
			"Shooters":       searchShooters(form.Fields[0].Value, form.Fields[1].Value, form.Fields[2].Value),
			"qty":            tblQty(tblShooter),
			"Grades":         globalGradesDataList,
			"AgeGroups":      dataListAgeGroup(),
		},
	})
}

func shooterUpdate(f forms.Form) (string, error) {
	return "", updateDocument(tblShooter, f.Fields[6].Value, &Shooter{
		FirstName: f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      f.Fields[2].Value,
		Grades:    f.Fields[3].ValueUintSlice,
		AgeGroup:  f.Fields[4].ValueUint,
		Sex:       f.Fields[5].Checked,
	}, &Shooter{}, updateShooterDetails)
}

func eventSearchShooters(w http.ResponseWriter, r *http.Request, f forms.Form) {
	render(w, page{
		template: "shooterSearch",
		Data: map[string]interface{}{
			"shooters": searchShooters(f.Fields[0].Value, f.Fields[1].Value, f.Fields[2].Value),
		},
	})
}

func shooterInsert(f forms.Form) (string, error) {
	//Add new club if there isn't already a club with that name
	clubID, err := clubInsertIfMissing(f.Fields[2].Value)
	if err != nil {
		return "", err
	}

	//Insert new shooter
	_, err = Shooter{
		FirstName: f.Fields[0].Value,
		Surname:   f.Fields[1].Value,
		Club:      clubID,
		Grades:    f.Fields[3].ValueUintSlice,
		AgeGroup:  f.Fields[4].ValueUint,
		Sex:       f.Fields[5].Checked,
	}.insert()
	return "", err
}

/*func importShooters(f forms.Form) (string, error) {
	//Form validation doesn't yet have a
	file, _, err := r.FormFile("f")
	if err != nil {
		warn.Println(err)
		return "", err
	}
	defer file.Close()

	//Read file contents into bytes buffer.
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		warn.Println(err)
	}

	//Convert file source into structs.
	var shooters []Shooter
	err = json.Unmarshal(buf.Bytes(), &shooters)
	if err != nil {
		warn.Println(err)
		return "", err
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
	return "", err
}*/

//Add new club if there isn't already a club with that name
func clubInsertIfMissing(clubName string) (string, error) {
	club, ok := getClubByName(clubName)
	if ok {
		//return existing club
		return club.ID, nil
	}
	//Club doesn't exist so try to insert it.
	return Club{Name: clubName}.insert()
}

//TODO move into a config file or database?
func dataListAgeGroup() []forms.Option {
	//TODO would changing option.Value to an interface reduce the amount of code to convert types?
	return []forms.Option{
		{},
		//{Value: "0"}, //None = 0
		//{Value: "", Label: "None"}, //None = 0
		{Value: "1", Label: "U21"},
		{Value: "2", Label: "U25"},
		{Value: "3", Label: "Veteran"},
		{Value: "4", Label: "Super Veteran"},
	}
}
