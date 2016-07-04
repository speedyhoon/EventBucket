package main

import (
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
	var ranges []Range
	if err == nil {
		ranges = findAggs(rangeID, event.Ranges)
	}
	if len(ranges) < 1 {
		errorHandler(w, r, http.StatusNotFound, "range")
		return
	}

	sortShooters(ids[1]).Sort(event.Shooters)

	templater(w, page{
		Title:    "Scoreboard",
		Menu:     urlEvents,
		MenuID:   ids[0],
		Heading:  event.Name,
		template: templateScoreboard,
		Data: map[string]interface{}{
			"Event":       event,
			"Ranges":      ranges,
			"SortByRange": ids[1],
			"Colspan":     5 + len(ranges),
		},
	})
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
