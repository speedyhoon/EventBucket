package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	club, err := getClub(clubID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "club")
		return
	}

	action, forms := sessionForms(w, r, clubDetails, clubMoundNew)

	//Club Details Form
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
		Title:    "Club",
		MenuID:   clubID,
		Menu:     urlClubs,
		template: 25,
		Error:    forms[2].Error,
		JS:       []string{"main", "clubSettings", "editRow"},
		Data: map[string]interface{}{
			"Club":        club,
			"debug":       debug,
			"ClubDetails": forms[0],
			"ClubMound":   forms[1],
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
		Title:    "Clubs",
		template: 25,
		JS:       []string{"main"},
		Data: map[string]interface{}{
			"NewClub":   forms[0],
			"ListClubs": listClubs,
			"debug":     debug,
		},
	})
}

func mapClubs(w http.ResponseWriter, r *http.Request, submittedForm form) {
	clubID := submittedForm.Fields[0].Value
	clubs, err := getClubs()
	if err != nil {
		warn.Println(err)
	}
	searchClub := clubID != ""
	var mapClubs []MapClub
	for _, club := range clubs {
		if searchClub && club.ID == clubID || (!searchClub && club.Latitude != 0 && club.Longitude != 0) {
			mapClubs = append(mapClubs, MapClub{
				Name:      club.Name,
				Latitude:  club.Latitude,
				Longitude: club.Longitude,
				URL:       club.URL,
				Address:   club.Address,
				Town:      club.Town,
				Postcode:  club.Postcode,
			})
		}
	}

	var jsonList []byte
	jsonList, err = json.Marshal(mapClubs)
	if err != nil {
		warn.Println(err)
	}
	fmt.Fprintf(w, "%s", jsonList)
}

//MapClub is exported
type MapClub struct {
	Name      string  `json:"n"`
	Latitude  float32 `json:"x,omitempty"`
	Longitude float32 `json:"y,omitempty"`
	URL       string  `json:"u,omitempty"`
	Address   string  `json:"a,omitempty"`
	Town      string  `json:"t,omitempty"`
	Postcode  string  `json:"p,omitempty"`
}

func clubInsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	name := submittedForm.Fields[0].Value
	var ID string
	defaultClub := getDefaultClub()

	//Check if a club with that name already exists.
	club, err := getClubByName(name)
	if err != nil {
		ID, err = insertClub(Club{Name: name, IsDefault: defaultClub.ID == ""}) //Set this club to default if no other clubs are the defualt
		if err != nil {
			formError(w, r, submittedForm, err)
			return
		}
	} else {
		/*Use a generic pageError form to pass the error message to the Club Settings page.
		TODO investigate if there is a simpler way to pass error messages between different pages. Maybe use a slice []string so several messages could be displayed if needed?
		It would also be handy to have success, warning and error statuses */
		setSession(w, form{action: pageError, Error: fmt.Errorf("A club with name '%v' already exists.", name)})
		ID = club.ID
	}
	http.Redirect(w, r, urlClub+ID+"#edit", http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	clubID := submittedForm.Fields[8].Value
	isDefault := submittedForm.Fields[6].Checked
	defaultClub := getDefaultClub()
	if isDefault && defaultClub.ID != clubID {
		//need to remove isDefault for the default club so there is only one default at a time.
		err := updateDocument(tblClub, defaultClub.ID, &Club{IsDefault: false}, &Club{}, updateClubDefault)
		if err != nil {
			warn.Println(err)
		}
	}
	err := updateDocument(tblClub, clubID, &Club{
		Name:      submittedForm.Fields[0].Value,
		Address:   submittedForm.Fields[1].Value,
		Town:      submittedForm.Fields[2].Value,
		Postcode:  submittedForm.Fields[3].Value,
		Latitude:  submittedForm.Fields[4].valueFloat32,
		Longitude: submittedForm.Fields[5].valueFloat32,
		IsDefault: isDefault,
		URL:       submittedForm.Fields[7].Value,
	}, &Club{}, updateClubDetails)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedForm form) {
	clubID := submittedForm.Fields[1].Value
	err := updateDocument(tblClub, clubID, submittedForm.Fields[0].Value, &Club{}, insertClubMound)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}

func trimFloat(num float32) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", num), "0"), ".")
}

func editClubMound(w http.ResponseWriter, r *http.Request, submittedForm form) {
	clubID := submittedForm.Fields[2].Value
	err := updateDocument(tblClub, clubID, &Mound{
		Name: submittedForm.Fields[0].Value,
		ID:   submittedForm.Fields[1].valueUint,
	}, &Club{}, editMound)
	if err != nil {
		formError(w, r, submittedForm, err)
		return
	}
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}
