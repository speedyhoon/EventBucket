package main

import (
	"fmt"
	"net/http"
	"strconv"
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
				Snippet: "Last updated: " + lastUpdated + ".",
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
	validatedValues := checkForm(searchShooterForm().Inputs, r)
	query := M{
	//Ignore Deleted shooters. Selects not modified, updated & merged shooters
	//"$or": []M{{"t": nil}, {"t": M{"$lt": 3}}},
	}
	if validatedValues["surname"] != "" {
		query["s"] = M{"$regex": fmt.Sprintf(`^%v`, validatedValues["surname"]), "$options": "i"}
	}
	if validatedValues["first"] != "" {
		query["f"] = M{"$regex": fmt.Sprintf(`^%v`, validatedValues["first"]), "$options": "i"}
	}
	if validatedValues["club"] != "" {
		query["c"] = M{"$regex": fmt.Sprintf(`^%v`, validatedValues["club"]), "$options": "i"}
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

func searchShooterForm() Form {
	return Form{
		Action: URL_shooterInsert,
		Title:  "Add Shooters",
		Inputs: []Inputs{
			{
				Name:  "first",
				Html:  "text",
				Label: "First Name",
			}, {
				Name:  "surname",
				Html:  "text",
				Label: "Surname",
			}, {
				Name:  "club",
				Html:  "text",
				Label: "Club",
			},
		},
	}
}

func searchShooterGrade(w http.ResponseWriter, r *http.Request) {
	output := ""
	validatedValues := checkForm(searchShooterGradeForm().Inputs, r)
	shooterId, err := strconv.Atoi(validatedValues["shooterid"])
	if err != nil {
		fmt.Fprint(w, output)
		return
	}
	shooter := getNraaShooter(shooterId) //TODO change to getShooter after it is moved
	output += fmt.Sprintf("%v %v", shooter.FirstName, shooter.Surname)
	if len(shooter.Grades) == 0 {
		output += "<div>No grades listed</div>"
	}
	for _, grade := range shooter.Grades {
		output += fmt.Sprintf("<div>Class: %v, Grade: %v, Threshold: %v</div>", grade.DisciplineName, grade.GradeName, grade.GradeThreshold)
	}
	fmt.Fprint(w, output)
}

func searchShooterGradeForm() Form {
	return Form{
		Action: URL_shooterInsert,
		Title:  "Shooters Grades",
		Inputs: []Inputs{
			{
				Name: "shooterid",
				Html: "number",
			},
		},
	}
}
