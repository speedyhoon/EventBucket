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
	orderedByEvent(sortByDate, sortByTime, sortByName).Sort(events)
	openEvents := []homeCalendar{}
	for _, event := range events {
		if !event.Closed {
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
			openEvents = append(openEvents, calendarEvent)
		}
	}
	hostname, ipAddresses := hostnameIPAddresses()
	//TODO change getClubs to simpler DB lookup getClubNames
	clubs := getClubs()
	return Page{
		TemplateFile: "home",
		Theme:        templateHome,
		Title:        "Home",
		Data: M{
			"FutureEvents": openEvents,
			"PageName":     "Calendar",
			"ArchiveLink":  urlArchive,
			"FormNewEvent": generateForm(homeFormNewEvent(clubs, (Event{}))),
			"Hostname":     hostname,
			"IpAddresses":  ipAddresses,
		},
		v8Url: vURLHome,
	}
}

func homeFormNewEvent(clubs []Club, event Event) Form {
	title := "Event Details"
	save := "Update Event"
	var nameFocus bool
	submitName := "eventid"
	if event.ID == "" {
		title = "New Event"
		save = "Save Event"
		nameFocus = true
		submitName = ""
		event.Date = time.Now().Format("2006-01-02")
		event.Time = time.Now().Format("15:04")
	}
	var clubList []Option
	for _, clubData := range clubs {
		if event.Club != "" && clubData.ID == event.Club {
			event.Club = clubData.Name
		}
		clubList = append(clubList, Option{
			Value: clubData.Name,
		})
	}
	return Form{
		action: urlEventInsert,
		title:  title,
		inputs: []Inputs{
			{
				name:      "name",
				html:      "text",
				label:     "Event Name",
				varType:   "string",
				maxLength: v8MaxStringInput,
				minLength: v8MinStringInput,
				required:  true,
				value:     event.Name,
				autofocus: nameFocus,
			}, {
				name:         "club",
				html:         "search",
				dataList:     true,
				id:           "clubSearch",
				autoComplete: "off",
				label:        "Club Name",
				//TODO auto set the club name to X if there is only one available
				options:   clubList,
				required:  true,
				value:     event.Club,
				varType:   "string",
				maxLength: v8MaxStringInput,
				minLength: v8MinStringInput,
			}, {
				name:     "date",
				html:     "date",
				label:    "Date",
				required: true,
				value:    event.Date,
			}, {
				name:  "time",
				html:  "time",
				label: "Time",
				value: event.Time,
			}, {
				html:      "submit",
				inner:     save,
				name:      submitName,
				value:     event.ID,
				varType:   "string",
				maxLength: v8MaxEventID,
				minLength: v8MinEventID,
			},
		},
	}
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	/*validatedValues, ok := v8(homeFormNewEvent([]Club{}, Event{}).inputs, r)
	if !ok {
		http.Error(w, "Invalid request form data", 400)
		return
	}*/
	//TODO merge this database functionality into an upsert
	validatedValues := checkForm(homeFormNewEvent([]Club{}, Event{}).inputs, r)
	//TODO one day change validated values to return the schema compatible data so it can be directly used add constants would by used more often to access the map items
	eventID := validatedValues["eventid"]
	if eventID == "" {
		newEvent := Event{
			Name: validatedValues["name"],
		}
		club, ok := getClubByName(validatedValues["club"])
		if ok {
			newEvent.Club = club.ID
		} else {
			clubID, isClub := insertClub(validatedValues["club"])
			if isClub != nil {
				newEvent.Club = clubID
			}
		}

		if validatedValues["date"] != "" {
			newEvent.Date = validatedValues["date"]
		}

		if validatedValues["time"] != "" {
			newEvent.Time = validatedValues["time"]
		}

		//Add default ranges and aggregate ranges
		var err error
		newEvent.ID, err = getNextID(tblEvent)
		newEvent.AutoInc.Range = 1
		if err == nil {
			insertDoc(tblEvent, newEvent)
			//redirect user to event settings
			http.Redirect(w, r, urlEventSettings+newEvent.ID, http.StatusSeeOther)
		} else {
			//TODO go to previous referer page (home or organisers both have the form)
			//http.Redirect(w, r, urlOrganisers, http.StatusSeeOther)
		}
	} else {
		tableUpdateData(tblEvent, eventID, M{
			"n": validatedValues["name"],
			"c": validatedValues["club"],
			"d": validatedValues["date"],
			"t": validatedValues["time"],
		})
		http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
	}
}

//Home and archive pages
type homeCalendar struct {
	ID, Name, Club, ClubID, Time string
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

func orderedByEvent(less ...lessFunc2) *multiSorter2 {
	return &multiSorter2{
		less: less,
	}
}
