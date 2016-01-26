package main

import "net/http"

func event(w http.ResponseWriter, r *http.Request, eventID string) {
	//	eventID := strings.TrimPrefix(r.URL.Path, urlEvent)
	//	if eventID == "" {
	//		http.Redirect(w, r, urlEvents, http.StatusNotFound)
	//	}
	cookie := r.Header.Get("Set-Cookie")
	if cookie != "" {

	}
	if r.URL.Path[len(urlEvent):] == "3C" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	templater(w, page{
		Title:  "Event",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"EventId": eventID,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, "/events", http.StatusSeeOther)
	}
	templater(w, page{
		Title: "Events",
		Data: M{
			"Stuff": "EVENTS page!",
		},
	})
}

func insertEvent(w http.ResponseWriter, r *http.Request) {
	submittedFields, isValid := isValid(r, GlobalForms[0].fields)
	if !isValid {
		setSession(w, form{
			action: "0",
			fields: submittedFields,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/event/1a", http.StatusSeeOther)
}
