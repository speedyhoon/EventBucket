package main

import (
	"fmt"
	"net/http"
	"strings"
	//	"bytes"
)

func startShooting(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_startShooting)
	templator(TEMPLATE_ADMIN, "start-shooting", startShooting_Data(data, false), w)
}

func startShootingAll(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_startShootingAll)
	templator(TEMPLATE_ADMIN, "start-shooting", startShooting_Data(data, true), w)
}
func startShooting_Data(data string, showAll bool) map[string]interface{} {
	arr := strings.Split(data, "/")
	event_id := arr[0]
	range_id := arr[1]
	event, _ := getEvent(event_id)

	available_class_shots := make([][]string, len(DEFAULT_CLASS_SETTINGS))
	html_available_class_shots := make([]string, len(DEFAULT_CLASS_SETTINGS))
	var sightersQty, shotsQty int
	var currentRangeClass RangeProperty
	var longest_shots_for_current_range int
	for index, classSetting := range DEFAULT_CLASS_SETTINGS{
		currentRangeClass = event.Ranges[range_id].Class[fmt.Sprintf("%v",index)]
		//If the range properties are set then use them to override the default shotsQty and sightersQty
		if currentRangeClass.ShotsQty > 0 || currentRangeClass.SightersQty > 0{
			sightersQty = currentRangeClass.SightersQty
			shotsQty = currentRangeClass.ShotsQty
		}else{
			sightersQty = classSetting.SightersQty
			shotsQty = classSetting.ShotsQty
		}
		if sightersQty + shotsQty > longest_shots_for_current_range{
			longest_shots_for_current_range = sightersQty + shotsQty
		}
		for i := 1; i <= sightersQty; i++ {
			available_class_shots[index] = append(available_class_shots[index], fmt.Sprintf("S%v",i))
			html_available_class_shots[index] += fmt.Sprintf("<td>S%v</td>",i)
		}
		for i := 1; i <= shotsQty; i++ {
			available_class_shots[index] = append(available_class_shots[index], fmt.Sprintf("%v",i))
			html_available_class_shots[index] += fmt.Sprintf("<td>%v</td>",i)
		}
	}

//	dump(available_class_shots)



//	available_class_shots := map[string][]string{
//		"F Class": []string{
//			"S1", "S2", "S3", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
//		},
//		"Target": []string{
//			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
//		},
//		"Match": []string{
//			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
//		},
//	}

	class_shots := map[string][]string{}
	class_shots_length := map[string]int{}
	var long_shots []string
	var temp_grade string
	var shooter_list []EventShooter
	for shooter_id, shooter_data := range event.Shooters {

		if showAll || (!showAll && ((event.IsPrizeMeet && len(shooter_data.Scores[range_id].Shots) <= 0)||(!event.IsPrizeMeet && shooter_data.Scores[range_id].Total <= 0))) {

			temp_grade = CLASS[shooter_data.Grade]
//			class_shots[temp_grade] = available_class_shots[temp_grade]
			class_shots[temp_grade] = available_class_shots[GRADE_TO_INT[shooter_data.Grade]]


			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club", "", -1)
			shooter_data.Id = shooter_id
			shooter_list = append(shooter_list, shooter_data)

		}


	}
	for temp_grade, shots_array := range class_shots {
		class_shots_length[temp_grade] = len(shots_array)
		if len(long_shots) < len(shots_array) {
			long_shots = shots_array
		}
	}
	first_class := ""
	first_class_int := 0
	for _, shooter := range event.Shooters {
		first_class = CLASS[shooter.Grade]
		first_class_int = GRADE_TO_INT[shooter.Grade]
		break
	}

	//Sort the list of shooters by grade only
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	name := func(c1, c2 *EventShooter) bool {
		return c1.FirstName < c2.FirstName
	}
	OrderedBy(grade, name).Sort(shooter_list, )


	var totalScores_link string
	if showAll{
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View Incompleted Shooters</a>", URL_startShooting, event_id, range_id)
	}else{
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View All Shooters</a>", URL_startShootingAll, event_id, range_id)
	}

	return map[string]interface{}{
		"Title":              "Start Shooting",
		"EventId":            event_id,
		"LinkToPage": 				totalScores_link,
		"RangeName":          event.Ranges[range_id].Name,
		"class_shots":        class_shots,
		"menu":               event_menu(event_id, event.Ranges, URL_startShooting, event.IsPrizeMeet),
		"RangeId":            range_id,
		"first_class":        first_class,
		"longest_shots":      long_shots,
		"class_shots_length": class_shots_length,
		//		"ListShooters": event.Shooters,
		"ListShooters": shooter_list,
		"Css":          "admin.css",
		"Js":           "start-shooting.js",
		"available_class_shots": available_class_shots,
		"first_class_shots": available_class_shots[first_class_int],
		"first_loaded_colspan": longest_shots_for_current_range - len(available_class_shots[first_class_int]) + 1,
		"target_heading_cells": html_available_class_shots[0],
		"fclass_heading_cells": html_available_class_shots[1],
		"match_heading_cells": html_available_class_shots[2],
	}
}

func updateShotScores(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(startShooting_Form("", "", "", "").Inputs, r) //total_scores_Form(event_id, range_id, shooter_id, shots
	event_id := validated_values["event_id"]
	event, _ := getEvent(event_id)
	range_id := validated_values["range_id"]
	if !event.Ranges[range_id].Locked {
		shooter_id := validated_values["shooter_id"]
		shots := validated_values["shots"]
		new_score := calc_total_centers(shots, GRADE_TO_INT[event.Shooters[shooter_id].Grade])
		if new_score.Centers > 0 {
			generator(w, fmt.Sprintf("%v.%v", new_score.Total, new_score.Centers), make(map[string]interface{}))
		} else {
			generator(w, fmt.Sprintf("%v", new_score.Total), make(map[string]interface{}))
		}
		event_total_score_update(event_id, range_id, shooter_id, new_score)
	} else {
		fmt.Println("BAD updateShotScores. Current Range is locked!")
	}
}

func calc_total_centers(shots string, class int) Score {
	//This function assumes all validation on input "shots" has at least been done!
	//AND input "shots" is verified to contain all characters in settings[class].validShots!
	//TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot

	total := 0
	centers := 0
	//	xs := 0
	countback1 := ""
	//	countback2 := ""

	relevant_settings := DEFAULT_CLASS_SETTINGS[class].ValidShots
	for _, shot := range strings.Split(shots[DEFAULT_CLASS_SETTINGS[class].SightersQty:], "") {
		total += relevant_settings[shot].Total
		centers += relevant_settings[shot].Centers
		countback1 = relevant_settings[shot].CountBack1 + countback1
		//		countback2 = relevant_settings[shot].CountBack2 + countback2
	}
	return Score{Total: total, Centers: centers, Shots: shots /*Xs: xs,*/, CountBack1: countback1 /*CountBack2: countback2*/}
}

func startShooting_Form(event_id, range_id, shooter_id, shots string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: map[string]Inputs{
			"shots": Inputs{
				Html:     "number",
				Label:    "Total",
				Required: true,
				Value:    shots,
				//TODO add min and max for validation on fclass and target
			},
			"shooter_id": Inputs{
				Html:  "hidden",
				Value: shooter_id,
			},
			"range_id": Inputs{
				Html:  "hidden",
				Value: range_id,
			},
			"event_id": Inputs{
				Html:  "hidden",
				Value: event_id,
			},
			"submit": Inputs{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}
