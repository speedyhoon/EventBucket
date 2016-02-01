package main

import "net/http"

func clubSettings(w http.ResponseWriter, r *http.Request, clubID string) {
	club, err := getClub(clubID)
	if err != nil {
		warn.Println("club id", clubID, "not found.")
		errorHandler(w, r, http.StatusNotFound, "club")
		return
	}

	templater(w, page{
		Title:  "Club Settings",
		menu:   urlClub,
		MenuID: club.ID,
		Data: M{
			"Club": club,
			"ClubDetails": []field{
				{Value: club.Name},
				{Value: club.Address},
				{Value: club.Town},
				{Value: club.PostCode},
				{Value: club.Latitude},
				{Value: club.Longitude},
				{Value: club.ID},
			},
			"ClubMound": []field{
				{},
				{},
				{},
				{Value: club.ID},
			},
		},
	})
}
