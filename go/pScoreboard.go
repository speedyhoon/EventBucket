package main

import (
	"net/http"
)

func scoreboard(w http.ResponseWriter, r *http.Request, eventID, rangeId string) {
	event, err := getEvent(eventID)

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	rangeIDstr, err := stoU(rangeId)
	var ranges []Range
	if err == nil {
		ranges = findAggs(rangeIDstr, event.Ranges)
	}
	if len(ranges) < 1 {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

	sortShooters(rangeId).Sort(event.Shooters)

	templater(w, page{
		Title:    "Scoreboard",
		Menu:     urlEvents,
		MenuID:   eventID,
		Heading:  event.Name,
		template: templateScoreboard,
		Data: map[string]interface{}{
			"Event":       event,
			"Ranges":      ranges,
			"SortByRange": rangeId,
			"Colspan":     5 + len(ranges),
		},
	})
}

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
