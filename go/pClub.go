package main

import (
	"fmt"
	"net/http"
)

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	club, err := getClub(clubID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "club")
		return
	}
	templater(w, page{
		Title:   "Club",
		MenuID:  clubID,
		Menu:    urlClubs,
		Heading: club.Name,
		Data: map[string]interface{}{
			"Club": club,
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	listClubs, err := getClubs()
	if err != nil {
		warn.Println(err)
	}
	_, forms := sessionForms(w, r, clubNew)
	templater(w, page{
		Title: "Clubs",
		Data: map[string]interface{}{
			"NewClub":   forms[0],
			"ListClubs": listClubs,
		},
	})
}

func clubInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	name := submittedForm.Fields[0].Value
	var ID string

	club, err := getClubByName(name)
	if err != nil {
		ID, err = insertClub(Club{Name: name, IsDefault: submittedForm.Fields[1].Checked})
		if err != nil {
			formError(w, submittedForm, redirect, err)
			return
		}
	} else {
		//Use a generic pageError form to pass the error message to the Club Settings page.
		/*TODO investigate if there is a simpler way to pass error messages between different pages. Maybe use a slice []string so several messages could be displayed if needed?
		It would also be handy to have success, warning and error statuses */
		setSession(w, form{action: pageError, Error: fmt.Errorf("A club with name '%v' already exists.", name)})
		ID = club.ID
	}
	http.Redirect(w, r, urlClubSettings+ID, http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[8].Value
	err := updateDocument(tblClub, clubID, &Club{
		Name:      submittedForm.Fields[0].Value,
		Address:   submittedForm.Fields[1].Value,
		Town:      submittedForm.Fields[2].Value,
		Postcode:  submittedForm.Fields[3].Value,
		Latitude:  submittedForm.Fields[4].valueFloat32,
		Longitude: submittedForm.Fields[5].valueFloat32,
		IsDefault: submittedForm.Fields[6].Checked,
		URL:       submittedForm.Fields[7].Value,
	}, &Club{}, updateClubDetails)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[1].Value
	err := updateDocument(tblClub, clubID, &Mound{
		Name: submittedForm.Fields[0].Value,
	}, &Club{}, insertClubMound)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}
