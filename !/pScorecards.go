package main

import "net/http"

func scorecards(w http.ResponseWriter, r *http.Request, eventID string) {
	//	eventID := strings.TrimPrefix(r.URL.Path, urlEvent)
	//	if eventID == "" {
	//		http.Redirect(w, r, urlEvents, http.StatusNotFound)
	//	}
	/*cookie := r.Header.Get("Set-Cookie")
	if cookie != "" {

	}
	if r.URL.Path[len(urlEvent):] == "3C" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}*/
	templater(w, page{
		Title:  "Scorecards",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"EventId": eventID,
			"Hi":      "boo",
		},
	})
}
