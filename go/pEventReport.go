package main

import "net/http"

func eventReport(w http.ResponseWriter, r *http.Request, event Event) {
	templater(w, page{
		Title:   "Event Report",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event": event,
		},
	})
}

func shootersReport(w http.ResponseWriter, r *http.Request, event Event) {
	templater(w, page{
		Title:   "Shooters Report",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"EventID":  event,
			"Shooters": event.Shooters,
		},
	})
}

func shooterReport(w http.ResponseWriter, r *http.Request, event Event, shooterID sID) {
	templater(w, page{
		Title:       "Shooter Report",
		Menu:        urlEvents,
		MenuID:      event.ID,
		Heading:     event.Name,
		SubTemplate: "shootersreport",
		Data: map[string]interface{}{
			"Shooters": []EventShooter{event.Shooters[shooterID]},
		},
	})
}
