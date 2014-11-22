package main

import(
	"net/http"
	"strings"
)

func range_report(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(get_id_from_url(r, URL_rangeReport), "/")
	event_id := arr[0]
	range_id := arr[1]
	event, _ := getEvent(event_id)

	// Closures that order the Change structure.
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	total := func(c1, c2 *EventShooter) bool {
		return c1.Scores[range_id].Total > c2.Scores[range_id].Total
	}
	centa := func(c1, c2 *EventShooter) bool {
		return c1.Scores[range_id].Centers > c2.Scores[range_id].Centers
	}
	cb := func(c1, c2 *EventShooter) bool {
		return c1.Scores[range_id].CountBack1 > c2.Scores[range_id].CountBack1
	}

	var shooter_list []EventShooter
	for shooter_id, shooterList := range event.Shooters {
		shooterList.Id = shooter_id
		/*for range_id, score := range shooterList.Scores {
			//			vardump(score)
			//			export(score)
			//			score.Position = 0
			shooterList.Scores[range_id] = score
			//			dump("\n\n\n")
		}*/
		shooter_list = append(shooter_list, shooterList)
		//		vardump(shooterList)
	}

	OrderedBy(grade, total, centa, cb).Sort(shooter_list)

	map_template := map[string]interface{}{
		"Range_id": range_id,
		"Range_name": event.Ranges[range_id].Name,
		"Event_name": event.Name,
		"List_shooters": shooter_list,
	}



	templator(TEMPLATE_ADMIN, "range-report", map_template, w)
}
