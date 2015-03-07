package main

import (
	"net/http"
	"sort"
	"strings"
	"time"
)

func home() Page {
	//Sort the list of shooters by grade only
	sortByDate := func(c1, c2 *Event) bool {
		return c1.Date < c2.Date
	}
	sortByTime := func(c1, c2 *Event) bool {
		return c1.Time < c2.Time
	}
	sortByName := func(c1, c2 *Event) bool {
		return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
	}
	events := getEvents() //TODO make custom select
	OrderedByEvent(sortByDate, sortByTime, sortByName).Sort(events)
	openEvents := []HomeCalendar{}
	for _, event := range events {
		if !event.Closed {
			var listRanges []string
			for _, rangeObj := range event.Ranges {
				listRanges = append(listRanges, rangeObj.Name)
			}
			club, _ := getClub(event.Club)
			calendarEvent := HomeCalendar{
				Id:     event.Id,
				Name:   event.Name,
				ClubId: event.Club,
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
			openEvents = append(openEvents, calendarEvent)
		}
	}
	hostname, ipAddresses := HostnameIpAddresses()
	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return Page{
		TemplateFile: "home",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"FutureEvents": openEvents,
			"PageName":     "Calendar",
			"ArchiveLink":  URL_archive,
			"Menu":         homeMenu("/", HOME_MENU_ITEMS),
			"FormNewEvent": generateForm(homeFormNewEvent(clubs, (Event{}))),
			"Hostname":     hostname,
			"IpAddresses":  ipAddresses,
		},
		v8Url: VURL_home,
	}
}

func homeFormNewEvent(clubs []Club, event Event) Form {
	title := "Event Details"
	save := "Update Event"
	if event.Id == "" {
		title = "New Event"
		save = "Save Event"
		event.Date = time.Now().Format("2006-01-02")
		event.Time = time.Now().Format("15:04")
	}
	var clubList []Option
	for _, clubData := range clubs {
		if event.Club != "" && clubData.Id == event.Club {
			event.Club = clubData.Name
		}
		clubList = append(clubList, Option{
			Value:   clubData.Id,
			Display: clubData.Name,
		})
	}
	inputs := []Inputs{
		{
			Name:     "name",
			Html:     "text",
			Label:    "Event Name",
			Required: true,
			//				AutoComplete: "off",
			Value:     event.Name,
			Autofocus: "on",
		}, {
			Name:        "club",
			Html:        "datalist",
			Label:       "Host Club",
			Placeholder: "Club Name",
			//TODO previous club names appear from browser cahce when they are not available
			//TODO auto set the club name to X if there is only one available
			Options:  clubList,
			Required: true,
			//				AutoComplete: "off",
			Value: event.Club,
		}, {
			Name:     "date",
			Html:     "date",
			Label:    "Date",
			Required: true,
			Value:    event.Date,
		}, {
			Name:  "time",
			Html:  "time",
			Label: "Time",
			Value: event.Time,
		}, {
			Html:  "submit",
			Value: save,
		},
	}
	if event.Id != "" {
		inputs = append(inputs, Inputs{
			Name:  "eventid",
			Html:  "hidden",
			Value: event.Id,
		})
	}

	return Form{
		Action: URL_eventInsert,
		Title:  title,
		Inputs: inputs,
	}
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	//TODO merge this database functionality into an upsert
	var clubs []Club
	emptyEvent := Event{Id: "1"}
	validatedValues := checkForm(homeFormNewEvent(clubs, emptyEvent).Inputs, r)
	//TODO one day change validated values to return the schema compatible data so it can be directly used add constants would by used more often to access the map items
	eventId := validatedValues["eventid"]
	if eventId == "" {
		newEvent := Event{
			Name: validatedValues["name"],
		}
		club, ok := getClubByName(validatedValues["club"])
		if ok {
			newEvent.Club = club.Id
		} else {
			clubId, _ := insertClub(validatedValues["club"])
			newEvent.Club = clubId
		}

		if validatedValues["date"] != "" {
			newEvent.Date = validatedValues["date"]
		}

		if validatedValues["time"] != "" {
			newEvent.Time = validatedValues["time"]
		}

		//Add default ranges and aggregate ranges
		var err error
		newEvent.Id, err = getNextId(TBLevent)
		newEvent.AutoInc.Range = 1
		if err == nil {
			InsertDoc(TBLevent, newEvent)
			//redirect user to event settings
			http.Redirect(w, r, URL_eventSettings+newEvent.Id, http.StatusSeeOther)
		} else {
			//TODO go to previous referer page (home or organisers both have the form)
			//http.Redirect(w, r, URL_organisers, http.StatusSeeOther)
		}
	} else {
		eventUpsertData(eventId, M{
			"n": validatedValues["name"],
			"c": validatedValues["club"],
			"d": validatedValues["date"],
			"t": validatedValues["time"],
		})
		http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
	}
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
