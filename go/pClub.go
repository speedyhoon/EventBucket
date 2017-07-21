package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func club(w http.ResponseWriter, r *http.Request, club Club) {
	action, forms := sessionForms(w, r, clubDetails, clubMoundNew)

	//Club Details form
	if action != clubDetails {
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
		Error:   forms[2].Error,
		Data: map[string]interface{}{
			"Club":        club,
			"debug":       debug,
			"clubDetails": forms[0],
			"clubMound":   forms[1],
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	listClubs, err := getClubs()
	_, forms := sessionForms(w, r, clubNew)
	templater(w, page{
		Title:   "Clubs",
		skipCSP: true,
		Error:   err,
		Data: map[string]interface{}{
			"clubNew":   forms[0],
			"ListClubs": listClubs,
			"debug":     debug,
		},
	})
}

func mapClubs(w http.ResponseWriter, r *http.Request, f form) {
	clubs, err := getMapClubs(f.Fields[0].Value)
	if err != nil {
		warn.Println(err)
	}

	var list []byte
	list, err = json.Marshal(clubs)
	if err != nil {
		warn.Println(err)
	}
	fmt.Fprintf(w, "%s", list)
}

func clubInsert(w http.ResponseWriter, r *http.Request, f form) {
	name := f.Fields[0].Value

	//Check if a club with that name already exists.
	club, ok := getClubByName(name)
	ID := club.ID
	
	//Club doesn't exist so try to insert it
	if !ok {
		var err error
		ID, err = Club{
			Name:      name,
			IsDefault: getDefaultClub().ID == "", //Set this club to the default if no other clubs are the default
		}.insert()
		if err != nil {
			formError(w, r, f, err)
			return
		}
	}
	//Else if a club already exists with the same name, redirect to the existing club
	
	http.Redirect(w, r, urlClub+ID+"#edit", http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, f form) {
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
	err := updateDocument(tblClub, clubID, &Club{
		Name:      f.Fields[0].Value,
		Address:   f.Fields[1].Value,
		Town:      f.Fields[2].Value,
		Postcode:  f.Fields[3].Value,
		Latitude:  f.Fields[4].valueFloat32,
		Longitude: f.Fields[5].valueFloat32,
		IsDefault: isDefault,
		URL:       f.Fields[7].Value,
	}, &Club{}, updateClubDetails)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, f form) {
	clubID := f.Fields[1].Value
	err := updateDocument(tblClub, clubID, f.Fields[0].Value, &Club{}, insertClubMound)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}

func trimFloat(num float32) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", num), "0"), ".")
}

func editClubMound(w http.ResponseWriter, r *http.Request, f form) {
	clubID := f.Fields[2].Value
	err := updateDocument(tblClub, clubID, &Mound{
		Name: f.Fields[0].Value,
		ID:   f.Fields[1].valueUint,
	}, &Club{}, editMound)
	if err != nil {
		formError(w, r, f, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}
