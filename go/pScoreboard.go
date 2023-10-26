package main

import "net/http"

func scoreboard(w http.ResponseWriter, r *http.Request, event Event, rangeID rID) {
	ranges := findAggs(uint(rangeID), event.Ranges)
	if len(ranges) < 1 {
		errorHandler(w, r, "range")
		return
	}

	sortShooters(rangeID).Sort(event.Shooters)

	render(w, page{
		Title:    "Scoreboard",
		Menu:     urlEvents,
		MenuID:   event.ID,
		Heading:  event.Name,
		template: "scoreboard",
		Data: map[string]interface{}{
			"Event":       event,
			"Ranges":      ranges,
			"RangeName":   event.Ranges[rangeID].Name,
			"SortByRange": rangeID.StrID(),
			"colspan":     5,
			"medalsLimit": 3,
		},
	})
}

// findAggs expands any aggregates within the slice supplied.
func findAggs(rangeID uint, ranges []Range) (rs []Range) {
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
