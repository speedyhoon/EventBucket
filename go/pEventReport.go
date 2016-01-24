package main

import "net/http"

func eventReport(w http.ResponseWriter, r *http.Request, eventId string) {
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
		Title: "Event Report",
		Data: M{
			"EventId": eventId,
		},
	})
}
