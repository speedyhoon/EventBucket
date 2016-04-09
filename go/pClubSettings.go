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
		forms[0].Fields[7].Value = club.URL
	}
	forms[0].Fields[8].Value = club.ID

	//Club Mound form
	forms[1].Fields[1].Value = club.ID

	templater(w, page{
		Title:   "Club Settings",
		Menu:    urlClubs,
		MenuID:  club.ID,
		Heading: club.Name,
		Data: map[string]interface{}{
			"Club":        club,
			"ClubDetails": forms[0],
			"ClubMound":   forms[1],
		},
	})
}

func trimFloat(num float32) string {
	//TODO 100 is returned as 1
	return strings.TrimRight(strings.Trim(fmt.Sprintf("%f", num), "0"), ".")
}
