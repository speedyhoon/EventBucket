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

func report(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Report",
	})
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
