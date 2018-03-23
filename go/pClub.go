package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/forms"
)

func club(w http.ResponseWriter, r *http.Request, club Club) {
	forms, submitted := session.Forms(w, r, getFields, clubEdit, clubMoundNew)

	if submitted.Action == clubEdit {
		submitted.Fields[0].Value = club.Name
		submitted.Fields[1].Value = club.Address
		submitted.Fields[2].Value = club.Town
		submitted.Fields[3].Value = club.Postcode
		submitted.Fields[4].Value = trimFloat(club.Latitude)
		submitted.Fields[5].Value = trimFloat(club.Longitude)
		submitted.Fields[6].Checked = club.IsDefault
		submitted.Fields[6].Disable = club.IsDefault
		submitted.Fields[7].Value = club.URL
	}
	forms[0].Fields[8].Value = club.ID

	//Club Mound form
	forms[1].Fields[1].Value = club.ID

	render(w, page{
		Title:   "Club",
		MenuID:  club.ID,
		Menu:    urlClubs,
		skipCSP: true,
		Data: map[string]interface{}{
			"Club":         club,
			"clubEdit":     forms[0],
			"clubMoundNew": forms[1],
			//#ifdef DEBUG
			//TODO move option into HTML generation
			"showMap":      true,
			//#endif
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	clubs, err := getClubs()
	f, _ := session.Forms(w, r, getFields, clubNew)
	render(w, page{
		Title:   "Clubs",
		skipCSP: true,
		Error:   err,
		Data: map[string]interface{}{
			"clubNew": f[0],
			"clubs":   clubs,
			//#ifdef DEBUG
			//TODO move option into HTML generation
			"showMap": true,
			//#endif
		},
	})
}

func clubsMap(w http.ResponseWriter, r *http.Request, f []forms.Field) {
	clubs, err := getMapClubs(f[0].Value)
	if err != nil {
		warn.Println(err)
	}

	var list []byte
	list, err = json.Marshal(clubs)
	if err != nil {
		warn.Println(err)
	}
	fmt.Fprint(w, list)
}

func clubInsert(f forms.Form) (string, error) {
	name := f.Fields[0].Value

	//Check if a club with that name already exists.
	club, ok := getClubByName(name)
	ID := club.ID

	if !ok {
		var err error
		ID, err = Club{
			Name:      name,
			IsDefault: getDefaultClub().ID == "", //Set this club to the default if no other clubs are the default
		}.insert()
		if err != nil {
			return "", err
		}
	}
	return urlClub + ID + "#edit", nil
}

func clubDetailsUpsert(f forms.Form) (string, error) {
	clubID := f.Fields[8].Value
	isDefault := f.Fields[6].Checked
	defaultClub := getDefaultClub()
	if isDefault && defaultClub.ID != clubID {
		//need to remove isDefault for the default club so there is only one default at a time.
		err := updateDocument(tblClub, defaultClub.ID, &Club{IsDefault: false}, &Club{}, updateClubDefault)
		if err != nil {
			warn.Println(err)
		}
	}
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Club{
			Name:      f.Fields[0].Value,
			Address:   f.Fields[1].Value,
			Town:      f.Fields[2].Value,
			Postcode:  f.Fields[3].Value,
			Latitude:  f.Fields[4].ValueFloat32,
			Longitude: f.Fields[5].ValueFloat32,
			IsDefault: isDefault,
			URL:       f.Fields[7].Value,
		}, &Club{}, updateClubDetails)
}

func clubMoundInsert(f forms.Form) (string, error) {
	clubID := f.Fields[1].Value
	return urlClub + clubID,
		updateDocument(tblClub, clubID, f.Fields[0].Value, &Club{}, insertClubMound)
}

func clubMoundUpsert(f forms.Form) (string, error) {
	clubID := f.Fields[2].Value
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Mound{
			Name: f.Fields[0].Value,
			ID:   f.Fields[1].ValueUint,
		}, &Club{}, editMound)
}
