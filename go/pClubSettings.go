package main

import "net/http"

func clubSettings(w http.ResponseWriter, r *http.Request, clubID string) {
	club, err := getClub(clubID)
	if err != nil {
		warn.Println("club id", clubID, "not found.")
		errorHandler(w, r, http.StatusNotFound, "club")
		return
	}

	var invalidForm, detailsForm, newMoundForm form
	invalidForm = getSession(w, r)
	if invalidForm.action == clubDetails {
		detailsForm = invalidForm
	} else {
		detailsForm = form{Fields: []field{
			{Value: club.Name},
			{Value: club.Address},
			{Value: club.Town},
			{Value: club.Postcode},
			{Value: club.Latitude},
			{Value: club.Longitude},
			{Value: club.ID},
		}}
	}
	if invalidForm.action == clubMoundNew {
		newMoundForm = invalidForm
	} else {
		newMoundForm = form{Fields: []field{
			{},
			{},
			{},
			{Value: club.ID},
		}}
	}

	templater(w, page{
		Title:  "Club Settings",
		menu:   urlClub,
		MenuID: club.ID,
		Data: M{
			"Club":        club,
			"ClubDetails": detailsForm,
			"ClubMound":   newMoundForm,
		},
	})
}
