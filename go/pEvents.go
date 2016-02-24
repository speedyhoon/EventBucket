package main

import "net/http"

func events(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r, []uint8{eventNew, eventDetails})
	listEvents, err := getEvents()
	templater(w, page{
		Title: "Events",
		Error: err,
		Data: M{
			"NewEvent":   eventNewDefaultValues(sessionForm),
			"ListEvents": listEvents,
		},
	})
}

func eventNewDefaultValues(form form) form {
	if len(form.Fields) == 0 {
		form.Fields = []field{
			{Required: hasDefaultClub()},
			{},
			{},
		}
	}
	if form.Fields[2].Value == "" {
		form.Fields[2].Value = defaultDateTime()[0]
	}
	return form
}
