package main

import (
	"fmt"
	"net/http"
	"strings"
)

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
	detailsForm.Fields[6].Value = club.ID

	newMoundForm := getFormSession(w, r, clubMoundNew)
	newMoundForm.Fields[0].Value = club.ID

	templater(w, page{
		Title:   "Club Settings",
		Menu:    urlClubs,
		MenuID:  club.ID,
		Heading: club.Name,
		Data: map[string]interface{}{
			"Club":        club,
			"ClubDetails": detailsForm,
			"ClubMound":   newMoundForm,
		},
	})
}

func trimFloat(num float32) string {
	return strings.TrimRight(strings.Trim(fmt.Sprintf("%f", num), "0"), ".")
}
