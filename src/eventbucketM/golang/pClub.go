package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func club(clubID string) Page {
	page := Page{
		TemplateFile: "club",
		Title:        "Club",
		Theme:        templateHome,
		Data: M{
			"MoundForm":   "",
			"TableMounds": "",
			"Latitude":    "",
			"Longitude":   "",
			"Details":     "",
		},
	}
	club, err := getClub(clubID)
	if err != nil {
		//TODO return a 404 error
		return page
	}
	var mounds []string
	for index, mound := range club.Mounds {
		mounds = append(mounds, generateForm(clubMoundInsertForm(clubID, mound, true, fmt.Sprintf("formMounds%v", index))))
	}
	tableMounds := "<p>No Mounds Listed</p>"
	if len(mounds) >= 1 {
		tableMounds = `<table>
			<tr>
				<th>Name
				<th>Distance
				<th colspan=2 class=text-left>Units`
		tableMounds += strings.Join(mounds, "") + "</table>"
	}

	page.Title = club.Name
	page.Data["MoundForm"] = generateForm(clubMoundInsertForm(clubID, Mound{}, false, ""))
	//	page.Data["MapForm"] =  generateForm2(clubMapUpsertForm(clubID, club.Latitude, club.Longitude))
	//	page.Data["ListMounds"] = mounds
	page.Data["TableMounds"] = tableMounds
	page.Data["Latitude"] = club.Latitude
	page.Data["Longitude"] = club.Longitude
	page.Data["Details"] = generateForm(clubDetailsForm(club))
	return page
}

//,"<input type=checkbox checked class=map>Display Map <abbr class=help title=\"Longitude and Latitude decimal format is three digits with six decimal places\ne.g. 000.000000 or -000.000000\n\nTo find your clubs longitude and latitude location go to Google Maps, then right click the Map and select What's here?\n\nTip:\nLines of longitude appear vertical (North-South)\nLines of latitude appear horizontal (East-West)\">?</abbr>"
func clubMoundInsertForm(clubID string, mound Mound, existing bool, formID string) Form {
	min := 10
	max := 9750
	submitLabel := "Create Mound"
	var formTable bool
	if existing {
		submitLabel = "Save Mound"
		formTable = true
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
			name:  "name",
			html:  "text",
			label: "Mound Name",
			value: mound.Name,
		}, {
			name:     "distance",
			html:     "number",
			label:    "Distance",
			required: true,
			min:      &min,
			max:      &max,
			step:     1,
			value:    fmt.Sprintf("%v", mound.Distance),
		}, {
			name:     "unit",
			html:     "select",
			required: true,
			label:    "Unit",
			options:  options,
		}, {
			html:  "submit",
			inner: submitLabel,
			name:  "clubid",
			value: clubID,
		},
	}
	if existing {
		inputs = append(inputs, Inputs{
			name:  "moundid",
			html:  "hidden",
			value: fmt.Sprintf("%v", mound.ID),
		})
	}
	return Form{
		action: urlClubMoundInsert,
		title:  "Add Shooting Mounds",
		inputs: inputs,
		table:  formTable,
		id:     formID,
	}
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(clubMoundInsertForm("", Mound{}, true, "").inputs, r)
	clubID := validatedValues["clubid"]
	club, err := getClub(clubID)
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	if err != nil {
		Error.Println("Unable to insert Mound")
		return
	}
	distance, _ := strconv.Atoi(validatedValues["distance"])
	moundID := validatedValues["moundid"]
	newMound := Mound{
		ID:       club.AutoInc.Mound,
		Name:     validatedValues["name"],
		Distance: distance,
		Unit:     validatedValues["unit"],
	}
	if moundID != "" {
		newMound.ID, err = strconv.Atoi(moundID)
		if err != nil {
			Error.Println("Unable to update club")
			return
		}
		for index, mound := range club.Mounds {
			if mound.ID == newMound.ID {
				club.Mounds[index] = newMound
				break
			}
		}
	} else {
		club.Mounds = append(club.Mounds, newMound)
		club.AutoInc.Mound++
	}
	updateDocByID(tblClub, clubID, club)
	http.Redirect(w, r, urlClub+clubID, http.StatusSeeOther)
}

