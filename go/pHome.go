package main

import (
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	listEvents, err := getEvents()
	templater(w, page{
		Title: "Home",
		Error: err,
		Data: map[string]interface{}{
			"NewEvent":   getFormSession(w, r, eventNew),
			"ListEvents": listEvents,
		},
	})
}

func report(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Report",
	})
}

func about(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	templater(w, page{
		Title: "About",
		Data: map[string]interface{}{
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
