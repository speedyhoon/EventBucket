package main

import (
	"net/http"
)

func scoreboard(w http.ResponseWriter, r *http.Request, eventID, rangeID string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, "event")
		return
	}

	uRangeID, err := stoU(rangeID)
	var ranges []Range
	if err == nil {
		ranges = findAggs(uRangeID, event.Ranges)
	}
	if len(ranges) < 1 {
		errorHandler(w, r, "range")
		return
	}

	sortShooters(rangeID).Sort(event.Shooters)

	templater(w, page{
		Title:    "Scoreboard",
		Menu:     urlEvents,
		MenuID:   eventID,
		Heading:  event.Name,
		template: "scoreboard",
		Data: map[string]interface{}{
			"Event":       event,
			"Ranges":      ranges,
			"SortByRange": rangeID,
			"colspan":     5 + len(ranges),
			"medalsLimit": 3,
		},
	})
}

//FindAggs expands any aggregates within the slice supplied
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
