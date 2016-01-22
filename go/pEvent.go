package main

import "net/http"

func event(w http.ResponseWriter, r *http.Request, eventId string) {
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
		Title: "Event",
		Data: M{
			"EventId": eventId,
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
	if r.Method != "POST" {
		/*405 Method Not Allowed
		A request was made of a resource using a request method not supported by that resource; for example,
		using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
		//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/

		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}
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
