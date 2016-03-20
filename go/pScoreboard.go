package main

import "net/http"

func scoreboard(w http.ResponseWriter, r *http.Request, eventID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	templater(w, page{
		Title:   "Scoreboard",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Event":       event,
			"Legend":      scoreBoardLegend(),
			"SortByRange": 3,
		},
	})
}

type legend struct {
	//To access a field in HTML a struct, it must start with an uppercase letter. Other wise it will output error: xxx is an unexported field of struct type main.legend
	Class, Name string
}

func scoreBoardLegend() []legend {
	//Constants are not able to be slices so using a function instead.
	//Using a Legend slice because a map[string]string would render with a random order.
	return []legend{
		{Class: "^highestPossibleScore^", Name: "Highest Possible Score"},
		{Class: "^shootOff^", Name: "Shoot Off"},
		{Class: "^incompleteScore^", Name: "Incomplete Score"},
		{Class: "^noScore^", Name: "Missing Score"},
		{Class: "^p1^", Name: "1st"},
		{Class: "^p2^", Name: "2nd"},
		{Class: "^p3^", Name: "3rd"},
	}
}
