package main

import (
	"net/http"
	"time"
	"sort"
//	"fmt"
	"strings"
)

type HomeCalendar struct {
	Id, Name, Club, ClubId, Time string
	Day                          string
	Date, Ranges                 string
	Month                        time.Month
	Year                         int
}

func home(w http.ResponseWriter, r *http.Request) {
	//TODO change sending a string filename to the URL_page and devmode handles it automatically
	templator(TEMPLATE_HOME, "home", homeData(getEvents()), w)
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




func homeData(events []Event) map[string]interface{} {
	closed_events := []HomeCalendar{}
	future_events := []HomeCalendar{}
	draft_events := []HomeCalendar{}


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
	OrderedByEvent(sort_by_date, sort_by_time, sort_by_name).Sort(events)



	currentTime := time.Now()

	for _, event := range events {
		calendar_event := HomeCalendar{
			Id:     event.Id,
			Name:   event.Name,
			ClubId: event.Club,
			Club:   getClub(event.Club).Name,
			Time:   event.Time,
		}
		var list_of_ranges []string
		for _, rangeObj := range event.Ranges{
			list_of_ranges = append(list_of_ranges, rangeObj.Name)
		}
		calendar_event.Ranges = strings.Join(list_of_ranges, ", ")
		if event.Date != "" {
			//			export( event.Date)
			date_obj, err := time.Parse("2006-01-02", event.Date)
			checkErr(err)
			calendar_event.Day = date_obj.Weekday().String()
			calendar_event.Date = ordinal(date_obj.Day())
			calendar_event.Month = date_obj.Month()
			calendar_event.Year = date_obj.Year()
			if currentTime.After(date_obj){
				closed_events = append([]HomeCalendar{calendar_event}, closed_events...)
			}else{
				future_events = append(future_events, calendar_event)
			}
		}else{
			draft_events = append(draft_events, calendar_event)
		}
	}

	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return map[string]interface{}{
		"ClosedEvents":   closed_events,
		"FutureEvents":   future_events,
		"DraftEvents":   draft_events,
		"PageName": "Calendar",
		"Menu":     home_menu("/", HOME_MENU_ITEMS),
		"FormNewEvent": generateForm2(home_form_new_event(clubs, "","","","")),
	}
}

func home_form_new_event(clubs []Club, name, club, date, eventTime string) Form {
	var action, title, save string
	if name != "" || club != "" || date != "" || eventTime != ""{
		title = "Event Details"
		save = "Update Event"
		//TODO change update to a new save function
		action = URL_eventInsert2
	}else {
		action = URL_eventInsert2
		title = "Create Event"
		save = "Save Event"
		date = time.Now().Format("2006-01-02")
		eventTime = time.Now().Format("15:04")
	}

	var clubName string

	var club_list []Option
	for _, club_data := range clubs {
		if club_data.Id == club {
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
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Event Name",
				Required: true,
//				AutoComplete: "off",
				Value: name,
			},
			"club": {
				Html: "datalist",
				Label: "Host Club",
				Placeholder: "Club Name",
				Options: club_list,
				Required: true,
				AutoComplete: "off",
				Value: clubName,
			},
			"date": {
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			"time": {
				Html:  "time",
				Label: "Time",
				Value: eventTime,
			},
			"submit": {
				Html:  "submit",
				Value: save,
			},
		},
	}
}


func eventInsert2(w http.ResponseWriter, r *http.Request) {
	var clubs []Club
	validated_values := check_form(home_form_new_event(clubs,"","","","").Inputs, r)

	var newEvent Event
	newEvent.Name = validated_values["name"]

	club_name := validated_values["club"]
	club, ok := getClub_by_name(club_name)
	if ok {
		newEvent.Club = club.Id
	}else{
		newEvent.Club = insert_new_club(club_name)
	}

	if validated_values["date"] != ""{
		newEvent.Date = validated_values["date"]
	}

	if validated_values["time"] != ""{
		newEvent.Time = validated_values["time"]
	}

	//Add default ranges and aggregate ranges
	newEvent = default_event_settings(newEvent)

	newEvent.Id = getNextId(TBLevent)
	InsertDoc(TBLevent, newEvent)

	//redirect user to event settings
	redirecter(URL_eventSettings+newEvent.Id, w, r)
}
