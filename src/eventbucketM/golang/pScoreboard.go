package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func addShooterIdsToShooterObjects(eventShooters []EventShooter) []EventShooter {
	for shooterId := range eventShooters {
		eventShooters[shooterId].Id = shooterId
	}
	return eventShooters
}

func scoreboard(url string) Page {
	event, _ := getEvent(strings.Split(url, "/")[0])

	//set the range to sort by
	var sortByRange string
	if event.SortScoreboard != "" {
		sortByRange = event.SortScoreboard
	} else if len(event.Ranges) >= 1 {
		sortByRange = "0"
	}

	//Add shooter ids to the shooter objects
	event.Shooters = addShooterIdsToShooterObjects(event.Shooters)

	if sortByRange != "" {
		//Closures that order the Change structure.
		grade := func(c1, c2 *EventShooter) bool {
			return c1.Grade < c2.Grade
		}
		total := func(c1, c2 *EventShooter) bool {
			return c1.Scores[sortByRange].Total > c2.Scores[sortByRange].Total
		}
		centres := func(c1, c2 *EventShooter) bool {
			return c1.Scores[sortByRange].Centers > c2.Scores[sortByRange].Centers
		}
		cb := func(c1, c2 *EventShooter) bool {
			return c1.Scores[sortByRange].CountBack1 > c2.Scores[sortByRange].CountBack1
		}
		OrderedBy(grade, total, centres, cb).Sort(event.Shooters)
	}

	previousGrade := -1
	shooterQty := len(event.Shooters)
	var position, shouldBePosition int
	var shootEqual bool
	var positionEqual string
	var thisShooterScore Score

	//Loop through all the shooters
	for index, shooter := range event.Shooters {
		shouldBePosition += 1
		if shooter.Grade != previousGrade {
			previousGrade = shooter.Grade
			//reset position back to 1st when the grade has changed
			position = 1
			shouldBePosition = 1

			//Add the grade separator row in the table
			event.Shooters[index].GradeSeparator = true
		} else if !shootEqual {
			position = shouldBePosition
			positionEqual = ""
		}

		if shootEqual {
			positionEqual = "="
			shootEqual = false
		}

		thisShooterScore = shooter.Scores[sortByRange]
		if thisShooterScore.Total == 0 && thisShooterScore.Centers == 0 {
			event.Shooters[index].Warning = LEGEND_NO_SCORE
			//Set the colour for the cell as well
			thisShooterScore.Warning = LEGEND_NO_SCORE
			shooter.Scores[sortByRange] = thisShooterScore
			if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
				position = 0
			}
		}

		//Calculate if there is a shoot off needed for the next shooter
		if index+1 < shooterQty {
			//Cache the next shooters details
			next_shooter := event.Shooters[index+1]
			nextShooterScore := next_shooter.Scores[sortByRange]

			//Check if the scores are exactly the same
			thisShooterScore.Shots = ""
			nextShooterScore.Shots = ""
			if shooter.Grade == next_shooter.Grade && thisShooterScore == nextShooterScore {
				positionEqual = "="
				shootEqual = true
				if thisShooterScore.Total != 0 {
					event.Shooters[index].Warning = LEGEND_SHOOT_OFF
					event.Shooters[index+1].Warning = LEGEND_SHOOT_OFF
					//Set the colour for the cell as well
					thisShooterScore.Warning = LEGEND_SHOOT_OFF
					shooter.Scores[sortByRange] = thisShooterScore
					nextShooterScore.Warning = LEGEND_SHOOT_OFF
					event.Shooters[index+1].Scores[sortByRange] = nextShooterScore
				}
			}
		}

		//generate the shooters position e.g. "=33rd"
		if position > 0 {
			event.Shooters[index].Position = fmt.Sprintf("%v%v", positionEqual, ordinal(position))
		}
	}

	intSortByRange, intErr := strconv.Atoi(sortByRange)
	if intErr != nil {
		intSortByRange = -1
	}
	return Page{
		TemplateFile: "scoreboard",
		Title:        "Scoreboard",
		Theme:        TEMPLATE_EMPTY,
		Data: M{
			"Title":          "Scoreboard",
			"EventId":        event.Id,
			"EventName":      event.Name,
			"ListShooters":   event.Shooters,
			"ListRanges":     event.Ranges,
			"Css":            "scoreboard.css",
			"Legend":         scoreBoardLegend(),
			"menu":           eventMenu(event.Id, event.Ranges, URL_scoreboard, event.IsPrizeMeet),
			"SortByRange":    intSortByRange,
			"SortScoreboard": generateForm(eventSettingsSortScoreboard(event)),
		},
	}
}

type lessFunc func(p1, p2 *EventShooter) bool

type multiSorter struct {
	changes []EventShooter
	less    []lessFunc
}

func (ms *multiSorter) Sort(changes []EventShooter) {
	ms.changes = changes
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.changes)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}
