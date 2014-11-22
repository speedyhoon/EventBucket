package main

import (
	"net/http"
	"time"
	"sort"
//	"fmt"
	"strings"
)

func home()M{
	//Sort the list of shooters by grade only
	sort_by_date := func(c1, c2 *Event) bool {
		return c1.Date < c2.Date
	}
	sort_by_time := func(c1, c2 *Event) bool {
		return c1.Time < c2.Time
	}
	sort_by_name := func(c1, c2 *Event) bool {
		return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
	}

	//TODO make custom select
	events := getEvents()
	OrderedByEvent(sort_by_date, sort_by_time, sort_by_name).Sort(events)

//	closed_events := []HomeCalendar{}
	openEvents := []HomeCalendar{}
//	draftEvents := []HomeCalendar{}
//	currentTime := time.Now()

	for _, event := range events {
		if !event.Closed{
			var list_of_ranges []string
			for _, rangeObj := range event.Ranges{
				list_of_ranges = append(list_of_ranges, rangeObj.Name)
			}
			calendar_event := HomeCalendar{
				Id:     event.Id,
				Name:   event.Name,
				ClubId: event.Club,
				Club:   getClub(event.Club).Name,
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
//				}else {
//					fmt.Printf("Event %v doesn't have a valid date", event.Name)
				}
//				if currentTime.After(date_obj){
//					closed_events = append([]HomeCalendar{calendar_event}, closed_events...)
//				}else{

//				}
//			}else{
//				draftEvents = append(draftEvents, calendar_event)
//				fmt.Printf("Event %v doesn't have a date", event.Name)
			}
			openEvents = append(openEvents, calendar_event)
		}
	}

	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return M{
//		"ClosedEvents":   closed_events,
		"FutureEvents":   openEvents,
//		"DraftEvents":   draftEvents,
		"PageName": "Calendar",
		"ArchiveLink": URL_archive,
		"Menu":     home_menu("/", HOME_MENU_ITEMS),
		"FormNewEvent": generateForm2(home_form_new_event(clubs, "","","","", true)),
	}
}

func home_form_new_event(clubs []Club, name, club, date, eventTime string, newEvent bool) Form {
	title := "Event Details"
	save := "Update Event"
	action := URL_eventSettings
	if newEvent {
		title = "Create Event"
		save = "Save Event"
		action = URL_eventInsert
		date = time.Now().Format("2006-01-02")
		eventTime = time.Now().Format("15:04")
	}

	var clubName string

	var club_list []Option
	for _, club_data := range clubs {
		if club != "" && club_data.Id == club {
			clubName = club_data.Name
		}
		club_list = append(club_list, Option{
			Value: club_data.Id,
			Display: club_data.Name,
		})
	}

	return Form{
		Action: action,
		Title:  title,
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Event Name",
				Required: true,
//				AutoComplete: "off",
				Value: name,
			},
			{
				Name: "club",
				Html: "datalist",
				Label: "Host Club",
				Placeholder: "Club Name",
				Options: club_list,
				Required: true,
//				AutoComplete: "off",
				Value: clubName,
			},
			{
				Name: "date",
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			{
				Name: "time",
				Html:  "time",
				Label: "Time",
				Value: eventTime,
			},
			{
				Html:  "submit",
				Value: save,
			},
		},
	}
}

func eventInsert2(w http.ResponseWriter, r *http.Request) {
	var clubs []Club
	validated_values := check_form(home_form_new_event(clubs,"","","","",true).Inputs, r)

//	export(validated_values)

	newEvent := Event{
		Name: validated_values["name"],
	}

//	var newEvent Event
//	newEvent.Name = validated_values["name"]

	club, ok := getClub_by_name(validated_values["club"])
	if ok {
		newEvent.Club = club.Id
	}else{
		newEvent.Club = insert_new_club(validated_values["club"])
	}

	if validated_values["date"] != ""{
		newEvent.Date = validated_values["date"]
	}

	if validated_values["time"] != ""{
		newEvent.Time = validated_values["time"]
	}

	//Add default ranges and aggregate ranges
	newEvent.Id = getNextId(TBLevent)
	newEvent.AutoInc.Range = 1
	InsertDoc(TBLevent, newEvent)

	//redirect user to event settings
	redirecter(URL_eventSettings+newEvent.Id, w, r)
}


//Home and archive pages
type HomeCalendar struct {
	Id, Name, Club, ClubId, Time string
	Day                          string
	Date, Ranges                 string
	Month                        time.Month
	Year                         int
}

type lessFunc2 func(p1, p2 *Event) bool

type multiSorter2 struct {
	changes []Event
	less    []lessFunc2
}

func (ms *multiSorter2) Sort(changes []Event) {
	ms.changes = changes
	sort.Sort(ms)
}

func (ms *multiSorter2) Len() int {
	return len(ms.changes)
}

func (ms *multiSorter2) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *multiSorter2) Less(i, j int) bool {
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

func OrderedByEvent(less ...lessFunc2) *multiSorter2 {
	return &multiSorter2{
		less: less,
	}
}
