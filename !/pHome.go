package main

import (
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r, []uint8{eventNew})
	listEvents, err := getEvents()
	templater(w, page{
		Title: "Home",
		Error: err,
		Data: M{
			"NewEvent":   eventNewDefaultValues(sessionForm),
			"ListEvents": listEvents,
		},
	})
}

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

func eventArchive(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Archive",
		Data: M{
			"Stuff": "Archive page!",
		},
	})
}

func about(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	templater(w, page{
		Title: "About",
		Data: M{
			"Hostname":    hostname,
			"IpAddresses": localIPs(),
		},
	})
}

func licence(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Licence",
	})
}
