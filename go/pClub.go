package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func club(w http.ResponseWriter, r *http.Request, club Club) {
	action, forms := sessionForms(w, r, clubEdit, clubMoundNew)

	//Club Details Form
	if action != clubEdit {
		forms[0].Fields[0].Value = club.Name
		forms[0].Fields[1].Value = club.Address
		forms[0].Fields[2].Value = club.Town
		forms[0].Fields[3].Value = club.Postcode
		forms[0].Fields[4].Value = trimFloat(club.Latitude)
		forms[0].Fields[5].Value = trimFloat(club.Longitude)
		forms[0].Fields[6].Checked = club.IsDefault
		forms[0].Fields[6].Disable = club.IsDefault
		forms[0].Fields[7].Value = club.URL
	}
	forms[0].Fields[8].Value = club.ID

	//Club Mound form
	forms[1].Fields[1].Value = club.ID

	templater(w, page{
		Title:   "Club",
		MenuID:  club.ID,
		Menu:    urlClubs,
		skipCSP: true,
		Data: map[string]interface{}{
			"Club":         club,
			"showMap":      !debug,
			"clubEdit":  forms[0],
			"clubMoundNew": forms[1],
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	clubs, err := getClubs()
	_, forms := sessionForms(w, r, clubNew)
	templater(w, page{
		Title:   "Clubs",
		skipCSP: true,
		Error:   err,
		Data: map[string]interface{}{
			"clubNew": forms[0],
			"clubs":   clubs,
			"showMap": !debug,
		},
	})
}

func clubsMap(w http.ResponseWriter, r *http.Request, f form) {
	clubs, err := getMapClubs(f.Fields[0].Value)
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

func clubInsert(f form) (string, error) {
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

func clubDetailsUpsert(f form) (string, error) {
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
			Latitude:  f.Fields[4].valueFloat32,
			Longitude: f.Fields[5].valueFloat32,
			IsDefault: isDefault,
			URL:       f.Fields[7].Value,
		}, &Club{}, updateClubDetails)
}

func clubMoundInsert(f form) (string, error) {
	clubID := f.Fields[1].Value
	return urlClub + clubID,
		updateDocument(tblClub, clubID, f.Fields[0].Value, &Club{}, insertClubMound)
}

func clubMoundUpsert(f form) (string, error) {
	clubID := f.Fields[2].Value
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Mound{
			Name: f.Fields[0].Value,
			ID:   f.Fields[1].valueUint,
		}, &Club{}, editMound)
}
