package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func scoreboard(w http.ResponseWriter, r *http.Request) {
	event_id := get_id_from_url(r, URL_scoreboard)
	templator(TEMPLATE_EMPTY, "scoreboard", scoreboard_Data(event_id), w)
}
func scoreboard_Data(url string) map[string]interface{} {
	//	export(url)
	arr := strings.Split(url, "/")
	event_id := arr[0]

	event, _ := getEvent(event_id)
	var sortByRange string
	if event.SortScoreboard != "" {
		sortByRange = event.SortScoreboard
	} else if len(event.Ranges) >= 1 {
		for event_range := range event.Ranges {
			sortByRange = event_range
			break
		}
	}

	score_board_legend_on_off := make(map[string]bool)
	for _, legend_name := range SCOREBOARD_LEGEND {
		score_board_legend_on_off[legend_name] = false
	}

	// Closures that order the Change structure.
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	total := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Total > c2.Scores[sortByRange].Total
	}
	centa := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].Centers > c2.Scores[sortByRange].Centers
	}
	cb := func(c1, c2 *EventShooter) bool {
		return c1.Scores[sortByRange].CountBack1 > c2.Scores[sortByRange].CountBack1
	}

	var shooter_list []EventShooter
	for shooter_id, shooterList := range event.Shooters {
		shooterList.Id = shooter_id
		for range_id, score := range shooterList.Scores {
			//			vardump(score)
			//			export(score)
//			score.Position = 0
			shooterList.Scores[range_id] = score
			//			dump("\n\n\n")
		}
		shooter_list = append(shooter_list, shooterList)
		//		vardump(shooterList)
	}
	if sortByRange != "" {
		OrderedBy(grade, total, centa, cb).Sort(shooter_list)
	}

	previous_grade := ""
	previous_class := ""
	position := 0
	should_be_position := 0
	shoot_off := false
	shoot_equ := false
	shooter_length := len(shooter_list)
	for index, shooter := range shooter_list {
		should_be_position += 1
		if shooter.Grade != previous_grade {
			//reset position back to 1st
			position = 1
			should_be_position = 1
			shooter_list[index].GradeSeparator = true
			previous_grade = shooter.Grade
			if class_translation(shooter.Grade) != previous_class {
				previous_class = class_translation(shooter.Grade)
				shooter_list[index].ClassSeparator = true
			}
		} else if !shoot_off && !shoot_equ {
			position = should_be_position
		}
		var display string
		if shoot_off {
			score_board_legend_on_off["ShootOff"] = true
			display = "="
			shoot_off = false
			shoot_equ = false
			shooter_list[index].Warning = 1
		}
		if shoot_equ {
			display = "="
			shoot_equ = false
		}

		this_shooter_score := shooter.Scores[sortByRange]
		if SCOREBOARD_SHOW_WARNING_FOR_ZERO_SCORES && this_shooter_score.Total == 0 && this_shooter_score.Centers == 0 {
			score_board_legend_on_off["NoScore"] = true
			shooter_list[index].Warning = 2
			if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
				position = 0
			}
		}
		if this_shooter_score.Centers == 10 && ((this_shooter_score.Total == 60 && class_translation(shooter.Grade) == "F Class") || (this_shooter_score.Total == 50 && class_translation(shooter.Grade) == "Target")) {
			shooter_list[index].Warning = 4
			score_board_legend_on_off["HighestPossibleScore"] = true
		}
		if index+1 < shooter_length {
			next_shooter := shooter_list[index+1]
			next_shooter_score := next_shooter.Scores[sortByRange]
			if shooter.Grade == next_shooter.Grade &&
				this_shooter_score.Total == next_shooter_score.Total &&
				this_shooter_score.Centers == next_shooter_score.Centers &&
				this_shooter_score.CountBack1 == next_shooter_score.CountBack1 {
				display = "="
				if this_shooter_score.Total == 0 {
					shoot_equ = true
					if SCOREBOARD_IGNORE_POSITION_FOR_ZERO_SCORES {
						position = 0
					}
				} else {
					shoot_off = true
					shooter_list[index].Warning = 1
					score_board_legend_on_off["ShootOff"] = true
				}
			}
		}
		if position > 0 {
			shooter_list[index].Position = fmt.Sprintf("%v%v", display, ordinal(position))
		}
	}

	return map[string]interface{}{
		"Title":        "Scoreboard",
		"EventId":      arr[0],
		"EventName":    event.Name,
		"ListShooters": shooter_list,
		"ListRanges":   event.Ranges,
		"Css":          "scoreboard.css",
		"Legend":       render_legend(score_board_legend_on_off),
		"SortByRange":  event.SortScoreboard,
		"menu":         event_menu(event_id, event.Ranges, URL_scoreboard),
	}
}

func render_legend(items_status map[string]bool) string {
	labels := []string{}
	for _, legend_name := range SCOREBOARD_LEGEND {
		if items_status[legend_name] {
			labels = append(labels, fmt.Sprintf("<label class=%v>%v</label>", SCOREBOARD_LEGEND_CSS_CLASSES[legend_name][0], SCOREBOARD_LEGEND_CSS_CLASSES[legend_name][1]))
		}
	}
	return strings.Join(labels, " ")
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
