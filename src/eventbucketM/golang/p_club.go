package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func club(clubId string) Page {
	page := Page{
		TemplateFile: "club",
		Theme:        TEMPLATE_HOME,
		Title:        "Club with id '" + clubId + "' not found",
		Data: M{
			"Menu": home_menu(URL_club, HOME_MENU_ITEMS),
		},
	}
	club, err := getClub(clubId)
	if err != nil {
		//TODO return a 404 error
		return page
	}
	var mounds []string
	for index, mound := range club.Mounds {
		mounds = append(mounds, generateForm2(clubMoundInsertForm(clubId, mound, true, fmt.Sprintf("formMounds%v", index))))
	}
	tableMounds := "<p>No Mounds Listed</p>"
	if len(mounds) >= 1 {
		tableMounds = `
		<table>
			<tr>
				<th>Name</th>
				<th>Distance</th>
				<th colspan=2 class=text-left>Units</th>
			</tr>`
		tableMounds += strings.Join(mounds, "") + "</table>"
	}

	page.Title = club.Name
	page.Data["MoundForm"] = generateForm2(clubMoundInsertForm(clubId, Mound{}, false, ""))
	//	page.Data["MapForm"] =  generateForm2(clubMapUpsertForm(clubId, club.Latitude, club.Longitude))
	//	page.Data["ListMounds"] = mounds
	page.Data["TableMounds"] = tableMounds
	page.Data["Latitude"] = club.Latitude
	page.Data["Longitude"] = club.Longitude
	page.Data["Details"] = generateForm2(clubDetailsForm(club))
	return page
}

