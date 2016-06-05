package main

import (
	"math"
	"net/http"
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

	rangeID, err := stoU(ids[1])
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

	sortShooters(ids[1]).Sort(event.Shooters)

	event.Shooters = addGradeSeparatorToShooterObject(event.Shooters)

	ranges := findAggs(rangeID, event.Ranges)
	if len(ranges) < 1 {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

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
			"Colspan":     5 + len(ranges),
		},
	})
}

func addGradeSeparatorToShooterObject(eventShooters []EventShooter) []EventShooter {
	//Add a boolean field to each shooter in a list of ordered shooters and is true for the first shooter that has a different grade than the last
	var previousShooterGrade uint = math.MaxUint32
	var previousShooterClass uint = math.MaxUint32

	for shooterID := range eventShooters {
		if eventShooters[shooterID].Grade != previousShooterGrade {
			eventShooters[shooterID].GradeSeparator = true
			previousShooterGrade = eventShooters[shooterID].Grade

			if globalGrades[eventShooters[shooterID].Grade].ClassID != previousShooterClass {
				eventShooters[shooterID].ClassSeparator = true
				previousShooterClass = globalGrades[eventShooters[shooterID].Grade].ClassID
			}
		}
	}
	return eventShooters
}

type legend struct {
	//To access a field in HTML a struct, it must start with an uppercase letter. Other wise it will output error: xxx is an unexported field of struct type main.legend
	Class, Name string
}

func findAggs(rangeID uint, ranges []Range) []Range {
	var rs []Range
	for _, r := range ranges {
		if r.ID == rangeID {
			if r.IsAgg {
				for _, id := range r.Aggs {
					rs = append(rs, findAggs(id, ranges)...)
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
		{Class: "a", Name: "Sorted By"},
		{Class: "l", Name: "Highest Possible Score"},
		{Class: "o", Name: "Shoot Off"},
		{Class: "n", Name: "Incomplete Score"},
		{Class: "m", Name: "Missing Score"},
		{Class: "d", Name: "1st"},
		{Class: "e", Name: "2nd"},
		{Class: "f", Name: "3rd"},
		/*{Class: "g", Name: "4th"},
		{Class: "h", Name: "5th"},
		{Class: "i", Name: "6th"},
		{Class: "j", Name: "7th"},
		{Class: "k", Name: "8th"},*/
	}
}
