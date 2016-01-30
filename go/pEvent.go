package main

import (
	"net/http"
)

func event(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getClub(eventID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		warn.Println(err)
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}
	templater(w, page{
		Title:  "Event",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"Event": event,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r)
	templater(w, page{
		Title: "Events",
		Data: M{
			"NewEvent": eventNewDefaultValues(sessionForm),
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {
	ID, err := getNextID(tblEvent)
	if err != nil {
		//TODO add error problems to form.
		redirect()
		return
	}
	err = upsertDoc(tblEvent, "", Event{
		ID: ID,
		//		Name: name,
	})
	if err != nil {
		//TODO add error problems to form.
		redirect()
		return
	}
	http.Redirect(w, r, urlEvent+ID, http.StatusSeeOther)
}

func eventNewDefaultValues(form form) []field {
	if form.action != eventNew && len(form.fields) == 0 {
		form.fields = []field{
			{Required: hasDefaultClub()},
			{},
			{Value: defaultDate()[0]},
			{Value: defaultTime()[0]},
		}
	}
	return form.fields
}
