package main

import "net/http"

func clubSettings(w http.ResponseWriter, r *http.Request, clubId string) {
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
		Title: "Club Settings",
		Data: M{
			"ClubId": clubId,
		},
	})
}

//parameters don't match the regex string - 404 enent id not found
//			errorHandler(w, r, http.StatusNotFound)

//club settings
//whoops an error occured
// that club id you supplied doesn't match anything
//here is a list of valid clubs - that link to the clubsettings page.
func whoops(w http.ResponseWriter, r *http.Request, url string) {
	var pageName, pageType string
	parameterType := "ID"
	switch url {
	case urlClubSettings:
		pageName = "Club Settings"
		pageType = "club"
	case urlEvent:
		pageName = "Event"
		pageType = "event"
	case urlEventSettings:
		pageName = "Event Settings"
		pageType = "event"
	}
	templater(w, page{
		Title: "noId",
		Data: M{
			"PageName":      pageName,
			"PageType":      pageType,
			"ParameterType": parameterType,
			"List":          "no data available right now",
		},
	})
}
