package main

import (
	"net/http"
	"fmt"
)

func club(clubId string) Page {
	club, err := getClub(clubId)
	if err != nil{
		//TODO return a 404 error
		return Page {
			TemplateFile: "club",
			Theme: TEMPLATE_HOME,
			Title:  "Club with id '" + clubId + "' not found",
			Data: M{
				"Menu":  home_menu(URL_club, HOME_MENU_ITEMS),
			},
		}
	}
	var mounds []string
	for _, mound := range club.Mounds{
		mounds = append(mounds, generateForm2(clubMoundInsertForm(clubId, mound, true)))
	}
	return Page {
		TemplateFile: "club",
		Theme: TEMPLATE_HOME,
		Title: club.Name,
		Data: M{
			"Menu": home_menu(URL_club, HOME_MENU_ITEMS),
			"MoundForm": generateForm2(clubMoundInsertForm(clubId, Mound{}, false)),
			"MapForm": generateForm2(clubMapUpsertForm(clubId, club.Latitude, club.Longitude)),
			"ListMounds": mounds,
			"Latitude": club.Latitude,
			"Longitude": club.Longitude,
			//TODO integrate with club settings page so it is simpler to use!
		},
	}
}

func clubMoundInsertForm(clubId string, mound Mound, existing bool) Form {
	min := 10
	max := 9750
	submitLabel := "Create Mound"
	if existing {
		submitLabel = "Save Mound"
	}
	options := []Option{
		{
			Display: "Yards",
			Value: "y",
		},
		{
			Display: "Metres",
			Value: "m",
		},
	}
	if mound.Unit == "y"{
		options[0].Selected = true
	}else if mound.Unit == "m"{
		options[1].Selected = true
	}
	inputs := []Inputs{
		{
			Name: "clubid",
			Html:  "hidden",
			Value: clubId,
		},
		{
			Name: "name",
			Html:  "text",
			Label: "Mound Name",
			Value: mound.Name,
		},
		{
			Name: "distance",
			Html:	"number",
			Label: "Distance",
			Required: true,
			Min: &min,
			Max: &max,
			Step: 1,
			Value: fmt.Sprintf("%v", mound.Distance),
		},
		{
			Name: "unit",
			Html: "select",
			Required: true,
			Label: "Unit",
			Options: options,
		},
		{
			Name: "submit",
			Html:  "submit",
			Value: submitLabel,
		},
	}
	if existing{
		inputs = append(inputs, Inputs{
			Name: "moundid",
			Html:  "hidden",
			Value: fmt.Sprintf("%v", mound.Id),
		})
	}
	return Form{
		Action: URL_clubMoundInsert,
		Title: "Add Shooting Mounds",
		Inputs: inputs,
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := check_form(clubMoundInsertForm("", Mound{}, true).Inputs, r)
	clubId := validatedValues["clubid"]
	club, err := getClub(clubId)
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	if err != nil{
		Error.Println("Unable to insert Mound")
		return
	}
	distance, _ := strToInt(validatedValues["distance"])
	moundId := validatedValues["moundid"]
	newMound := Mound{
		Id: club.AutoInc.Mound,
		Name: validatedValues["name"],
		Distance: distance,
		Unit: validatedValues["unit"],
	}
	vardump(validatedValues)
	if moundId != "" {
		newMound.Id, err = strToInt(moundId)
		if err != nil{
			Error.Println("Unable to update club")
			return
		}
		for index, mound := range club.Mounds{
			if mound.Id == newMound.Id{
				club.Mounds[index] = newMound
				break
			}
		}
	}else{
		club.Mounds = append(club.Mounds, newMound)
		club.AutoInc.Mound += 1
	}
	UpdateDoc_by_id(TBLclub, clubId, club)
	http.Redirect(w, r, URL_club + clubId, http.StatusSeeOther)
}

func clubMapUpsertForm(clubId, latitude, longitude string) Form {
	latMin := -90
	latMax := 90
	longMin := -180
	longMax := 180
	return Form{
		Action: URL_clubMapUpsert,
		Title:  "Update Map Location",
		Inputs: []Inputs{
			{
				Name: "clubid",
				Html:  "hidden",
				Value: clubId,
			},
			{
				Name: "latitude",
				Html:     "number",
				Label:    "Latitude",
				Required: true,
				Min:      &latMin,
				Step:     0.000001,
				Max:      &latMax,
				Value:    latitude,
			},
			{
				Name: "longitude",
				Html:     "number",
				Required: true,
				Label:    "Longitude",
				Min:      &longMin,
				Step:     0.000001,
				Max:      &longMax,
				Value:    longitude,
			},
			{
				Name: "submit",
				Html:  "submit",
				Value: "Update Co-ordinates",
			},
		},
	}
}

func clubMapUpsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := check_form(clubMapUpsertForm("", "", "").Inputs, r)
	clubId := validatedValues["clubid"]
	club, err := getClub(clubId)
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	http.Redirect(w, r, URL_club + clubId, http.StatusSeeOther)
	if err != nil{
		Error.Println("Unable to update club map")
		return
	}
	club.Latitude = validatedValues["latitude"]
	club.Longitude = validatedValues["longitude"]
	UpdateDoc_by_id(TBLclub, clubId, club)
}