func clubDetailsForm(club Club) Form {
	return Form{
		action: urlClubDetailsUpsert,
		title:  "Club Details",
		inputs: []Inputs{
			{
				name:     "name",
				html:     "text",
				label:    "Name",
				required: true,
				value:    club.Name,
			}, {
				name:     "address",
				html:     "text",
				label:    "Address",
				required: true,
				value:    club.Address,
			}, {
				name:     "town",
				html:     "text",
				label:    "Town",
				required: true,
				value:    club.Town,
			}, {
				name:  "postcode",
				html:  "text",
				label: "Post Code",
				value: club.PostCode,
			}, {
				name:     "latitude",
				html:     "number",
				label:    "Latitude",
				required: true,
				min:      &latitudeMin,
				step:     0.000001,
				max:      &latitudeMax,
				value:    club.Latitude,
			}, {
				name:     "longitude",
				html:     "number",
				required: true,
				label:    "Longitude",
				min:      &longitudeMin,
				step:     0.000001,
				max:      &longitudeMax,
				value:    club.Longitude,
				help: `To find your clubs latitude &amp; longitude, go to Google Maps, then right click the Map and select What's here?
Longitude and Latitude decimal format is three digits with six decimal places e.g. 000.000000 or -000.000000
Tip: Lines of longitude appear vertical (North-South), Lines of latitude appear horizontal (East-West).`,
			}, {
				html:  "submit",
				inner: "Save Club Details",
				name:  "clubid",
				value: club.ID,
			}, {
				snippet: "<a href=//maps.google.com.au/ target=_blank>Google Maps</a>",
			},
		},
	}
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request) {
	var club Club
	var err error
	validatedValues := checkForm(clubDetailsForm(club).inputs, r)
	club, err = getClub(validatedValues["clubid"])
	//TODO instead return true/false as success/failure and on failure return a filled out form (bool, FormInvalid, error.Message as string)
	//TODO FormInvalid.Inputs { Name: "", Html, "number, Message: "Number is greater than 9750", }...
	if err != nil {
		Error.Println("Unable to update club details")
		http.Redirect(w, r, urlClubs, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, urlClub+club.ID, http.StatusSeeOther)
	club.Name = validatedValues["name"]
	club.Address = validatedValues["address"]
	club.PostCode = validatedValues["postcode"]
	club.Town = validatedValues["town"]
	club.Latitude = validatedValues["latitude"]
	club.Longitude = validatedValues["longitude"]
	updateDocByID(tblClub, club.ID, club)
}

func clubs() Page {
	return Page{
		TemplateFile: "clubs",
		Title:        "Clubs",
		Theme:        templateHome,
		Data: M{
			"Title":    "Clubs",
			"Clubs":    generateForm(organisersClubForm()),
			"ClubList": getClubs(),
		},
	}
}

func organisersClubForm() Form {
	//TODO add validation to
	return Form{
		action: urlClubInsert,
		title:  "New Club",
		inputs: []Inputs{
			{
				name:      "name",
				html:      "text",
				label:     "Club Name",
				required:  true,
				autofocus: true,
			},
			{
				html:  "submit",
				inner: "Add Club",
			},
		},
	}
}

func clubInsert(w http.ResponseWriter, r *http.Request) {
	insertClub(checkForm(organisersClubForm().inputs, r)["name"])
}

func insertClub(clubName string) (string, error) {
	nextID, err := getNextID(tblClub)
	if err == nil {
		newClub := Club{
			Name: clubName,
			ID:   nextID,
		}
		newClub.AutoInc.Mound = 1
		insertDoc(tblClub, newClub)
		return newClub.ID, nil
	}
	return "", errors.New("Unable to generate club id")
}
