package main

import (
	"net/http"
	"fmt"
	"strconv"
)
func shooterInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(event_add_shooterForm("", []int{}).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.FirstName = validated_values["first"]
	new_shooter.Surname = validated_values["surname"]
	new_shooter.Club = validated_values["club"]
	new_shooter.Grade, _ = strconv.Atoi(validated_values["grade"])
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	event_shooter_insert(event_id, new_shooter)
}

func event(event_id string) M {
	event, err := getEvent(event_id)
	if err{
		return M{
			"Title": "Event not found: "+event_id,
			"Menu":  standard_menu(ORGANISERS_MENU_ITEMS),
			"Valid": false,
		}
	}
	return M{
		"Title": event.Name,
		"EventId": event_id,
		"ListRanges": event.Ranges,
		"ListShooters": event.Shooters,
		"Menu": event_menu(event_id, event.Ranges, URL_event, event.IsPrizeMeet),
		"AddRange": generateForm2(eventSettings_add_rangeForm(event_id)),
		"ExistingShooterEntry": URL_shooterListInsert,
		"NewShooterEntry": URL_shooterInsert,
		"GradeOptions": draw_options(Inputs{Options:eventGradeOptions(event.Grades)}, ""),
		//TODO add ClubOptions when club textbox is changed to a datalist
		"AgeOptions": draw_options(Inputs{Options:AGE_GROUPS2()}, ""),
		"Valid": true,
		"EventGrades":    generateForm2(eventSettingsClassGrades(event.Id, event.Grades)),
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(event_add_shooterListForm("", []int{}).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.Grade, _ = strconv.Atoi(validated_values["grade"])
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	var success bool
	new_shooter.SID, success = strToInt(validated_values["sid"])
	if success {
		temp_shooter := getShooterList(new_shooter.SID)
		new_shooter.FirstName = temp_shooter.NickName
		new_shooter.Surname = temp_shooter.Surname
		new_shooter.Club = temp_shooter.Club
	}
	event_shooter_insert(event_id, new_shooter)
}

func event_add_shooterForm(eventId string, grades []int) Form {
	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: []Inputs{
			{
				Name: "first",
				Html:      "text",
				Label:   "First Name",
				Required: true,
			},
			{
				Name: "surname",
				Html:      "text",
				Label:   "Surname",
				Required: true,
			},
			{
				Name: "club",
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
				Required: true,
			},
			{
				Name: "age",
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2(),
				Required: true,
			},
			{
				Name: "grade",
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: eventGradeOptions(grades),
			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Add Shooter",
//			},
			{
				Name: "event_id",
				Html: "hidden",
				Value: eventId,
			},
		},
	}
}

func event_add_shooterListForm(eventId string, grades []int) Form {
	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: []Inputs{
//			"first":Inputs{
//				Html:      "text",
//				Label:   "First Name",
//				Required: true,
//			},
//			"surname":Inputs{
//				Html:      "text",
//				Label:   "Surname",
//				Required: true,
//			},
//			"club":Inputs{
//				Html:      "text",
//				TODO change club to a data-list
//				//SelectValues:   getClubSelectBox(eventsCollection),
//				Label:   "Club",
//				Required: true,
//			},
			{
				Name: "age",
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2(),
				Required: true,
			},
			{
				Name: "grade",
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: eventGradeOptions(grades),
			},
			{
				Name: "sid",
				Html:      "select",
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
				Name: "event_id",
				Html: "hidden",
				Value: eventId,
			},
		},
	}
}

func eventGradeOptions(eventGrades []int)[]Option{
	//TODO add club defaults here
	var options []Option
	allGrades := grades()
	if len(eventGrades) < 1 {
		//Use system default event grades
		eventGrades = gradeList()
	}
	for _, gradeId := range eventGrades{
		options = append(options, Option{
			Value: fmt.Sprintf("%v", gradeId),
			Display: allGrades[gradeId].LongName,
		})
	}
	return options
}
