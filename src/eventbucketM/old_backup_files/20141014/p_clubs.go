package main

//
//import (
//	"net/http"
//	"fmt"
//)
//
//func clubs(w http.ResponseWriter, r *http.Request) {
//	templator(TEMPLATE_HOME, "club", clubs_Data(), w)
//}
//
//func clubs_Data(club_id string) map[string]interface{} {
//	this_club := getClub(club_id)
//	var temporary []string
//	for _, mound := range this_club.Mounds{
//		temporary = append(temporary, generateForm("clubMoundUpdate", clubMoundUpdateForm(this_club.Id, mound)))
//	}
//	return map[string]interface{}{
//		"Title": this_club.Name,
//		"Id": this_club.Id,
//		"ListMounds": temporary,
//		"Latitude": this_club.Latitude,
//		"Longitude": this_club.Longitude,
//		"InsertRangeForm": generateForm("clubMoundInsert", clubMoundInsertForm(this_club.Id)),
//	}
//}
//
//func clubMoundInsertForm(club_id string) map[string]Inputs {
//	return map[string]Inputs{
//		"clubid": {
//			Html:  "hidden",
//			Value: club_id,
//		},
//		"distance": {
//			Html:	"number",
//			Label: "Distance",
//			Required: true,
//			Min: 1,
//		},
//		"unit": {
//			Html: "select",
//			Required: true,
//			Label: "Unit",
//			Select: []string{"Yards","Metres"},
//		},
//		"submit": {
//			Html:  "submit",
//			Label: "Insert New Mound",
//		},
//	}
//}
//
//func clubMoundUpdateForm(club_id string, mound Mound) map[string]Inputs {
//	return map[string]Inputs{
//		"clubid": {
//			Html:  "hidden",
//			Value: club_id,
//		},
//		"distance": {
//			Html:	"number",
//			Label: "Distance",
//			Required: true,
//			Value: echo(mound.Distance),
//		},
//		"unit": {
//			Html: "select",
//			Required: true,
//			Label: "Unit",
//			Select: []string{"Yards","Metres"},
//		},
//		"submit": {
//			Html:  "submit",
//			Label: "Update Mound",
//		},
//	}
//}
//
//func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
//	validated_values := check_form(clubMoundInsertForm(""), r)
//	club_id := validated_values["clubid"]
//	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
//	var new_mound Mound
//	new_mound.Distance = str_to_int(validated_values["distance"])
//	new_mound.Unit = validated_values["unit"]
//	club_insert_mound(club_id, new_mound)
//}
