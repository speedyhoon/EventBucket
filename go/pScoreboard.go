package main

import (
	"net/http"
	"strconv"
	"strings"
)

func scoreboard(w http.ResponseWriter, r *http.Request, parameters string) {
	//eventID/rangeID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	rangeID, err := strconv.ParseUint(ids[1], 10, 64)
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

	ranges := findRanges(rangeID, event.Ranges)

	templater(w, page{
		Title:    "Scoreboard",
		Menu:     urlEvents,
		MenuID:   ids[0],
		Heading:  event.Name,
		template: templateScoreboard,
		Data: map[string]interface{}{
			"Event":       event,
			"Ranges":      ranges,
			"Legend":      scoreBoardLegend(),
			"SortByRange": rangeID,
		},
	})
}

type legend struct {
	//To access a field in HTML a struct, it must start with an uppercase letter. Other wise it will output error: xxx is an unexported field of struct type main.legend
	Class, Name string
}

func findRanges(rangeID uint64, ranges []Range) []Range {
	var rs []Range
	for _, r := range ranges {
		if r.ID == rangeID {
			if r.IsAgg {
				for _, id := range r.Aggs {
					rs = append(rs, findRanges(id, ranges)...)
				}
			}
			rs = append(rs, r)
		}
	}
	return rs
}

func scoreBoardLegend() []legend {
	//Constants are not able to be slices so using a function instead.
	//Using a Legend slice because a map[string]string would render with a random order.
	return []legend{
		{Class: "^sortBy^", Name: "Sorted By"},
		{Class: "^highestPossibleScore^", Name: "Highest Possible Score"},
		{Class: "^shootOff^", Name: "Shoot Off"},
		{Class: "^incompleteScore^", Name: "Incomplete Score"},
		{Class: "^noScore^", Name: "Missing Score"},
		{Class: "^p1^", Name: "1st"},
		{Class: "^p2^", Name: "2nd"},
		{Class: "^p3^", Name: "3rd"},
		{Class: "^p4^", Name: "4th"},
		{Class: "^p5^", Name: "5th"},
		{Class: "^p6^", Name: "6th"},
		{Class: "^p7^", Name: "7th"},
		{Class: "^p8^", Name: "8th"},
	}
}
