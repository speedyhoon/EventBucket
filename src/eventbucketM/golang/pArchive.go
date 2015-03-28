package main

import (
	"strings"
	"time"
)

func archive() Page {
	//Sort the list of shooters by grade only
	sortByDate := func(c1, c2 *Event) bool {
		return c1.Date > c2.Date
	}
	sortByTime := func(c1, c2 *Event) bool {
		return c1.Time > c2.Time
	}
	sortByName := func(c1, c2 *Event) bool {
		return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
	}
	//TODO make custom mongodb query to get a more flexible list of events?
	events := getEvents()
	orderedByEvent(sortByDate, sortByTime, sortByName).Sort(events)
	closedEvents := []homeCalendar{}
	for _, event := range events {
		if event.Closed {
			var listRanges []string
			for _, rangeObj := range event.Ranges {
				listRanges = append(listRanges, rangeObj.Name)
			}
			club, _ := getClub(event.Club)
			calendarEvent := homeCalendar{
				ID:     event.ID,
				Name:   event.Name,
				ClubID: event.Club,
				Club:   club.Name,
				Time:   event.Time,
				Ranges: strings.Join(listRanges, ", "),
			}
			if event.Date != "" {
				dateObj, err := time.Parse("2006-01-02", event.Date)
				if err == nil {
					calendarEvent.Day = dateObj.Weekday().String()
					calendarEvent.Date = ordinal(dateObj.Day())
					calendarEvent.Month = dateObj.Month()
					calendarEvent.Year = dateObj.Year()
				}
			}
			closedEvents = append(closedEvents, calendarEvent)
		}
	}
	return Page{
		TemplateFile: "archive",
		Title:        "Archive",
		Theme:        templateHome,
		Data: M{
			"ClosedEvents": closedEvents,
			"PageName":     "Calendar",
		},
	}
}
