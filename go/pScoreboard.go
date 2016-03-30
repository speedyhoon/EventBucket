package main

import (
	"math"
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

	sortShooters(ids[1]).Sort(event.Shooters)

	event.Shooters = addGradeSeparatorToShooterObject(event.Shooters)

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
			"Colspan":     4 + len(ranges),
		},
	})
}

func addGradeSeparatorToShooterObject(eventShooters []EventShooter) []EventShooter {
	//Add a boolean field to each shooter in a list of ordered shooters and is true for the first shooter that has a different grade than the last
	var previousShooterGrade uint64 = math.MaxUint64
	var previousShooterClass uint64 = math.MaxUint64
	allGrades := grades()
	for shooterID := range eventShooters {
		if eventShooters[shooterID].Grade != previousShooterGrade {
			eventShooters[shooterID].GradeSeparator = true
			previousShooterGrade = eventShooters[shooterID].Grade

			if allGrades[eventShooters[shooterID].Grade].classID != previousShooterClass {
				eventShooters[shooterID].ClassSeparator = true
				previousShooterClass = allGrades[eventShooters[shooterID].Grade].classID
			}
		}
	}
	return eventShooters
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
