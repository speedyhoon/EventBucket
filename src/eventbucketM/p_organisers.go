package main

import (
	"net/http"
)

func organisers() Page {
	clubs := getClubs()
	return Page {
		TemplateFile: "organisers",
		Theme: TEMPLATE_HOME,
		Data: M{
			"Title":        "Organisers",
			//		"Events":       generateForm2(organisers_eventForm(clubs)),
			"Events":       generateForm2(home_form_new_event(clubs, "","","","",true)),
			"EventList":    eventList(),
			"Clubs":        generateForm2(organisers_clubForm()),
			"ClubList":     clubs,
			//		"Championship": generateForm2(organisers_champForm()),
			//		"Menu":         standard_menu(ORGANISERS_MENU_ITEMS),
			"Menu":     home_menu(URL_organisers, HOME_MENU_ITEMS),
			"ShooterList":  generateForm2(organisers_update_shooter_list("")),
		},
	}
}

func organisers_clubForm() Form {
	//TODO add validation to
	return Form{
		Action: URL_clubInsert,
		Title:  "New Club",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Club Name",
				Required: true,
			},
			{
				Html:  "submit",
				Value: "Add Club",
			},
		},
	}
}

/*func organisers_eventForm(clubs []Club) Form {
//	club_name := "club"
	club := Inputs{
		Name: "club",
		Label:    "Host Club",
		Required: true,
	}
	if len(clubs) > 0 {
		club.Html = "select"
		club.Options = getClubSelectionBox(clubs)
//		club.Placeholder = "Select Club
//		club.Html = "select"
//		club.SelectValues = getClubSelectionBox(clubs)
		if len(clubs) > 1 {
			club.Placeholder = "Select Club"
		}
	} else {
		club.Name = "club_insert"
		club.Html = "text"
		club.Label = "Host Club"
		club.Placeholder = "Club Name"
	}
	return Form{
		Action: URL_eventInsert,
		Title:  "Create Event",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Event Name",
				Required: true,
			},

			club,

			{
				Html:  "submit",
				Value: "Add Event",
			},
		},
	}
}*/

func organisers_update_shooter_list(last_updated string) Form {
	if last_updated == "" {
		last_updated = "Never"
	}
	return Form{
		Action: URL_updateShooterList,
		Title:  "Update Shooter List",
		Inputs: []Inputs{
			{
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
	nextId, err := getNextId(TBLclub)
	if err == nil {
		var newClub Club
		newClub.Name = club_name
		newClub.Id = nextId
		newClub.AutoInc.Mound = 1
		InsertDoc(TBLclub, newClub)
		return newClub.Id
	}
	return "-1"
}

func getClubSelectionBox(club_list []Club) []Option {
	var drop_down []Option
	for _, club := range club_list {
		drop_down = append(drop_down, Option{Display: club.Name, Value:club.Id})
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

/*func organisers_champForm() Form {
	return Form{
		Action: URL_champInsert,
		Title:  "Create Championship",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Championship Name",
				Required: true,
			},
			{
				Html:  "submit",
				Value: "Add Championship",
			},
		},
	}
}*/
