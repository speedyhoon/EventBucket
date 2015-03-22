package main

import (
	"fmt"
	"sort"
	"strings"
)

func addShooterIdsToShooterObjects(eventShooters []EventShooter) []EventShooter {
	for shooterId := range eventShooters {
		eventShooters[shooterId].Id = shooterId
	}
	return eventShooters
}

func scoreboard(url string) Page {
	arr := strings.Split(url, "/")

	//TODO handle the get event error properly
	event, _ := getEvent(arr[0])

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
	previousClass := -1
	shooterQty := len(event.Shooters)
	var position, shouldBePosition int
	var shoot_off, shoot_equ bool
	allGrades := grades()
	var thisShooterScore Score
	scoreBoardLegendOnOff := scoreBoardLegend()
	//classHPS := calculateHighestPossibleScores(10)

	//Loop through all the shooters
	for index, shooter := range event.Shooters {
		shouldBePosition += 1
		if shooter.Grade != previousGrade {
			//reset position back to 1st when the grade has changed
			position = 1
			shouldBePosition = 1

			//Add the grade separator row in the table
			event.Shooters[index].GradeSeparator = true
			previousGrade = shooter.Grade

			//Add a class label for the following shooters in the table
			if allGrades[shooter.Grade].ClassId != previousClass {
				previousClass = allGrades[shooter.Grade].ClassId
				event.Shooters[index].ClassSeparator = true
			}
		} else if !shoot_off && !shoot_equ {
			position = shouldBePosition
		}
		var positionEqual string
		if shoot_equ || shoot_off {
			positionEqual = "="
			shoot_equ = false

			if shoot_off {
				shoot_off = false
			}
		}

		thisShooterScore = shooter.Scores[sortByRange]
		if thisShooterScore.Total == 0 && thisShooterScore.Centers == 0 {
			scoreBoardLegendOnOff[LEGEND_NO_SCORE].On = true
			event.Shooters[index].Warning = 2
			if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
				position = 0
			}
		}

		//Add highest possible score warning to shooters score
		/*if classHPS[previousClass].Total == thisShooterScore.Total && classHPS[previousClass].Centers == thisShooterScore.Centers {
			//if thisShooterScore.Centers == 10 && ((thisShooterScore.Total == 60 && allGrades[shooter.Grade].ClassName == "F Class") || (thisShooterScore.Total == 50 && allGrades[shooter.Grade].ClassName == "Target")) {
			//event.Shooters[index].Warning = 4
			//			shooter.Scores[sortByRange].Warning = 4
			thisShooterScore.Warning = 4
			shooter.Scores[sortByRange] = thisShooterScore
			scoreBoardLegendOnOff[LEGEND_HIGHEST_POSSIBLE_SCORE].On = true
		}*/

		//Calculate if there is a shoot off needed for the next shooter
		if index+1 < shooterQty {
			//Cache the next shooters details
			next_shooter := event.Shooters[index+1]
			next_shooter_score := next_shooter.Scores[sortByRange]

			//Check if the scores are exactly the same
			thisShooterScore.Shots = ""
			next_shooter_score.Shots = ""
			if shooter.Grade == next_shooter.Grade && thisShooterScore == next_shooter_score {
				positionEqual = "="
				if thisShooterScore.Total == 0 {
					shoot_equ = true
					if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
						position = 0
					}
				} else {
					shoot_off = true
					event.Shooters[index].Warning = 1
					event.Shooters[index+1].Warning = 1
					scoreBoardLegendOnOff[LEGEND_SHOOT_OFF].On = true
				}
			}
		}

		//generate the shooters position e.g. "=33rd"
		if position > 0 {
			event.Shooters[index].Position = fmt.Sprintf("%v%v", positionEqual, ordinal(position))
		}
	}

	outputer := M{
		"Title":          "Scoreboard",
		"EventId":        event.Id,
		"EventName":      event.Name,
		"ListShooters":   event.Shooters,
		"ListRanges":     event.Ranges,
		"Css":            "scoreboard.css",
		"Legend":         scoreBoardLegendOnOff,
		"menu":           eventMenu(event.Id, event.Ranges, URL_scoreboard, event.IsPrizeMeet),
		"SortByRange":    -1,
		"SortScoreboard": "",
	}
	if len(event.Ranges) >= 1 {
		outputer["SortByRange"], _ = strToInt(sortByRange)
		outputer["SortScoreboard"] = generateForm(eventSettingsSortScoreboard(event))
	}
	return Page{
		TemplateFile: "scoreboard",
		Title:        "Scoreboard",
		Theme:        TEMPLATE_EMPTY,
		Data:         outputer,
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
