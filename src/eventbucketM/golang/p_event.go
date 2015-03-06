package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func shooterInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := checkForm(event_add_shooterForm("", []int{}).Inputs, r)
	eventId := validated_values["eventid"]
	http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)
	new_shooter := EventShooter{
		FirstName: validated_values["first"],
		Surname:   validated_values["surname"],
		Club:      validated_values["club"],
	}
	new_shooter.Grade, _ = strconv.Atoi(validated_values["grade"])
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	eventShooterInsert(eventId, new_shooter)
}

func event(eventId string) Page {
	event, err := getEvent20Shooters(eventId)
	if err {
		return Page{
			TemplateFile: "event",
			Theme:        TEMPLATE_ADMIN,
			Data: M{
				"Title": "Event not found: " + eventId,
				"Menu":  eventMenu("", []Range{}, "", false),
				"Valid": false,
			},
		}
	}
	return Page{
		TemplateFile: "event",
		Theme:        TEMPLATE_ADMIN,
		Title:        "Event",
		Data: M{
			"EventName":       event.Name,
			"Menu":            eventMenu(eventId, event.Ranges, URL_event, event.IsPrizeMeet),
			"Valid":           true,
			"NewShooterEntry": URL_shooterInsert,
			"GradeOptions":    drawOptions(Inputs{Options: eventGradeOptions(event.Grades)}, ""),
			//TODO add ClubOptions when club textbox is changed to a datalist
			"AgeOptions":           drawOptions(Inputs{Options: AGE_GROUPS2()}, ""),
			"ExistingShooterEntry": URL_shooterListInsert,
			"EventGrades":          generateForm(eventSettingsClassGrades(event.Id, event.Grades)),
			"ListShooters":         event.Shooters,
			"EventId":              eventId,
		},
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := checkForm(event_add_shooterListForm("", []int{}).Inputs, r)
	eventId := validated_values["eventid"]
	http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)

	var new_shooter EventShooter
	new_shooter.Grade, _ = strconv.Atoi(validated_values["grade"])
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	var err error
	new_shooter.SID, err = strToInt(validated_values["sid"])
	if err == nil {
		temp_shooter := getShooterList(new_shooter.SID)
		new_shooter.FirstName = temp_shooter.NickName
		new_shooter.Surname = temp_shooter.Surname
		new_shooter.Club = temp_shooter.Club
	}
	eventShooterInsert(eventId, new_shooter)
}

func event_add_shooterForm(eventId string, grades []int) Form {
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
				Options:  AGE_GROUPS2(),
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

func event_add_shooterListForm(eventId string, grades []int) Form {
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
				Options:  AGE_GROUPS2(),
				Required: true,
			}, {
				Name:        "grade",
				Html:        "select",
				Label:       "Class & Grade",
				Placeholder: "Class & Grade",
				Required:    true,
				Options:     eventGradeOptions(grades),
			}, {
				Name: "sid",
				Html: "select",
				//				Label: "Class & Grade",
				//				Placeholder: "Class & Grade",
				Required: true,
				//				SelectedValues: eventGradeOptions(event),
			},
			//			"submit":Inputs{
			//				Html:      "submit",
			//				Value:   "Add Shooter",
			//			},
			{
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			},
		},
	}
}

func eventGradeOptions(eventGrades []int) []Option {
	//TODO add club defaults here
	var options []Option
	allGrades := grades()
	if len(eventGrades) < 1 {
		//Use system default event grades
		eventGrades = gradeList()
	}
	for _, gradeId := range eventGrades {
		options = append(options, Option{
			Value:   fmt.Sprintf("%v", gradeId),
			Display: allGrades[gradeId].LongName,
		})
	}
	return options
}
