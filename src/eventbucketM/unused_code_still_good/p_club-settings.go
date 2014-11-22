package main

import (
	"fmt"
	"net/http"
)

func club_settings(w http.ResponseWriter, r *http.Request) {
//	club_id := get_id_from_url(r, URL_club)
//	templator(TEMPLATE_ADMIN, "club-settings", club_settings_Data(club_id), w)
}

func club_settings_Data(club_id string) map[string]interface{} {
	this_club := getClub(club_id)
	var temporary []string
	for _, mound := range this_club.Mounds {
		temporary = append(temporary, generateForm2(clubMoundUpdateForm(this_club.Id, mound)))
	}
	return map[string]interface{}{
		"Name":            this_club.Name,
		"Id":              this_club.Id,
		"ListMounds":      temporary,
		"Latitude":        this_club.Latitude,
		"Longitude":       this_club.Longitude,
		"InsertRangeForm": generateForm2(clubMoundInsertForm(this_club.Id)),
		"MapForm":         generateForm2(clubMapUpsertForm(this_club.Id)),
	}
}

//func clubs(w http.ResponseWriter, r *http.Request) {
//	templator(TEMPLATE_HOME, "clubs", clubs_Data(), w)
//}

func clubMoundInsertForm(club_id string) Form {
	return Form{
		Action: "clubMoundInsert",
		Title:  "Insert Mound",
		Inputs: []Inputs{{
			Name: "clubid",
				Html:  "hidden",
				Value: club_id,
			},{
			Name: "distance",
				Html:     "number",
				Label:    "Distance",
				Required: true,
				Min:      "1",
			},{
			Name: "unit",
				Html:     "select",
				Required: true,
				Label:    "Unit",
				Options:   []Option{Option{Value:"Yards",Display:"Yards"},Option{Value:"Metres",Display:"Metres"}},
			},{
				Html:  "submit",
				Value: "Insert New Mound",
			},
		},
	}
}

//if inputData.Min != 0 {
//output += fmt.Sprintf(" min=%f",inputData.Min)
//}
//if inputData.Max != 0{
//output += fmt.Sprintf(" max=%f",inputData.Max)
//}
//if inputData.Step != 0{
//output += fmt.Sprintf(" step=%f",inputData.Step)
//}

func clubMapUpsertForm(club_id string) Form {
	return Form{
		Action: "clubMapUpsert",
		Title:  "Update Map Location",
		Inputs: []Inputs{{
			Name: "clubid",
				Html:  "hidden",
				Value: club_id,
			},{
			Name: "latitude",
				Html:     "number",
				Label:    "Latitude",
				Required: true,
				Min:      "-90",
				Step:     0.000001,
				Max:      "90",
			},{
				Name: "longitude",
				Html:     "number",
				Required: true,
				Label:    "Longitude",
				Min:      "-180",
				Step:     0.000001,
				Max:      "180",
			},
			{
				Html:  "submit",
				Value: "Update Co-ordinates",
			},
		},
	}
}

func clubMoundUpdateForm(club_id string, mound Mound) Form {
	return Form{
		Action: "clubMoundUpdate",
		Inputs: []Inputs{{
			Name: "clubid",
				Html:  "hidden",
				Value: club_id,
			},{
				Name: "distance",
				Html:     "number",
				Label:    "Distance",
				Required: true,
				Value:    echo(mound.Distance),
			},{
			Name: "unit",
				Html:     "select",
				Required: true,
				Label:    "Unit",
				Options:   []Option{Option{Value:"Yards", Display:"Yards"},Option{Value:"Metres", Display:"Metres"}},
			},{
				Html:  "submit",
				Value: "Update Mound",
			},
		},
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(clubMoundInsertForm("").Inputs, r)
	club_id := validated_values["clubid"]
	redirecter(fmt.Sprintf("/club/%v", club_id), w, r)
	var new_mound Mound
	new_mound.Distance = str_to_int(validated_values["distance"])
	new_mound.Unit = validated_values["unit"]
	//	club_insert_mound(club_id, new_mound)
}
