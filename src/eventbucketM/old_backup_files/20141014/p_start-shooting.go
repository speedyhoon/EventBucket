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

	available_class_shots := map[string][]string{
		"F Class": []string{
			"S1", "S2", "S3", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
		},
		"Target": []string{
			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		},
		"Match": []string{
			"S1", "S2", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		},
	}

	class_shots := map[string][]string{}
	class_shots_length := map[string]int{}
	var long_shots []string
	var temp_grade string
	var shooter_list []EventShooter
	for shooter_id, shooter_data := range event.Shooters {

		if showAll || (!showAll && ((event.IsPrizeMeet && len(shooter_data.Scores[range_id].Shots) <= 0)||(!event.IsPrizeMeet && shooter_data.Scores[range_id].Total <= 0))) {

			temp_grade = CLASS[shooter_data.Grade]
			class_shots[temp_grade] = available_class_shots[temp_grade]


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
	for _, shooter := range event.Shooters {
		first_class = CLASS[shooter.Grade]
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
		"menu":               event_menu(event_id, event.Ranges, URL_startShooting),
		"RangeId":            range_id,
		"first_class":        first_class,
		"longest_shots":      long_shots,
		"class_shots_length": class_shots_length,
		//		"ListShooters": event.Shooters,
		"ListShooters": shooter_list,
		"Css":          "admin.css",
		"Js":           "start-shooting.js",
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
		//TODO change this to use grade to Int
		new_score := calc_total_centers(shots, 0)
		dump(new_score)
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
	//TODO need validation to check that the shots given match the required vailation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot

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
				Label: "Save",
			},
		},
	}
}
