package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/speedyhoon/forms"
	"github.com/speedyhoon/session"
)

func club(w http.ResponseWriter, r *http.Request, club Club) {
	fs, action := session.Forms(w, r, getFields, clubEdit, clubMoundNew)

	if action == clubEdit {
		fs[action].Fields[0].Value = club.Name
		fs[action].Fields[1].Value = club.Address
		fs[action].Fields[2].Value = club.Town
		fs[action].Fields[3].Value = club.Postcode
		fs[action].Fields[4].Value = trimFloat(club.Latitude)
		fs[action].Fields[5].Value = trimFloat(club.Longitude)
		fs[action].Fields[6].Value = club.IsDefault
		fs[action].Fields[6].Disable = club.IsDefault
		fs[action].Fields[7].Value = club.URL
	}

	//always set the clubID
	fs[clubEdit].Fields[8].Value = club.ID

	//Club Mound form
	fs[clubMoundNew].Fields[1].Value = club.ID

	render(w, page{
		Title:   "Club",
		MenuID:  club.ID,
		Menu:    urlClubs,
		skipCSP: true,
		Data: map[string]interface{}{
			"Club":         club,
			"clubEdit":     fs[clubEdit],
			"clubMoundNew": fs[clubMoundNew],
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
		},
	})
}

func clubsMap(w http.ResponseWriter, r *http.Request, f []forms.Field) {
	clubs, err := getMapClubs(f[0].Str())
	if err != nil {
		log.Println(err)
	}

	var list []byte
	list, err = json.Marshal(clubs)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(w, list)
}

func clubInsert(f forms.Form) (string, error) {
	name := f.Fields[0].Str()

	//Check if a club with that name already exists.
	club, ok := getClubByName(name)
	ID := club.ID

	if !ok {
		var err error
		ID, err = Club{
			Name:      name,
			IsDefault: !defaultClub().IsDefault, //Set this club to the default if no other clubs are the default
		}.insert()
		if err != nil {
			return "", err
		}
	}
	return urlClub + ID + "#edit", nil
}

func clubDetailsUpsert(f forms.Form) (string, error) {
	clubID := f.Fields[8].Str()
	isDefault := f.Fields[6].Checked()
	defaultClub := defaultClub()
	if isDefault && defaultClub.ID != clubID {
		//Need to remove isDefault for the default club so there is only one default at a time.
		err := updateDocument(tblClub, defaultClub.ID, &Club{IsDefault: false}, &Club{}, updateClubDefault)
		if err != nil {
			log.Println(err)
		}
	}
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Club{
			Name:      f.Fields[0].Str(),
			Address:   f.Fields[1].Str(),
			Town:      f.Fields[2].Str(),
			Postcode:  f.Fields[3].Str(),
			Latitude:  f.Fields[4].Float(),
			Longitude: f.Fields[5].Float(),
			IsDefault: isDefault,
			URL:       f.Fields[7].Str(),
		}, &Club{}, updateClubDetails)
}

func clubMoundInsert(f forms.Form) (string, error) {
	clubID := f.Fields[1].Str()
	return urlClub + clubID,
		updateDocument(tblClub, clubID, f.Fields[0].Value, &Club{}, insertClubMound)
}

func clubMoundUpsert(f forms.Form) (string, error) {
	clubID := f.Fields[2].Str()
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Mound{
			Name: f.Fields[0].Str(),
			ID:   f.Fields[1].Uint(),
		}, &Club{}, editMound)
}
