package main

import "net/http"

func event(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)
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
	listEvents, err := getEvents()
	if err != nil {
		warn.Println(err)
	}
	templater(w, page{
		Title: "Events",
		Data: M{
			"NewEvent":   eventNewDefaultValues(sessionForm),
			"ListEvents": listEvents,
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	ID, err := getNextID(tblEvent)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}

	//Insert new event into database.
	err = upsertDoc(tblEvent, ID, Event{
		ID:   ID,
		Club: submittedForm.Fields[0].Value,
		Name: submittedForm.Fields[1].Value,
		Date: submittedForm.Fields[2].Value,
		Time: submittedForm.Fields[3].Value,
	})

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEvent+ID, http.StatusSeeOther)
}

func eventNewDefaultValues(form form) form {
	if form.action != eventNew && len(form.Fields) == 0 {
		form.Fields = []field{
			{Required: hasDefaultClub()},
			{},
			{Value: defaultDate()[0]},
			{Value: defaultTime()[0]},
		}
	}
	return form
}
