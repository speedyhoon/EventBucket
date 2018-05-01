package main

import (
	"net/http"

	"github.com/speedyhoon/forms"
	"github.com/speedyhoon/session"
)

func shooters(w http.ResponseWriter, r *http.Request, fields []forms.Field) {
	fs, _ := session.Forms(w, r, getFields, shooterNew, shootersImport, shooterSearch)

	render(w, page{
		Title: "Shooters",
		Data: map[string]interface{}{
			"shooterNew":     fs[shooterNew],
			"shootersImport": fs[shootersImport],
			"shooterSearch":  forms.Form{Fields: fields},
			"Shooters":       searchShooters(fields[0].Str(), fields[1].Str(), fields[2].Str()),
			"qty":            tblQty(tblShooter),
			"Grades":         globalGradesDataList,
			"AgeGroups":      dataListAgeGroup(),
		},
	})
}

func shooterUpdate(f forms.Form) (string, error) {
	return "", updateDocument(tblShooter, f.Fields[6].Str(), &Shooter{
		FirstName: f.Fields[0].Str(),
		Surname:   f.Fields[1].Str(),
		Club:      f.Fields[2].Str(),
		Grades:    f.Fields[3].UintSlice(),
		AgeGroup:  f.Fields[4].Uint(),
		Sex:       f.Fields[5].Checked(),
	}, &Shooter{}, updateShooterDetails)
}

func eventSearchShooters(w http.ResponseWriter, r *http.Request, f forms.Form) {
	render(w, page{
		template: "shooterSearch",
		Data: map[string]interface{}{
			"shooters": searchShooters(f.Fields[0].Str(), f.Fields[1].Str(), f.Fields[2].Str()),
		},
	})
}

func shooterInsert(f forms.Form) (string, error) {
	//Add new club if there isn't already a club with that name
	clubID, err := clubInsertIfMissing(f.Fields[2].Str())
	if err != nil {
		return "", err
	}

	//Insert new shooter
	_, err = Shooter{
		FirstName: f.Fields[0].Str(),
		Surname:   f.Fields[1].Str(),
		Club:      clubID,
		Grades:    f.Fields[3].UintSlice(),
		AgeGroup:  f.Fields[4].Uint(),
		Sex:       f.Fields[5].Checked(),
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
