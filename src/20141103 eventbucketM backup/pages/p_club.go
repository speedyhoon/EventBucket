package main

import (
	"net/http"
	"fmt"
)

func club(w http.ResponseWriter, r *http.Request) {
//	page_url := len("/club/")
	//	var validID = regexp.MustCompile(`\A`+page_url+`[0-9a-f]{24}\z`)
	url := echo(r.URL)
	//	if validID.MatchString(url){
	templator("admin", loadHTM("club"), club_Data(url[6:]), w)
	//	}else{
	//		redirectPermanent("/clubs")
	//		fmt.Println("redirected user "+url)
	//	}
}

func club_Data(club_id string) map[string]interface{} {
	this_club := getClub(club_id)
	var temporary []string
	for _, mound := range this_club.Mounds{
		temporary = append(temporary, generateForm("clubMoundUpdate", clubMoundUpdateForm(this_club.Id, mound)))
	}
	return map[string]interface{}{
		"Title": this_club.Name,
		"Id": this_club.Id,
		"ListMounds": temporary,
		"Latitude": this_club.Latitude,
		"Longitude": this_club.Longitude,
		"InsertRangeForm": generateForm("clubMoundInsert", clubMoundInsertForm(this_club.Id)),
	}
}

func clubMoundInsertForm(club_id string) map[string]Inputs {
	return map[string]Inputs{
		"clubid": {
			Html:  "hidden",
			Value: club_id,
		},
		"distance": {
			Html:	"number",
			Label: "Distance",
			Required: true,
			Min: 1,
		},
		"unit": {
			Html: "select",
			Required: true,
			Label: "Unit",
			Select: []string{"Yards","Metres"},
		},
		"submit": {
			Html:  "submit",
			Label: "Insert New Mound",
		},
	}
}

func clubMoundUpdateForm(club_id string, mound Mound) map[string]Inputs {
	return map[string]Inputs{
		"clubid": {
			Html:  "hidden",
			Value: club_id,
		},
		"distance": {
			Html:	"number",
			Label: "Distance",
			Required: true,
			Value: echo(mound.Distance),
		},
		"unit": {
			Html: "select",
			Required: true,
			Label: "Unit",
			Select: []string{"Yards","Metres"},
		},
		"submit": {
			Html:  "submit",
			Label: "Update Mound",
		},
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(clubMoundInsertForm(""), r)
	club_id := validated_values["clubid"]
	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
	var new_mound Mound
	new_mound.Distance = str_to_int(validated_values["distance"])
	new_mound.Unit = validated_values["unit"]
	club_insert_mound(club_id, new_mound)
}
