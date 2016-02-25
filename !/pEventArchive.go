package main

import "net/http"

func eventArchive(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Archive",
		Data: M{
			"ListEvents": []Event{},
		},
	})
}
