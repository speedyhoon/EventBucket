package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/utl"
)

func club(w http.ResponseWriter, r *http.Request, club Club) {
	fs, action := session.Get(w, r, getFields, frmClubEdit, frmClubMoundNew)

	if action == frmClubEdit {
		fs[action].Fields[0].Value = club.Name
		fs[action].Fields[1].Value = club.Address
		fs[action].Fields[2].Value = club.Town
		fs[action].Fields[3].Value = club.Postcode
		fs[action].Fields[4].Value = utl.TrimFloat(club.Latitude)
		fs[action].Fields[5].Value = utl.TrimFloat(club.Longitude)
		fs[action].Fields[6].Value = club.IsDefault
		fs[action].Fields[6].Disable = club.IsDefault
		fs[action].Fields[7].Value = club.URL
	}

	//always set the clubID
	fs[frmClubEdit].Fields[8].Value = club.ID

	//Club Mound form
	fs[frmClubMoundNew].Fields[1].Value = club.ID

	render(w, page{
		Title:  "Club",
		MenuID: club.ID,
		Menu:   urlClubs,
		//#ifndef DEBUG
		skipCSP: true,
		//#endif
		Data: map[string]interface{}{
			"Club":         club,
			"clubEdit":     fs[frmClubEdit],
			"clubMoundNew": fs[frmClubMoundNew],
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	clubs, err := getClubs()
	f, _ := session.Get(w, r, getFields, frmClubNew)
	render(w, page{
		Title: "Clubs",
		Error: err,
		//#ifndef DEBUG
		skipCSP: true,
		//#endif
		Data: map[string]interface{}{
			"clubNew": f[0],
			"clubs":   clubs,
		},
	})
}

func clubsMap(w http.ResponseWriter, _ *http.Request, f []frm.Field) {
	clubs, err := getMapClubs(f[0].Str())
	if err != nil {
		log.Println(err)
	}

	var list []byte
	list, err = json.Marshal(clubs)
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprint(w, list)
	if err != nil {
		log.Println(err)
	}
}

func clubInsert(f frm.Form) (ID string, err error) {
	ID, err = clubInsertIfNone(f.Fields[0].Str())
	if err != nil {
		return "", err
	}
	return urlClub + ID + "#edit", nil
}

//Add new club if there isn't already a club with that name
func clubInsertIfNone(clubName string) (string, error) {
	club, ok := getClubByName(clubName)
	if ok {
		//return existing club
		return club.ID, nil
	}
	//Club doesn't exist so try to insert it.
	return Club{Name: clubName}.insert()
}

func clubDetailsUpsert(f frm.Form) (string, error) {
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

func clubMoundInsert(f frm.Form) (string, error) {
	clubID := f.Fields[1].Str()
	return urlClub + clubID,
		updateDocument(tblClub, clubID, f.Fields[0].Value, &Club{}, insertClubMound)
}

func clubMoundUpsert(f frm.Form) (string, error) {
	clubID := f.Fields[2].Str()
	return urlClub + clubID,
		updateDocument(tblClub, clubID, &Mound{
			Name: f.Fields[0].Str(),
			ID:   f.Fields[1].Uint(),
		}, &Club{}, editMound)
}
