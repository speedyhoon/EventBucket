package main

import "net/http"

func clubSettings(w http.ResponseWriter, r *http.Request, clubID string) {
	club, err := getClub(clubID)
	if err != nil {
		warn.Println("club id", clubID, "not found.")
		errorHandler(w, r, http.StatusNotFound, "club")
		return
	}

	detailsForm := getFormSession(w, r, clubDetails)
	detailsForm.Fields[0].Value = club.Name
	detailsForm.Fields[1].Value = club.Address
	detailsForm.Fields[2].Value = club.Town
	detailsForm.Fields[3].Value = club.Postcode
	detailsForm.Fields[4].Value = trimFloat(club.Latitude)
	detailsForm.Fields[5].Value = trimFloat(club.Longitude)
	detailsForm.Fields[6].Value = toB36(club.ID)
	newMoundForm := getFormSession(w, r, clubMoundNew)
	newMoundForm.Fields[0].Value = toB36(club.ID)

	/*
		var invalidForm, detailsForm, newMoundForm form
		invalidForm = getSession(w, r, []uint8{clubDetails})
		if invalidForm.action == clubDetails {
			detailsForm = invalidForm
		} else {
			detailsForm = form{Fields: []field{
				{Value: club.Name},
				{Value: club.Address},
				{Value: club.Town},
				{Value: club.Postcode},
				{Value: fmt.Sprintf("%f", club.Latitude)},
				{Value: fmt.Sprintf("%f", club.Longitude)},
				{Value: toB36(club.ID)},
			}}
		}
		if invalidForm.action == clubMoundNew {
			newMoundForm = invalidForm
		} else {
			newMoundForm = form{Fields: []field{
				{},
				{},
				{Value: toB36(club.ID)},
			}}
		}*/

	templater(w, page{
		Title:   "Club Settings",
		Menu:    urlClubs,
		MenuID:  toB36(club.ID),
		Heading: club.Name,
		Data: M{
			"Club":        club,
			"ClubDetails": detailsForm,
			"ClubMound":   newMoundForm,
		},
	})
}
