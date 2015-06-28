package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func shooterInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterForm(Event{}, true).inputs, r)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)
	shooter := EventShooter{
		FirstName: validatedValues["first"],
		Surname:   validatedValues["surname"],
		Club:      validatedValues["club"],
		Grade:     str2Int(validatedValues["grade"]),
		AgeGroup:  validatedValues["age"],
	}
	eventShooterInsert(eventID, shooter)
}

func event(eventID string) Page {
	event, err := getEvent(eventID)
	if err != nil {
		return Page{
			TemplateFile: "event",
			Theme:        templateAdmin,
			Title:        "Event not found: " + eventID,
			Data: M{
				"Menu":  eventMenu("", []Range{}, "", false),
				"Valid": false,
			},
		}
	}
	var editShooterList string
	for _, shooter := range event.Shooters {
		editShooterList += generateForm(eventUpdateShooterForm(event, shooter))
	}

	//TODO add ClubOptions when club textbox is changed to a datalist
	return Page{
		TemplateFile: "event",
		Theme:        templateAdmin,
		Title:        "Event",
		Data: M{
			"EventName":            event.Name,
			"Menu":                 eventMenu(eventID, event.Ranges, urlEvent, event.IsPrizeMeet),
			"Valid":                true,
			"NewShooterEntry":      urlShooterInsert,
			"GradeOptions":         drawOptions(Inputs{options: eventGradeOptions(event.Grades)}),
			"EntryOptions":         drawOptions(Inputs{options: eventGradeEntry(event.Grades)}),
			"AgeOptions":           drawOptions(Inputs{options: ageGroups()}),
			"ExistingShooterEntry": urlShooterListInsert,
			"EventGrades":          generateForm(eventSettingsClassGrades(event)),
			"ListShooters":         event.Shooters,
			"EditShooterList":      editShooterList,
			"EventID":              eventID,
			"DataListClubs":        drawOptions(Inputs{id: "clubSearch", dataList: true, options: dataListShooterClubNames()}), //TODO optimise by saving as a global var so it is only executed once per club update/insert
			"ShooterQty":           len(event.Shooters),
		},
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterForm(Event{}, false).inputs, r)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)

	shooterID, err := strconv.Atoi(validatedValues["sid"])
	if err != nil {
		warning.Printf("shooter id %v is not valid", validatedValues["sid"])
		return
	}
	shooter := getShooterList(shooterID)
	eventShooter := EventShooter{
		Grade:     str2Int(validatedValues["grade"]),
		AgeGroup:  validatedValues["age"],
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
	}
	eventShooterInsert(eventID, eventShooter)
}

func eventAddShooterForm(event Event, new bool) Form {
	return Form{
		inputs: []Inputs{
			{
				name:     "first",
				html:     "text",
				required: new,
			},
			{
				name:     "surname",
				html:     "text",
				required: new,
			},
			{
				name:     "club",
				html:     "text",
				required: new,
			},
			{
				name:     "age",
				html:     "select",
				required: true,
			}, {
				name:     "grade",
				html:     "select",
				required: true,
			}, {
				name:     "sid",
				html:     "select",
				required: !new,
			}, {
				name:  "eventid",
				html:  "hidden",
				value: event.ID,
			},
		},
	}
}

//Available Classes & Grades
func eventGradeOptions(eventGrades []int) []Option {
	//TODO add club defaults here
	allSelected := false
	if len(eventGrades) == 0 {
		allSelected = true
	}
	var options []Option
	for gradeID, grade := range grades() {
		options = append(options, Option{ //TODO change Value to an interface so type conversion isn't needed
			Value:    fmt.Sprintf("%v", gradeID),
			Display:  grade.longName,
			Selected: isValueInSlice(gradeID, eventGrades) || allSelected,
		})
	}
	return options
}

//shooter Entry grade to select from
func eventGradeEntry(gradesToSelectFrom []int) []Option {
	var options []Option
	allGrades := grades()
	if len(gradesToSelectFrom) == 0 {
		gradesToSelectFrom = gradeList()
	}
	for _, gradeID := range gradesToSelectFrom {
		options = append(options, Option{ //TODO change Value to an interface so type conversion isn't needed
			Value:   fmt.Sprintf("%v", gradeID),
			Display: allGrades[gradeID].longName,
		})
	}
	return options
}

func eventUpdateShooter(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventUpdateShooterForm(Event{}, EventShooter{}).inputs, r)
	eventID := validatedValues["eventid"]
	shooterID := validatedValues["sid"]
	if eventID != "" && shooterID != "" {
		http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)
		updateData := M{
			dot(schemaSHOOTER, shooterID, "f"): validatedValues["first"],
			dot(schemaSHOOTER, shooterID, "s"): validatedValues["surname"],
			dot(schemaSHOOTER, shooterID, "c"): validatedValues["club"],
			dot(schemaSHOOTER, shooterID, "g"): str2Int(validatedValues["grade"]),
			dot(schemaSHOOTER, shooterID, "b"): validatedValues["disabled"] != "",
			dot(schemaSHOOTER, shooterID, "a"): validatedValues["age"],
		}
		tableUpdateData(tblEvent, eventID, updateData)
		return
	}
	http.Redirect(w, r, urlHome, http.StatusSeeOther)
}

func dataListShooterClubNames() []Option {
	var clubList []Option
	for _, club := range getClubs() {
		if club.Name != "" {
			clubList = append(clubList, Option{
				Value: club.Name,
			})
		}
	}
	return clubList
}

func eventUpdateShooterForm(event Event, shooter EventShooter) Form {
	var options []Option
	if len(event.Grades) == 0 {
		event.Grades = gradeList()
	}
	allGrades := grades()
	for _, gradeID := range event.Grades {
		options = append(options, Option{
			Display:  allGrades[gradeID].longName,
			Value:    gradeID,
			Selected: shooter.Grade == gradeID,
		})
	}
	return Form{
		action: urlEventUpdateShooter,
		table:  true,
		id:     fmt.Sprintf("update%v", shooter.ID),
		inputs: []Inputs{
			{
				snippet: shooter.ID,
				name:    "eventid",
				html:    "hidden",
				value:   event.ID,
			}, {
				name:    "disabled",
				html:    "checkbox",
				checked: shooter.Disabled,
			}, {
				name:  "first",
				html:  "text",
				value: shooter.FirstName,
			}, {
				name:  "surname",
				html:  "text",
				value: shooter.Surname,
			}, {
				name:         "club",
				html:         "text",
				dataList:     true,
				autoComplete: "off",
				id:           "clubSearch",
				value:        shooter.Club,
			}, {
				name:    "grade",
				html:    "select",
				options: options,
			}, {
				name:    "age",
				html:    "select",
				options: shooterAgeGroupSelect(shooter),
			}, {
				html:  "submit",
				inner: "Save",
				name:  "sid",
				value: shooter.ID,
			},
		},
	}
}
