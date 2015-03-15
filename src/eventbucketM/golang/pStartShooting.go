package main

import (
	"fmt"
	"net/http"
	"strconv"
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

	linkName := "All"
	pageLink := URL_startShootingAll
	if showAll {
		linkName = "Incompleted"
		pageLink = URL_startShooting
	}

	return Page{
		TemplateFile: "start-shooting",
		Theme:        TEMPLATE_ADMIN,
		Title:        "Start Shooting",
		Data: M{
			"EventId":            eventId,
			"pageLink":           pageLink,
			"linkName":           linkName,
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

func updateShotScores2(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(startShootingForm("", "", "", "").Inputs, r)
	eventId := validatedValues["eventid"]
	rangeId, rangeErr := strconv.Atoi(validatedValues["rangeid"])
	shooterId, shooterErr := strconv.Atoi(validatedValues["shooterid"])
	event, eventErr := getEvent(eventId)
	//Check the score can be saved
	if rangeErr != nil || shooterErr != nil || eventErr != nil || rangeId >= len(event.Ranges) || event.Ranges[rangeId].Locked || event.Ranges[rangeId].IsAgg {
		http.NotFound(w, r)
		return
	}
	newScore := calcTotalCentres(validatedValues["shots"], grades()[event.Shooters[shooterId].Grade].ClassId)
	//Return the score to the client
	if newScore.Centers > 0 {
		fmt.Fprintf(w, "%v.%v", newScore.Total, newScore.Centers)
	} else {
		fmt.Fprintf(w, "%v", newScore.Total)
	}
	go updateShootersScores(newScore, shooterId, rangeId, event)
}

func updateShootersScores(newScore Score, shooterId, rangeId int, event Event) {
	updateBson := M{Dot(schemaSHOOTER, shooterId, rangeId): newScore}
	shooterIds := []int{shooterId}
	//Add any linked shooters to this update
	if event.Shooters[shooterId].LinkedId != nil {
		shooterIds = append(shooterIds, *event.Shooters[shooterId].LinkedId)
		updateBson[Dot(schemaSHOOTER, event.Shooters[shooterId].LinkedId, rangeId)] = newScore
	}
	//Find all the aggs that this rangeId is in & update the scores
	aggsFound := searchForAggs(event.Ranges, rangeId)
	if len(aggsFound) > 0 {
		for index, data := range calculateAggs(event.Shooters[shooterId].Scores, aggsFound, shooterIds, event.Ranges) {
			updateBson[index] = data
		}
	}
	tableUpdateData(TBLevent, event.Id, updateBson)
}

//This function assumes all validation on input "shots" has at least been done!
//AND input "shots" is verified to contain all characters in settings[class].validShots!
func calcTotalCentres(shots string, class int) Score {
	//TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot
	var total, centres int
	var countBack string
	if class >= 0 && class < len(DEFAULT_CLASS_SETTINGS) {
		classRules := DEFAULT_CLASS_SETTINGS[class].ValidShots
		//Ignore the first sighter shots from being added to the total score. Unused sighters should be still be present in the data passed
		for _, shot := range strings.Split(shots[DEFAULT_CLASS_SETTINGS[class].SightersQty:], "") {
			total += classRules[shot].Total
			centres += classRules[shot].Centers
			countBack = classRules[shot].CountBack1 + countBack
		}
	}
	return Score{Total: total, Centers: centres, Shots: shots, CountBack1: countBack}
}

func startShootingForm(eventId, rangeId, shooterId, shots string) Form {
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
			}, {
				Name:      "rangeid",
				Html:      "hidden",
				Value:     rangeId,
				VarType:   "int",
				VarMaxLen: 999, //TODO needs better parameters
				VarMinLen: 0,
			}, {
				Name:    "shooterid",
				Html:    "hidden",
				Value:   shooterId,
				VarType: "int",
			}, {
				Name:     "shots",
				Html:     "number",
				Required: true,
				Value:    shots,
				//TODO add min and max for validation on fclass and target
				VarType:   "string",
				VarMaxLen: 90, //TODO make dynamic by getting the shooters class
				VarMinLen: 1,
			},
		},
	}
}
