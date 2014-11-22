package main

import (
	"net/http"
	"strings"
)
func shooterInsert(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(event_add_shooterForm(event).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.FirstName = validated_values["first"]
	new_shooter.Surname = validated_values["surname"]
	new_shooter.Club = validated_values["club"]
	new_shooter.Grade = validated_values["grade"]
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	event_shooter_insert(event_id, new_shooter)
}

func event(w http.ResponseWriter, r *http.Request) {
	event_id := get_id_from_url(r, "/event/")
	templator(TEMPLATE_ADMIN, "event", event_Data(event_id), w)
}

func event_Data(event_id string) map[string]interface{} {
	event, err := getEvent(event_id)
	if err{
		return map[string]interface{}{
			"Title": "Event not found",
			"Menu":  standard_menu(ORGANISERS_MENU_ITEMS),
			"Valid": false,
		}
	}
	return map[string]interface{}{
		"Title": event.Name,
		"EventId": event_id,
		"ListRanges": event.Ranges,
//		"AddShooter": generateForm2(event_add_shooterForm(event)),
		"ListShooters": event.Shooters,
		"Menu": event_menu(event_id, event.Ranges, URL_event, event.IsPrizeMeet),
		"AddRange": generateForm2(eventSettings_add_rangeForm(event_id)),
		"ExistingShooterEntry": URL_shooterListInsert,
		"NewShooterEntry": URL_shooterInsert,
//		"GradeOptions": build_options(available_classes_grades(event)),
		"GradeOptions": draw_options(Inputs{Options:available_classes_grades(event)}, ""),
		//TODO add ClubOptions when club textbox is changed to a datalist
//		"AgeOptions": build_options(AGE_GROUPS2),
		"AgeOptions": draw_options(Inputs{Options:AGE_GROUPS2}, ""),
		"Valid": true,

		"EventGrades":    generateForm2(eventSettings_class_grades(event)),
	}
}

func shooterListInsert(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(event_add_shooterListForm(event).Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_event + event_id, w, r)
	var new_shooter EventShooter
	new_shooter.SID = str_to_int(validated_values["sid"])

	//	new_shooter.FirstName = validated_values["first"]
	//	new_shooter.Surname = validated_values["surname"]
	//	new_shooter.Club = validated_values["club"]
	new_shooter.Grade = validated_values["grade"]
	if validated_values["age"] != "" {
		new_shooter.AgeGroup = validated_values["age"]
	}
	temp_shooter := getShooterList(new_shooter.SID)
	new_shooter.FirstName = temp_shooter.NickName
	new_shooter.Surname = temp_shooter.Surname
	new_shooter.Club = temp_shooter.Club
	event_shooter_insert(event_id, new_shooter)
}

func event_add_shooterForm(event Event) Form {

	var event_id string
	if event.Id != ""{
		event_id = event.Id
	}

	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
			"first":Inputs{
				Html:      "text",
				Label:   "First Name",
				Required: true,
			},
			"surname":Inputs{
				Html:      "text",
				Label:   "Surname",
				Required: true,
			},
			"club":Inputs{
				Html:      "text",
				//TODO change club to a data-list
				//SelectValues:   getClubSelectBox(eventsCollection),
				Label:   "Club",
				Required: true,
			},
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: available_classes_grades(event),
			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Add Shooter",
//			},
			"event_id":Inputs{
				Html: "hidden",
				Value: event_id,
			},
		},
	}
}

func event_add_shooterListForm(event Event) Form {

	var event_id string
	if event.Id != ""{
		event_id = event.Id
	}

	return Form{
		Action: URL_shooterInsert,
		Title: "Add Shooters",
		Inputs: map[string]Inputs{
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
			"age":Inputs{
				Html:      "select",
				Label: "Age Group",
				Options: AGE_GROUPS2,
				Required: true,
			},
			"grade":Inputs{
				Html:      "select",
				Label: "Class & Grade",
				Placeholder: "Class & Grade",
				Required: true,
				Options: available_classes_grades(event),
			},
			"sid":Inputs{
				Html:      "select",
//				Label: "Class & Grade",
//				Placeholder: "Class & Grade",
				Required: true,
//				SelectedValues: available_classes_grades(event),
			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Add Shooter",
//			},
			"event_id":Inputs{
				Html: "hidden",
				Value: event_id,
			},
		},
	}
}

func available_classes_grades(event Event)[]Option{
	var grades []Option
	var grade_list map[string]bool
	selected_grades := strings.Split(event.Grades, ",")
	no_grades_selected := event.Grades == ""
	if !no_grades_selected {
		grade_list = slice_to_map_bool(selected_grades)
	}

	for _, class_settings := range DEFAULT_CLASS_SETTINGS {
		for _, grade_id := range class_settings.Grades {
			if grade_list[grade_id] || no_grades_selected {
				grades = append(grades, Option{
						Value: grade_id,
						Display: CLASSES[grade_id],
					})
			}
		}
	}
	return grades
}
