package main

import (
	"net/http"
)

func organisers(w http.ResponseWriter, r *http.Request) {
	templator(TEMPLATE_ADMIN, "organisers", organisers_Data(), w)
}

func organisers_Data() map[string]interface{} {
	clubs := getClubs()
	return map[string]interface{}{
		"Title":        "Organisers",
		"Events":       generateForm2(organisers_eventForm(clubs)),
		"EventList":    eventList(),
		"Clubs":        generateForm2(organisers_clubForm()),
		"ClubList":     clubs,
		"Championship": generateForm2(organisers_champForm()),
		"Menu":         standard_menu(ORGANISERS_MENU_ITEMS),
		"ShooterList":  generateForm2(organisers_update_shooter_list("")),
	}
}

func organisers_clubForm() Form {
	//TODO add validation to
	return Form{
		Action: URL_clubInsert,
		Title:  "Create Club",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Club Name",
				Required: true,
			},
			"submit": {
				Html:  "submit",
				Label: "Add Club",
			},
		},
	}
}

func organisers_eventForm(clubs []Club) Form {
	club_name := "club"
	club := Inputs{
		Label:    "Host Club",
		Required: true,
	}
	if len(clubs) > 0 {
		club.Html = "select"
		club.SelectValues = getClubSelectionBox(clubs)
//		club.Placeholder = "Select Club
//		club.Html = "select"
//		club.SelectValues = getClubSelectionBox(clubs)
		if len(clubs) > 1 {
			club.Placeholder = "Select Club"
		}
	} else {
		club_name = "club_insert"
		club.Html = "text"
		club.Label = "Host Club"
		club.Placeholder = "Club Name"
	}
	return Form{
		Action: URL_eventInsert,
		Title:  "Create Event",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Event Name",
				Required: true,
			},

			club_name: club,

			"submit": {
				Html:  "submit",
				Label: "Add Event",
			},
		},
	}
}

func organisers_update_shooter_list(last_updated string) Form {
	if last_updated == "" {
		last_updated = "Never"
	}
	return Form{
		Action: URL_updateShooterList,
		Title:  "Update Shooter List",
		Inputs: map[string]Inputs{
			"submit": {
				Html:  "submit",
				Label: "Last updated: " + last_updated,
				Value: "Update",
			},
		},
	}
}

func clubInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_clubForm().Inputs, r)
	insert_new_club(validated_values["name"])
}

func insert_new_club(club_name string) string {
	var newClub Club
	newClub.Name = club_name
	newClub.Id = getNextId(TBLclub)
	newClub.AutoInc.Mound = 1
	InsertDoc(TBLclub, newClub)
	return newClub.Id
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_eventForm(getClubs()).Inputs, r)
	//	export(validated_values)
	var newEvent Event
	newEvent.Name = validated_values["name"]
	if club_name, ok := validated_values["club_insert"]; ok {
		insert_new_club(club_name)
	} else if club_name, ok := validated_values["club"]; ok {
		newEvent.Club = club_name
	}
	newEvent = default_event_settings(newEvent)
	newEvent.Id = getNextId(TBLevent)
	InsertDoc(TBLevent, newEvent)
}

func default_event_settings(event Event) Event {
	//TODO add club settings for default ranges and aggs to create
	event.Ranges = map[string]Range{
		"1": Range{Name: "New Range 1"},
		"2": Range{Name: "New Range 2"},
		"3": Range{Name: "New Aggregate 1", Aggregate: "1,2"},
	}
	event.SortScoreboard = "3"
	event.AutoInc.Range = 4
	return event
}

func getClubSelectionBox(club_list []Club) map[string]string {
	drop_down := make(map[string]string)
	for _, club := range club_list {
		drop_down[club.Id] = club.Name
	}
	return drop_down
}

func eventList() []Club {
	events := getEvents()
	event_list := []Club{}
	for _, row := range events {
		event_list = append(event_list, Club{
			Name: row.Name,
			Url:  "/event/" + row.Id,
		})
	}
	return event_list
}

func organisers_champForm() Form {
	return Form{
		Action: URL_champInsert,
		Title:  "Create Championship",
		Inputs: map[string]Inputs{
			"name": {
				Html:     "text",
				Label:    "Championship Name",
				Required: true,
			},
			"submit": {
				Html:  "submit",
				Label: "Add Championship",
			},
		},
	}
}