//,"<input type=checkbox checked class=map>Display Map <abbr class=help title=\"Longitude and Latitude decimal format is three digits with six decimal places\ne.g. 000.000000 or -000.000000\n\nTo find your clubs longitude and latitude location go to Google Maps, then right click the Map and select What's here?\n\nTip:\nLines of longitude appear vertical (North-South)\nLines of latitude appear horizontal (East-West)\">?</abbr>"
func clubMoundInsertForm(clubId string, mound Mound, existing bool, formId string) Form {
	min := 10
	max := 9750
	submitLabel := "Create Mound"
	formType := ""
	if existing {
		submitLabel = "Save Mound"
		formType = "table"
	}
	options := []Option{
		{
			Display: "Yards",
			Value:   "y",
		},
		{
			Display: "Metres",
			Value:   "m",
		},
	}
	if mound.Unit == "y" {
		options[0].Selected = true
	} else if mound.Unit == "m" {
		options[1].Selected = true
	}
	inputs := []Inputs{
		{
			Name:  "name",
			Html:  "text",
			Label: "Mound Name",
			Value: mound.Name,
		}, {
			Name:     "distance",
			Html:     "number",
			Label:    "Distance",
			Required: true,
			Min:      &min,
			Max:      &max,
			Step:     1,
			Value:    fmt.Sprintf("%v", mound.Distance),
		}, {
			Name:     "unit",
			Html:     "select",
			Required: true,
			Label:    "Unit",
			Options:  options,
		}, {
			Name:  "submit",
			Html:  "submit",
			Value: submitLabel,
		}, {
			Name:  "clubid",
			Html:  "hidden",
			Value: clubId,
		},
	}
	if existing {
		inputs = append(inputs, Inputs{
			Name:  "moundid",
			Html:  "hidden",
			Value: fmt.Sprintf("%v", mound.Id),
		})
	}
	return Form{
		Action: URL_clubMoundInsert,
		Title:  "Add Shooting Mounds",
		Inputs: inputs,
		Type:   formType,
		Id:     formId,
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := check_form(clubMoundInsertForm("", Mound{}, true, "").Inputs, r)
	clubId := validatedValues["clubid"]
	club, err := getClub(clubId)
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	if err != nil {
		Error.Println("Unable to insert Mound")
		return
	}
	distance, _ := strToInt(validatedValues["distance"])
	moundId := validatedValues["moundid"]
	newMound := Mound{
		Id:       club.AutoInc.Mound,
		Name:     validatedValues["name"],
		Distance: distance,
		Unit:     validatedValues["unit"],
	}
	if moundId != "" {
		newMound.Id, err = strToInt(moundId)
		if err != nil {
			Error.Println("Unable to update club")
			return
		}
		for index, mound := range club.Mounds {
			if mound.Id == newMound.Id {
				club.Mounds[index] = newMound
				break
			}
		}
	} else {
		club.Mounds = append(club.Mounds, newMound)
		club.AutoInc.Mound += 1
	}
	UpdateDoc_by_id(TBLclub, clubId, club)
	http.Redirect(w, r, URL_club+clubId, http.StatusSeeOther)
}

func clubDetailsForm(club Club) Form {
	return Form{
		Action: URL_clubDetailsUpsert,
		Title:  "Club Details",
		Inputs: []Inputs{
			{
				Name:  "clubid",
				Html:  "hidden",
				Value: club.Id,
			}, {
				Name:     "name",
				Html:     "text",
				Label:    "Name",
				Required: true,
				Value:    club.Name,
			}, {
				Name:     "address",
				Html:     "text",
				Label:    "Address",
				Required: true,
				Value:    club.Address,
			}, {
				Name:     "town",
				Html:     "text",
				Label:    "Town",
				Required: true,
				Value:    club.Town,
			}, {
				Name:  "postcode",
				Html:  "text",
				Label: "Post Code",
				Value: club.PostCode,
			}, {
				Name:     "latitude",
				Html:     "number",
				Label:    "Latitude",
				Required: true,
				Min:      &LATITUDE_MIN,
				Step:     0.000001,
				Max:      &LATITUDE_MAX,
				Value:    club.Latitude,
			}, {
				Name:     "longitude",
				Html:     "number",
				Required: true,
				Label:    "Longitude",
				Min:      &LONGITUDE_MIN,
				Step:     0.000001,
				Max:      &LONGITUDE_MAX,
				Value:    club.Longitude,
				Help: `To find your clubs latitude &amp; longitude, go to Google Maps, then right click the Map and select What's here?
Longitude and Latitude decimal format is three digits with six decimal places e.g. 000.000000 or -000.000000
Tip: Lines of longitude appear vertical (North-South), Lines of latitude appear horizontal (East-West).`,
			}, {
				Html:  "submit",
				Value: "Save Club Details",
			}, {
				Snippet: "<a href=//maps.google.com.au/ target=_blank>Google Maps</a>",
			},
		},
	}
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request) {
	var club Club
	var err error
	validatedValues := check_form(clubDetailsForm(club).Inputs, r)
	club, err = getClub(validatedValues["clubid"])
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	if err != nil {
		Error.Println("Unable to update club details")
		http.Redirect(w, r, URL_clubs, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, URL_club+club.Id, http.StatusSeeOther)
	club.Name = validatedValues["name"]
	club.Address = validatedValues["address"]
	club.PostCode = validatedValues["postcode"]
	club.Town = validatedValues["town"]
	club.Latitude = validatedValues["latitude"]
	club.Longitude = validatedValues["longitude"]
	UpdateDoc_by_id(TBLclub, club.Id, club)
}

func clubs() Page {
	return Page{
		TemplateFile: "clubs",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Title":    "Clubs",
			"Clubs":    generateForm2(organisers_clubForm()),
			"ClubList": getClubs(),
			"Menu":     home_menu(URL_clubs, HOME_MENU_ITEMS),
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
				Name:      "name",
				Html:      "text",
				Label:     "Club Name",
				Required:  true,
				Autofocus: "on",
			},
			{
				Html:  "submit",
				Value: "Add Club",
			},
		},
	}
}

func getClubSelectionBox(club_list []Club) []Option {
	var drop_down []Option
	for _, club := range club_list {
		drop_down = append(drop_down, Option{Display: club.Name, Value: club.Id})
	}
	return drop_down
}

func clubInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(organisers_clubForm().Inputs, r)
	insertClub(validated_values["name"])
}

func insertClub(clubName string) (string, error) {
	nextId, err := getNextId(TBLclub)
	if err == nil {
		newClub := Club{
			Name: clubName,
			Id:   nextId,
		}
		newClub.AutoInc.Mound = 1
		InsertDoc(TBLclub, newClub)
		return newClub.Id, nil
	}
	return "", errors.New("Unable to generate club id")
}
