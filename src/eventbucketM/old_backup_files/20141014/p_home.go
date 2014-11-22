package main

import (
	"net/http"
	"time"
)

type HomeCalendar struct {
	Id, Name, Club, ClubId, Time string
	Day                          string
	Date                         string
	Month                        time.Month
	Year                         int
}

func home(w http.ResponseWriter, r *http.Request) {
	//TODO change sending a string filename to the URL_page and devmode handles it automatically
	templator(TEMPLATE_HOME, "home", homeData(getEvents()), w)
}

func homeData(events []Event) map[string]interface{} {
	calendar_events := []HomeCalendar{}
	for _, event := range events {
		calendar_event := HomeCalendar{
			Id:     event.Id,
			Name:   event.Name,
			ClubId: event.Club,
			Club:   getClub(event.Club).Name,
			Time:   event.Time,
		}
		if event.Date != "" {
			//			export( event.Date)
			date_obj, err := time.Parse("2006-01-02", event.Date)
			checkErr(err)
			calendar_event.Day = date_obj.Weekday().String()
			calendar_event.Date = ordinal(date_obj.Day())
			calendar_event.Month = date_obj.Month()
			calendar_event.Year = date_obj.Year()
		}
		calendar_events = append(calendar_events, calendar_event)
	}

	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return map[string]interface{}{
		"Events":   calendar_events,
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

	var club_list []string
	for _, club_data := range clubs {
		if club_data.Id == club {
			export(club_data)
			clubName = club_data.Name
		}
		club_list = append(club_list, club_data.Name)
	}

	return Form{
		Action: action,
		Title:  title,
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Event Name",
				Required: true,
				AutoComplete: "off",
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
			"date": Inputs{
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			"time": Inputs{
				Html:  "time",
				Label: "Time",
				Value: eventTime,
			},
			"submit": {
				Html:  "submit",
				Label: save,
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
