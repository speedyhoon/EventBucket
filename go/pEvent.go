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
	formID := 0
	submittedFields, isValid := isValid(r, GlobalForms[formID].fields)
	goToPage := func() { http.Redirect(w, r, "/", http.StatusSeeOther) }
	if !isValid {
		setSession(w, form{
			action: formID,
			fields: submittedFields,
		})
		goToPage()
		return
	}

	ID, err := getNextID(tblEvent)
	if err != nil {
		//TODO add error problems to form.
		goToPage()
		return
	}
	err = upsertDoc(tblEvent, "", Event{
		ID: ID,
		//		Name: name,
	})
	if err != nil {
		//TODO add error problems to form.
		goToPage()
		return
	}
	http.Redirect(w, r, "/event/"+ID, http.StatusSeeOther)
}
