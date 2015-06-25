package main

import (
	"strconv"
	"strings"
)

func addShooterIDsToShooterObjects(eventShooters []EventShooter) []EventShooter {
	//TODO this should be saved to the database & this function removed
	for shooterID := range eventShooters {
		eventShooters[shooterID].ID = shooterID
	}
	return eventShooters
}

func addGradeSeparatorToShooterObject(eventShooters []EventShooter) []EventShooter {
	//Add a boolean field to each shooter in a list of ordered shooters and is true for the first shooter that has a different grade than the last
	var previousShooterGrade int
	for shooterID := range eventShooters {
		if eventShooters[shooterID].Grade != previousShooterGrade {
			eventShooters[shooterID].GradeSeparator = true
			previousShooterGrade = eventShooters[shooterID].Grade
		}
	}
	return eventShooters
}

func scoreboard(url string) Page {
	event, _ := getEvent(strings.Split(url, "/")[0])

	//set the range to sort by
	var sortByRange string
	if event.SortScoreboard != "" {
		sortByRange = event.SortScoreboard
	} else if len(event.Ranges) > 0 {
		sortByRange = "0"
	}

	//TODO eventually remove this!
	//Add shooter ids to the shooter objects
	event.Shooters = addShooterIDsToShooterObjects(event.Shooters)

	sortShooters(sortByRange).Sort(event.Shooters)

	event.Shooters = addGradeSeparatorToShooterObject(event.Shooters)

	intSortByRange, intErr := strconv.Atoi(sortByRange)
	if intErr != nil {
		intSortByRange = -1
	}
	return Page{
		TemplateFile: "scoreboard",
		Title:        "Scoreboard",
		Theme:        templateEmpty,
		Data: M{
			"Title":          "Scoreboard",
			"EventID":        event.ID,
			"EventName":      event.Name,
			"ListShooters":   event.Shooters,
			"ListRanges":     event.Ranges,
			"Css":            "scoreboard.css",
			"Legend":         scoreBoardLegend(),
			"menu":           eventMenu(event.ID, event.Ranges, urlScoreboard, event.IsPrizeMeet),
			"SortByRange":    intSortByRange,
			"SortScoreboard": generateForm(eventSettingsSortScoreboard(event)),
		},
	}
}
