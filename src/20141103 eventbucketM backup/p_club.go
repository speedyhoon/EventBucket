package main

import (
	"net/http"
	//	"fmt"
)

func club(w http.ResponseWriter, r *http.Request) {
	club_id := get_id_from_url(r, URL_club)
	//TODO change club url to be the club.Url instead of club.Id
	templator(TEMPLATE_HOME, "club", club_Data(club_id), w)
}

func club_Data(club_id string) map[string]interface{} {
	this_club := getClub(club_id)
	//	menu_items := []Menu{
	//		Menu{
	//			Name:		"Home",
	//			Link:		"/",
	//		},
	//		Menu{
	//			Name:		"Clubs",
	//			Link:		"/clubs",
	//		},
	//		Menu{
	//			Name:		"Events",
	//			Link:		"/events",
	//		},
	//		Menu{
	//			Name:		"Event Archive",
	//			Link:		"/archive",
	//		},
	//		Menu{
	//			Name:		"Organisers",
	//			Link:		"/organisers",
	//		},
	//	}
	var temporary []string
	for _, mound := range this_club.Mounds {
		temporary = append(temporary, generateForm2(clubMoundUpdateForm(this_club.Id, mound)))
	}
	return map[string]interface{}{
		"Name":       this_club.Name,
		"Menu":       home_menu("/clubs",HOME_MENU_ITEMS),
		"ClubId":     this_club.Id,
		"ListMounds": this_club.Mounds,
		"Latitude":   this_club.Latitude,
		"Longitude":  this_club.Longitude,
	}
}

//func clubMoundInsertForm(club_id string) Form{
//	return Form{
//		Action: URL_club_mound_insert,
//		Title: "foobar",
//		Inputs: map[string]Inputs{
//			"clubid": {
//				Html:  "hidden",
//				Value: club_id,
//			},
//			"distance": {
//				Html:   "number",
//				Label: "Distance",
//				Required: true,
//				Min: 1,
//			},
//			"unit": {
//				Html: "select",
//				Required: true,
//				Label: "Unit",
//				Select: []string{"Yards", "Metres"},
//			},
//			"submit": {
//				Html:  "submit",
//				Value: "Insert New Mound",
//			},
//		},
//	}
//}
//
//func clubMoundUpdateForm(club_id string, mound Mound) Form {
//	return Form{
//		Action: URL_club_mound_update,
//		Title: "ellh!!",
//		Inputs: map[string]Inputs{
//			"clubid": {
//				Html:  "hidden",
//				Value: club_id,
//			},
//			"distance": {
//				Html:   "number",
//				Label: "Distance",
//				Required: true,
//				Value: echo(mound.Distance),
//			},
//			"unit": {
//				Html: "select",
//				Required: true,
//				Label: "Unit",
//				Select: []string{"Yards", "Metres"},
//			},
//			"submit": {
//				Html:  "submit",
//				Value: "Update Mound",
//			},
//		},
//	}
//}

//func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
//	validated_values := check_form(clubMoundInsertForm(""), r)
//	club_id := validated_values["clubid"]
//	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
//	var new_mound Mound
//	new_mound.Distance = str_to_int(validated_values["distance"])
//	new_mound.Unit = validated_values["unit"]
//	club_insert_mound(club_id, new_mound)
//}
