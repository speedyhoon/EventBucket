package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func shooters() Page {
	return Page{
		TemplateFile: "shooters",
		Theme:        TEMPLATE_HOME,
		Title:        "Shooters",
		Data: M{
			"Menu":        homeMenu(URL_shooters, HOME_MENU_ITEMS),
			"ShooterList": generateForm(updateShooterList()),
		},
	}
}

func updateShooterList() Form {
	lastUpdated := nraaGetLastUpdated()
	if lastUpdated == "" {
		lastUpdated = "Never"
	}
	return Form{
		Action: URL_updateShooterList,
		Title:  "Update Shooter List",
		Inputs: []Inputs{
			{
				Snippet: "Last updated: " + lastUpdated+".",
			}, {
				Html:      "submit",
				Value:     "Update",
				Autofocus: "on",
			},
		},
	}
}

//Search for a shooter by first name, surname or club
func searchShooter(w http.ResponseWriter, r *http.Request) {
	var shooters Shooter
	err := json.NewDecoder(r.Body).Decode(&shooters)
	if err != nil {
		Error.Println(err)
		return
	}
	query := M{
		//Ignore Deleted shooters. Selects not modified, updated & merged shooters
		"$or": []M{{"t": nil}, {"t": M{"$lt": 3}}},
	}
	if shooters.Surname != "" {
		query["s"] = M{"$regex": fmt.Sprintf(`^%v`, shooters.Surname), "$options": "i"}
	}
	if shooters.FirstName != "" {
		query["f"] = M{"$regex": fmt.Sprintf(`^%v`, shooters.FirstName), "$options": "i"}
	}
	if shooters.Club != "" {
		query["c"] = M{"$regex": fmt.Sprintf(`^%v`, shooters.Club), "$options": "i"}
	}
	var optionList []Option
	for _, shooter := range searchShooters(query) {
		optionList = append(optionList, Option{
			Value:   fmt.Sprintf("%v", shooter.SID),
			Display: fmt.Sprintf("%v %v, ~~ %v", shooter.FirstName, shooter.Surname, shooter.Club),
		})
	}
	fmt.Fprint(w, drawOptions(Inputs{Options: optionList}, ""))
}
