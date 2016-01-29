package main

import "net/http"

func clubSettings(w http.ResponseWriter, r *http.Request, clubID string) {
	cookie := r.Header.Get("Set-Cookie")
	if cookie != "" {
	}

	club, err := getClub(clubID)
	if err != nil {
		warn.Println("club id", clubID, "not found.")
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	templater(w, page{
		Title:  "Club Settings",
		menu:   urlClub,
		MenuID: club.ID,
		Data: M{
			"Club": club,
		},
	})
}
