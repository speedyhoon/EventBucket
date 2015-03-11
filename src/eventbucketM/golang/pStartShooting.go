package main

import (
	"fmt"
	"mgo"
	"net/http"
	"strings"
)

func startShooting(data string) Page {
	return startShooting_Data(data, false)
}

func startShootingAll(data string) Page {
	return startShooting_Data(data, true)
}
func startShooting_Data(data string, showAll bool) Page {
	arr := strings.Split(data, "/")
	eventId := arr[0]
	rangeId, err := strToInt(arr[1])
	event, _ := getEvent(eventId)

	if event.Ranges[rangeId].Aggregate != "" || err != nil {
		return Page{
			TemplateFile: "start-shooting",
			Theme:        TEMPLATE_ADMIN,
			Title:        "Start Shooting",
			Data: M{
				"menu":                 eventMenu(eventId, event.Ranges, URL_startShooting, event.IsPrizeMeet),
				"target_heading_cells": "",
				"fclass_heading_cells": "",
				"match_heading_cells":  "",
				"Aggregate":            true,
			},
		}
	}

	available_class_shots := make([][]string, len(DEFAULT_CLASS_SETTINGS))
	html_available_class_shots := make([]string, len(DEFAULT_CLASS_SETTINGS))
	var sightersQty, shotsQty int
	var currentRangeClass RangeProperty
	var longest_shots_for_current_range int
	for index, classSetting := range DEFAULT_CLASS_SETTINGS {
		currentRangeClass = event.Ranges[rangeId].Class[fmt.Sprintf("%v", index)]
		//If the range properties are set then use them to override the default shotsQty and sightersQty
		if currentRangeClass.ShotsQty > 0 || currentRangeClass.SightersQty > 0 {
			sightersQty = currentRangeClass.SightersQty
			shotsQty = currentRangeClass.ShotsQty
		} else {
			sightersQty = classSetting.SightersQty
			shotsQty = classSetting.ShotsQty
		}
		if sightersQty+shotsQty > longest_shots_for_current_range {
			longest_shots_for_current_range = sightersQty + shotsQty
		}
		for i := 1; i <= sightersQty; i++ {
			available_class_shots[index] = append(available_class_shots[index], fmt.Sprintf("S%v", i))
			html_available_class_shots[index] += fmt.Sprintf("<td>S%v</td>", i)
		}
		for i := 1; i <= shotsQty; i++ {
			available_class_shots[index] = append(available_class_shots[index], fmt.Sprintf("%v", i))
			html_available_class_shots[index] += fmt.Sprintf("<td>%v</td>", i)
		}
	}

	class_shots := map[string][]string{}
	class_shots_length := map[string]int{}
	var long_shots []string
	var temp_grade Grade
	var shooter_list []EventShooter
	allGrades := grades()
	for shooterId, shooter_data := range event.Shooters {
		if showAll || (!showAll && ((event.IsPrizeMeet && len(shooter_data.Scores[fmt.Sprintf("%v", rangeId)].Shots) <= 0) || (!event.IsPrizeMeet && shooter_data.Scores[fmt.Sprintf("%v", rangeId)].Total <= 0))) {
			temp_grade = allGrades[shooter_data.Grade]
			class_shots[temp_grade.ClassName] = available_class_shots[temp_grade.ClassId]
			//TODO add ignore case here!!!!!!!!
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club Inc", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club.", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Rifle Club", "", -1)
			shooter_data.Club = strings.Replace(shooter_data.Club, " Ex-Services Memorial", "", -1)
			shooter_data.Id = shooterId
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
		first_class = allGrades[shooter.Grade].ClassName
		first_class_int = allGrades[shooter.Grade].ClassId
		break
	}

	//Sort the list of shooters by grade only
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	name := func(c1, c2 *EventShooter) bool {
		return c1.FirstName < c2.FirstName
	}
	OrderedBy(grade, name).Sort(shooter_list)

	var totalScores_link string
	if showAll {
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View Incompleted Shooters</a>", URL_startShooting, eventId, rangeId)
	} else {
		totalScores_link = fmt.Sprintf("<a href=%v/%v/%v>View All Shooters</a>", URL_startShootingAll, eventId, rangeId)
	}

	return Page{
		TemplateFile: "start-shooting",
		Theme:        TEMPLATE_ADMIN,
		Title:        "Start Shooting",
		Data: M{
			"EventId":            eventId,
			"LinkToPage":         totalScores_link,
			"RangeName":          event.Ranges[rangeId].Name,
			"class_shots":        class_shots,
			"menu":               eventMenu(eventId, event.Ranges, URL_startShooting, event.IsPrizeMeet),
			"strRangeId":         fmt.Sprintf("%v", rangeId),
			"RangeId":            rangeId,
			"first_class":        first_class,
			"longest_shots":      long_shots,
			"class_shots_length": class_shots_length,
			"ListShooters":       shooter_list,
			"Js":                 "start-shooting.js",
			"available_class_shots": available_class_shots,
			"first_class_shots":     available_class_shots[first_class_int],
			"first_loaded_colspan":  longest_shots_for_current_range - len(available_class_shots[first_class_int]) + 1,
			"target_heading_cells":  html_available_class_shots[0],
			"fclass_heading_cells":  html_available_class_shots[1],
			"match_heading_cells":   html_available_class_shots[2],
			"Aggregate":             false,
		},
	}
}

func updateShotScores(w http.ResponseWriter, r *http.Request) {
	validated_values, passed := valid8(startShooting_Form("", "", "", "").Inputs, r)
	if passed {
		event := validated_values["event"].(Event)
		rangeId := validated_values["rangeid"].(int)
		//TODO check the range exists (is not nill) before accessing Locked or Aggregate
		if !event.Ranges[rangeId].Locked && event.Ranges[rangeId].Aggregate == "" {
			eventId := validated_values["eventid"].(string)
			shooterId := validated_values["shooterid"].(int)
			shots := validated_values["shots"].(string)

			new_score := calc_total_centers(shots, grades()[event.Shooters[shooterId].Grade].ClassId)
			var temp Page
			if new_score.Centers > 0 {
				generator(w, fmt.Sprintf("%v.%v", new_score.Total, new_score.Centers), temp)
			} else {
				generator(w, fmt.Sprintf("%v", new_score.Total), temp)
			}
			//Add any linked shooters to this update
			shooterIds := []int{shooterId}
			if event.Shooters[shooterId].LinkedId != nil {
				shooterIds = append(shooterIds, *event.Shooters[shooterId].LinkedId)
			}
			//Find all the aggs that this rangeId is in
			aggsFound := searchForAggs(event.Ranges, rangeId)
			var updateBson = make(M)
			if len(aggsFound) > 0 {
				updateBson = calculateAggs(event.Shooters[shooterId].Scores, aggsFound, shooterIds, event.Ranges)
			}
			//			event/shooters/shooterid/rangeId

			updateBson[Dot(schemaSHOOTER, shooterId, rangeId)] = new_score
			if event.Shooters[shooterId].LinkedId != nil {
				updateBson[Dot(schemaSHOOTER, event.Shooters[shooterId].LinkedId, rangeId)] = new_score
			}
			//			eventTotalScoreUpdate(eventId, rangeId, shooterIds, new_score)
			//			updateSetter := make(M)
			//			for _, shooterId := range shooterIds{
			//				updateSetter[Dot(schemaSHOOTER, shooterId, rangeId)] = score
			//			}
			change := mgo.Change{
				Upsert: true,
				Update: M{
					"$set": updateBson,
				},
			}
			var event Event
			_, err := conn.C(TBLevent).FindId(eventId).Apply(change, &event)
			//TODO better error handling would be nice
			if err != nil {
				Warning.Println(err)
			}
		} else {
			Warning.Println("BAD updateShotScores. Current Range is locked or is an aggreate range.")
		}
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
	//TODO check the class int is valid
	if class >= 0 && class < len(DEFAULT_CLASS_SETTINGS) {
		relevant_settings := DEFAULT_CLASS_SETTINGS[class].ValidShots
		for _, shot := range strings.Split(shots[DEFAULT_CLASS_SETTINGS[class].SightersQty:], "") {
			total += relevant_settings[shot].Total
			centers += relevant_settings[shot].Centers
			countback1 = relevant_settings[shot].CountBack1 + countback1
			//		countback2 = relevant_settings[shot].CountBack2 + countback2
		}
		/////////////////////////////		return Score{Total: total, Centers: centers, Shots: shots /*Xs: xs,*/, CountBack1: countback1 /*CountBack2: countback2*/}
	}
	return Score{}
}

func startShooting_Form(eventId, rangeId, shooterId, shots string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: []Inputs{
			{
				Name:      "eventid",
				Html:      "hidden",
				Value:     eventId,
				VarType:   "string",
				VarMaxLen: 999,
				VarMinLen: 1,
			},
			{
				Name:      "rangeid",
				Html:      "hidden",
				Value:     rangeId,
				VarType:   "int",
				VarMaxLen: 999, //TODO needs better parameters
				VarMinLen: 0,
			},
			{
				Name:    "shooterid",
				Html:    "hidden",
				Value:   shooterId,
				VarType: "int",
			},
			{
				Name:     "shots",
				Html:     "number",
				Label:    "Total",
				Required: true,
				Value:    shots,
				//TODO add min and max for validation on fclass and target
				VarType:   "string",
				VarMaxLen: 90, //TODO make dynamic by getting the shooters class
				VarMinLen: 1,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}