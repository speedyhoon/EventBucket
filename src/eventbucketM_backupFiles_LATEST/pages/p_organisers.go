package main

import (
	"net/http"
)

func organisers(w http.ResponseWriter, r *http.Request) {
	html := loadHTM("organisers")
	template_data := organisers_Data()
	templator("admin", html, template_data, w)
}

func organisers_Data() map[string]interface{} {
	clubs := getClubs()
	return map[string]interface{}{
		"Title":  "Organisers",
		"Events": generateForm("eventInsert", organisers_eventForm(clubs)),
		"EventList": eventList(),
		"Clubs": generateForm("clubInsert", organisers_clubForm()),
		"ClubList": clubs,
		"Championship": generateForm("champInsert", organisers_champForm()),
	}
}

func organisers_clubForm() map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Club Name",
		},
		"submit": {
			Html:  "submit",
			Label: "Add Club",
		},
	}
}

func organisers_champForm() map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Championship Name",
		},
		"submit": {
			Html:  "submit",
			Label: "Add Championship",
		},
	}
}

func organisers_eventForm(clubs []Club) map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Event Name",
			Required: true,
		},
		"club": {
			Html:         "select",
			SelectValues: getClubSelectionBox(clubs),
			Label:        "Host Club",
			Placeholder:	"Host Club",
			Required: true,
		},
		"submit": {
			Html:  "submit",
			Label: "Add Event",
		},
	}
}

func clubInsert(w http.ResponseWriter, r *http.Request) {
	//TODO it would be nice to pass into the form functions an instance of the struct object and in each form element return the struct property associated with it.
	validated_values := check_form(organisers_clubForm(), r)
	var newClub Club
	newClub.Name = validated_values["name"]
	newClub.Id = getNextId(TBLclub)
	InsertDoc2(TBLclub, newClub)
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_eventForm(getClubs()), r)
	var newEvent Event
	newEvent.Name = validated_values["name"]
	newEvent.Club = validated_values["club"]
	newEvent.Id = getNextId(TBLevent)
	InsertDoc2(TBLevent, newEvent)
}

func InsertDoc2(collection_name string, data interface{}){
	err := conn.C(collection_name).Insert(data)
	checkErr(err)
}

func getClubSelectionBox(club_list []Club) map[string]string {
	drop_down := make(map[string]string)
	for _, club := range club_list {
		drop_down[club.Id] = club.Name
	}
	return drop_down
}
