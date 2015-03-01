package main

import (
	"time"
	"strings"
)

func archive()Page{
	//Sort the list of shooters by grade only
	sort_by_date := func(c1, c2 *Event) bool {
		return c1.Date > c2.Date
	}
	sort_by_time := func(c1, c2 *Event) bool {
		return c1.Time > c2.Time
	}
	sort_by_name := func(c1, c2 *Event) bool {
		return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
	}
	//TODO make custom mongodb query to get a more flexible list of events?
	events := getEvents()
	OrderedByEvent(sort_by_date, sort_by_time, sort_by_name).Sort(events)
	closedEvents := []HomeCalendar{}
	for _, event := range events {
		if event.Closed{
			var list_of_ranges []string
			for _, rangeObj := range event.Ranges{
				list_of_ranges = append(list_of_ranges, rangeObj.Name)
			}
			club, _ := getClub(event.Club)
			calendar_event := HomeCalendar{
				Id:     event.Id,
				Name:   event.Name,
				ClubId: event.Club,
				Club:   club.Name,
				Time:   event.Time,
				Ranges: strings.Join(list_of_ranges, ", "),
			}
			if event.Date != "" {
				date_obj, err := time.Parse("2006-01-02", event.Date)
				if err == nil {
					calendar_event.Day = date_obj.Weekday().String()
					calendar_event.Date = ordinal(date_obj.Day())
					calendar_event.Month = date_obj.Month()
					calendar_event.Year = date_obj.Year()
				}
			}
			closedEvents = append(closedEvents, calendar_event)
		}
	}
	return Page {
		TemplateFile: "archive",
		Theme: TEMPLATE_HOME,
		Data: M{
			"ClosedEvents":   closedEvents,
			"PageName": "Calendar",
			"Menu":     home_menu(URL_archive, HOME_MENU_ITEMS),
		},
	}
}
