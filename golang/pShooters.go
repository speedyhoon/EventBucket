package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func shooters() Page {
	return Page{
		TemplateFile: "shooters",
		Theme:        templateHome,
		Title:        "Shooters",
		Data: M{
			//"makeList":   generateForm(makeShooterList()),
			"updateList": generateForm(makeShooterList()),
			//"updateList": generateForm(updateShooterList()),
		},
	}
}

func makeShooterList() Form {
	//TODO move support for loading nraa data to dev mode only
	return Form{
		action: urlMakeShooterList,
		title:  "Generate Shooter List",
		inputs: []Inputs{
			{
				html:      "submit",
				inner:     "Update",
				autofocus: true,
			},
		},
	}
}

/*func updateShooterList() Form {
	return Form{
		action: urlUpdateShooterList,
		title:  "Update Shooter List",
		inputs: []Inputs{
			{
				snippet: "Last updated: <strong>" + nraaGetLastUpdated() + "</strong>",
			}, {
				//TODO add support for loading a JSON file
				//	html:      "file",
				//	name:      "source",
				//}, {
				html:      "submit",
				inner:     "Update",
				autofocus: "on",
			},
		},
	}
}*/

//Search for a shooter by first name, surname or club
func searchShooter(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(searchShooterForm().inputs, r)
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
	fmt.Fprint(w, drawOptions(Inputs{options: optionList}))
}

func searchShooterForm() Form {
	return Form{
		action: urlShooterInsert,
		title:  "Add Shooters",
		inputs: []Inputs{
			{
				name:  "first",
				html:  "text",
				label: "First Name",
			}, {
				name:  "surname",
				html:  "text",
				label: "Surname",
			}, {
				name:  "club",
				html:  "text",
				label: "Club",
			},
		},
	}
}

func searchShooterGrade(w http.ResponseWriter, r *http.Request) {
	output := ""
	validatedValues := checkForm(searchShooterGradeForm().inputs, r)
	shooterID, err := strconv.Atoi(validatedValues["shooterid"])
	if err != nil {
		fmt.Fprint(w, output)
		return
	}
	shooter := getNraaShooter(shooterID) //TODO change to getShooter after it is moved
	output += fmt.Sprintf("%v %v", shooter.FirstName, shooter.Surname)
	if len(shooter.Grades) == 0 {
		//TODO put the div tag in the html template
		output += "<div>No grades listed</div>"
	}
	for _, grade := range shooter.Grades {
		//TODO remove the div tag and just join the text with <br>'s
		output += fmt.Sprintf("<div>Class: %v, Grade: %v, Threshold: %v</div>", grade.Class, grade.Grade, grade.Threshold)
	}
	fmt.Fprint(w, output)
}

func searchShooterGradeForm() Form {
	return Form{
		action: urlShooterInsert,
		title:  "Shooters Grades",
		inputs: []Inputs{
			{
				name: "shooterid",
				html: "number",
			},
		},
	}
}
