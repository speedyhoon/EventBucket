package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func shooterInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterForm(Event{}, true).Inputs, r)
	eventId := validatedValues["eventid"]
	http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)
	shooter := EventShooter{
		FirstName: validatedValues["first"],
		Surname:   validatedValues["surname"],
		Club:      validatedValues["club"],
		Grade:     str2Int(validatedValues["grade"]),
		AgeGroup:  validatedValues["age"],
	}
	eventShooterInsert(eventId, shooter)
}

func event(eventId string) Page {
	event, err := getEvent(eventId)
	if err != nil {
		return Page{
			TemplateFile: "event",
			Theme:        TEMPLATE_ADMIN,
			Title:        "Event not found: " + eventId,
			Data: M{
				"Menu":  eventMenu("", []Range{}, "", false),
				"Valid": false,
			},
		}
	}
	clubList := dataListShooterClubNames()
	var editShooterList string
	for _, shooter := range event.Shooters {
		editShooterList += generateForm(eventUpdateShooterForm(event, shooter, clubList))
	}

	//TODO add ClubOptions when club textbox is changed to a datalist
	return Page{
		TemplateFile: "event",
		Theme:        TEMPLATE_ADMIN,
		Title:        "Event",
		Data: M{
			"EventName":            event.Name,
			"Menu":                 eventMenu(eventId, event.Ranges, URL_event, event.IsPrizeMeet),
			"Valid":                true,
			"NewShooterEntry":      URL_shooterInsert,
			"GradeOptions":         drawOptions(Inputs{Options: eventGradeOptions(event.Grades)}),
			"EntryOptions":         drawOptions(Inputs{Options: eventGradeEntry(event.Grades)}),
			"AgeOptions":           drawOptions(Inputs{Options: AgeGroups()}),
			"ExistingShooterEntry": URL_shooterListInsert,
			"EventGrades":          generateForm(eventSettingsClassGrades(event)),
			"ListShooters":         event.Shooters,
			"EditShooterList":      editShooterList,
			"EventId":              eventId,
			"DataListClubs":        drawOptions(Inputs{Id: "clubSearch", DataList: true, Options: dataListShooterClubNames()}), //TODO optimise by saving as a global var so it is only executed once per club update/insert
			"ShooterQty":           len(event.Shooters),
		},
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterForm(Event{}, false).Inputs, r)
	eventId := validatedValues["eventid"]
	http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)

	shooterId, err := strconv.Atoi(validatedValues["sid"])
	if err != nil {
		Error.Printf("shooter id %v is not valid", validatedValues["sid"])
		return
	}
	shooter := getShooterList(shooterId)
	eventShooter := EventShooter{
		Grade:     str2Int(validatedValues["grade"]),
		AgeGroup:  validatedValues["age"],
		FirstName: shooter.NickName,
		Surname:   shooter.Surname,
		Club:      shooter.Club,
	}
	eventShooterInsert(eventId, eventShooter)
}

func eventAddShooterForm(event Event, new bool) Form {
	return Form{
		Inputs: []Inputs{
			{
				Name:     "first",
				Html:     "text",
				Required: new,
			},
			{
				Name:     "surname",
				Html:     "text",
				Required: new,
			},
			{
				Name:     "club",
				Html:     "text",
				Required: new,
			},
			{
				Name:     "age",
				Html:     "select",
				Required: true,
			}, {
				Name:     "grade",
				Html:     "select",
				Required: true,
			}, {
				Name:     "sid",
				Html:     "select",
				Required: !new,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: event.Id,
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
	for gradeId, grade := range grades() {
		options = append(options, Option{ //TODO change Value to an interface so type conversion isn't needed
			Value:    fmt.Sprintf("%v", gradeId),
			Display:  grade.LongName,
			Selected: isValueInSlice(gradeId, eventGrades) || allSelected,
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
	for _, gradeId := range gradesToSelectFrom {
		options = append(options, Option{ //TODO change Value to an interface so type conversion isn't needed
			Value:   fmt.Sprintf("%v", gradeId),
			Display: allGrades[gradeId].LongName,
		})
	}
	return options
}

func eventUpdateShooter(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventUpdateShooterForm(Event{}, EventShooter{}, []Option{}).Inputs, r)
	eventId := validatedValues["eventid"]
	shooterId := validatedValues["sid"]
	if eventId != "" && shooterId != "" {
		http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)
		updateData := M{Dot(schemaSHOOTER, shooterId): M{
			"f": validatedValues["first"],
			"s": validatedValues["surname"],
			"c": validatedValues["club"],
			"g": str2Int(validatedValues["grade"]),
			"b": validatedValues["disabled"] != "",
			"a": validatedValues["age"],
		}}
		tableUpdateData(TBLevent, eventId, updateData)
		return
	}
	http.Redirect(w, r, URL_home, http.StatusSeeOther)
}

func dataListShooterClubNames() []Option {
	var clubList []Option
	for _, club := range getClubs() {
		clubList = append(clubList, Option{
			Value: club.Name,
		})
	}
	return clubList
}

func eventUpdateShooterForm(event Event, shooter EventShooter, clubList []Option) Form {
	var options []Option
	allGrades := grades()
	if len(event.Grades) == 0 {
		event.Grades = gradeList()
	}
	for _, gradeId := range event.Grades {
		options = append(options, Option{
			Display:  allGrades[gradeId].LongName,
			Value:    gradeId,
			Selected: shooter.Grade == gradeId,
		})
	}
	return Form{
		Action: URL_eventUpdateShooter,
		Type:   "table",
		Id:     fmt.Sprintf("update%v", shooter.Id),
		Inputs: []Inputs{
			{
				Snippet: shooter.Id,
			}, {
				Name:    "disabled",
				Html:    "checkbox",
				Checked: shooter.Disabled,
				//AccessKey: "d",
			}, {
				Name:  "first",
				Html:  "search",
				Value: shooter.FirstName,
				//AccessKey: "f",
			}, {
				Name:  "surname",
				Html:  "search",
				Value: shooter.Surname,
				//AccessKey: "s",
			}, {
				Name:         "club",
				Html:         "search",
				DataList:     true,
				AutoComplete: "off",
				Id:           fmt.Sprintf("club%v%v", shooter.Id, event.Id),
				Value:        shooter.Club,
				Options:      clubList,
				//AccessKey: "c",
			}, {
				Name:    "grade",
				Html:    "select",
				Options: options,
				//AccessKey: "g",
			}, {
				Name:    "age",
				Html:    "select",
				Options: shooterAgeGroupSelectbox(shooter),
				//AccessKey: "a",
			}, {
				Html:  "submit",
				Inner: "Save",
				Name:  "sid",
				Value: shooter.Id,
				//AccessKey: "x",
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: event.Id,
			},
		},
	}
}
