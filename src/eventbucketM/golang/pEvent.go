package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func shooterInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterForm("", []int{}).Inputs, r)
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
	event, err := getEvent20Shooters(eventId)
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
			"GradeOptions":         drawOptions(Inputs{Options: eventGradeOptions(event.Grades)}, ""),
			"EntryOptions":         drawOptions(Inputs{Options: eventGradeEntry(event.Grades)}, ""),
			"AgeOptions":           drawOptions(Inputs{Options: AgeGroups()}, ""),
			"ExistingShooterEntry": URL_shooterListInsert,
			"EventGrades":          generateForm(eventSettingsClassGrades(event)),
			"ListShooters":         event.Shooters,
			"EventId":              eventId,
		},
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventAddShooterListForm(Event{}).Inputs, r)
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

func eventAddShooterForm(eventId string, grades []int) Form {
	return Form{
		Action: URL_shooterInsert,
		Title:  "Add Shooters",
		Inputs: []Inputs{
			{
				Name:     "first",
				Html:     "text",
				Label:    "First Name",
				Required: true,
			}, {
				Name:     "surname",
				Html:     "text",
				Label:    "Surname",
				Required: true,
			}, {
				Name: "club",
				Html: "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:    "Club",
				Required: true,
			}, {
				Name:     "age",
				Html:     "select",
				Label:    "Age Group",
				Options:  AgeGroups(),
				Required: true,
			}, {
				Name:        "grade",
				Html:        "select",
				Label:       "Class & Grade",
				Placeholder: "Class & Grade",
				Required:    true,
				Options:     eventGradeOptions(grades),
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			},
		},
	}
}

func eventAddShooterListForm(event Event) Form {
	return Form{
		Action: URL_shooterInsert,
		Title:  "Add Shooters",
		Inputs: []Inputs{
			//"first":Inputs{
			//	Html:      "text",
			//	Label:   "First Name",
			//	Required: true,
			//},
			//"surname":Inputs{
			//	Html:      "text",
			//	Label:   "Surname",
			//	Required: true,
			//},
			//"club":Inputs{
			//	Html:      "text",
			//	TODO change club to a data-list
			//	//SelectValues:   getClubSelectBox(eventsCollection),
			//	Label:   "Club",
			//	Required: true,
			//},
			{
				Name:     "age",
				Html:     "select",
				Label:    "Age Group",
				Options:  AgeGroups(),
				Required: true,
			}, {
				Name:        "grade",
				Html:        "select",
				Label:       "Class & Grade",
				Placeholder: "Class & Grade",
				Required:    true,
				Options:     eventGradeOptions(event.Grades),
			}, {
				Name: "sid",
				Html: "select",
				//				Label: "Class & Grade",
				//				Placeholder: "Class & Grade",
				Required: true,
				//				SelectedValues: eventGradeOptions(event),
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
	for _, gradeId := range gradesToSelectFrom {
		options = append(options, Option{ //TODO change Value to an interface so type conversion isn't needed
			Value:   fmt.Sprintf("%v", gradeId),
			Display: allGrades[gradeId].LongName,
		})
	}
	return options
}
