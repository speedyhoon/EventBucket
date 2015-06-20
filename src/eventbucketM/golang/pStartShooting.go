package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func startShooting(data string) Page {
	return startShootingData(data, false)
}

func startShootingAll(data string) Page {
	return startShootingData(data, true)
}
func startShootingData(data string, showAll bool) Page {
	arr := strings.Split(data, "/")
	eventID := arr[0]
	var titleAll string
	var rangeID int
	var err error
	if len(arr) >= 2 {
		rangeID, err = strconv.Atoi(arr[1])
	} else {
		return Page{
			TemplateFile: "start-shooting",
			Theme:        templateAdmin,
			Title:        "Start Shooting" + titleAll,
			Data: M{
				"menu":                 "",
				"target_heading_cells": "",
				"fclass_heading_cells": "",
				"match_heading_cells":  "",
				"Aggregate":            true,
			},
		}
	}
	event, _ := getEvent(eventID)
	if showAll {
		titleAll = " Show All"
	}
	if event.Ranges == nil {
		return Page{
			TemplateFile: "start-shooting",
			Theme:        templateAdmin,
			Title:        "Start Shooting" + titleAll,
			Data: M{
				"menu":                 eventMenu(eventID, event.Ranges, urlStartShooting, event.IsPrizeMeet),
				"target_heading_cells": "",
				"fclass_heading_cells": "",
				"match_heading_cells":  "",
				"Aggregate":            true,
			},
		}
	}

	if event.Ranges[rangeID].Aggregate != "" || err != nil {
		return Page{
			TemplateFile: "start-shooting",
			Theme:        templateAdmin,
			Title:        "Start Shooting" + titleAll,
			Data: M{
				"menu":                 eventMenu(eventID, event.Ranges, urlStartShooting, event.IsPrizeMeet),
				"target_heading_cells": "",
				"fclass_heading_cells": "",
				"match_heading_cells":  "",
				"Aggregate":            true,
			},
		}
	}

	availableClassShots := make([][]string, len(defaultClassSettings))
	htmlAvailableClassShots := make([]string, len(defaultClassSettings))
	var sightersQty, shotsQty int
	var currentRangeClass RangeProperty
	var longestShotsForCurrentRange int
	for index, classSetting := range defaultClassSettings {
		currentRangeClass = event.Ranges[rangeID].Class[fmt.Sprintf("%v", index)]
		//If the range properties are set then use them to override the default shotsQty and sightersQty
		if currentRangeClass.ShotsQty > 0 || currentRangeClass.SightersQty > 0 {
			sightersQty = currentRangeClass.SightersQty
			shotsQty = currentRangeClass.ShotsQty
		} else {
			sightersQty = classSetting.SightersQty
			shotsQty = classSetting.ShotsQty
		}
		if sightersQty+shotsQty > longestShotsForCurrentRange {
			longestShotsForCurrentRange = sightersQty + shotsQty
		}
		for i := 1; i <= sightersQty; i++ {
			availableClassShots[index] = append(availableClassShots[index], fmt.Sprintf("S%v", i))
			htmlAvailableClassShots[index] += fmt.Sprintf("<td>S%v", i)
		}
		for i := 1; i <= shotsQty; i++ {
			availableClassShots[index] = append(availableClassShots[index], fmt.Sprintf("%v", i))
			htmlAvailableClassShots[index] += fmt.Sprintf("<td>%v", i)
		}
	}

	classShots := map[string][]string{}
	classShotsLength := map[string]int{}
	var longShots []string
	var tempGrade Grade
	var shooterList []EventShooter
	allGrades := grades()
	for shooterID, shooterData := range event.Shooters {
		if showAll || (!showAll && ((event.IsPrizeMeet && len(shooterData.Scores[fmt.Sprintf("%v", rangeID)].Shots) <= 0) || (!event.IsPrizeMeet && shooterData.Scores[fmt.Sprintf("%v", rangeID)].Total <= 0))) {
			tempGrade = allGrades[shooterData.Grade]
			classShots[tempGrade.ClassName] = availableClassShots[tempGrade.ClassID]
			//TODO add ignore case here!!!!!!!!
			shooterData.Club = strings.Replace(shooterData.Club, " Rifle Club Inc.", "", -1)
			shooterData.Club = strings.Replace(shooterData.Club, " Rifle Club Inc", "", -1)
			shooterData.Club = strings.Replace(shooterData.Club, " Rifle Club.", "", -1)
			shooterData.Club = strings.Replace(shooterData.Club, " Rifle Club", "", -1)
			shooterData.Club = strings.Replace(shooterData.Club, " Ex-Services Memorial", "", -1)
			shooterData.ID = shooterID
			shooterList = append(shooterList, shooterData)
		}
	}
	for tempGrade2, shotsArray := range classShots {
		classShotsLength[tempGrade2] = len(shotsArray)
		if len(longShots) < len(shotsArray) {
			longShots = shotsArray
		}
	}
	var firstClass string
	var firstClassInt int
	for _, shooter := range event.Shooters {
		firstClass = allGrades[shooter.Grade].ClassName
		firstClassInt = allGrades[shooter.Grade].ClassID
		break
	}

	//Sort the list of shooters by grade, then first name
	orderShooters("", sorterGrade, sorterFirstName)

	linkName := "All"
	pageLink := urlStartShootingAll
	if showAll {
		linkName = "Incompleted"
		pageLink = urlStartShooting
	}

	return Page{
		TemplateFile: "start-shooting",
		Theme:        templateAdmin,
		Title:        "Start Shooting" + titleAll,
		Data: M{
			"EventID":            eventID,
			"pageLink":           pageLink,
			"linkName":           linkName,
			"RangeName":          event.Ranges[rangeID].Name,
			"class_shots":        classShots,
			"menu":               eventMenu(eventID, event.Ranges, urlStartShooting, event.IsPrizeMeet),
			"strRangeID":         fmt.Sprintf("%v", rangeID),
			"RangeID":            rangeID,
			"first_class":        firstClass,
			"longest_shots":      longShots,
			"class_shots_length": classShotsLength,
			"ListShooters":       shooterList,
			"Js":                 "start-shooting.js",
			"available_class_shots": availableClassShots,
			"first_class_shots":     availableClassShots[firstClassInt],
			"first_loaded_colspan":  longestShotsForCurrentRange - len(availableClassShots[firstClassInt]) + 1,
			"target_heading_cells":  htmlAvailableClassShots[0],
			"fclass_heading_cells":  htmlAvailableClassShots[1],
			"match_heading_cells":   htmlAvailableClassShots[2],
			"Aggregate":             false,
		},
	}
}

func hasShootFinished(shots string, grade int) bool {
	classSettings := defaultClassSettings[grade2Class(grade)]
	return len(strings.Replace(shots[classSettings.SightersQty:], "-", "", -1)) == classSettings.ShotsQty
}

func updateShotScores2(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(startShootingForm("", "", "", "").inputs, r)
	rangeID, rangeErr := strconv.Atoi(validatedValues["rangeid"])
	shooterID, shooterErr := strconv.Atoi(validatedValues["shooterid"])
	event, eventErr := getEvent(validatedValues["eventid"])

	//Check the score can be saved
	if rangeErr != nil || shooterErr != nil || eventErr != nil || rangeID >= len(event.Ranges) || event.Ranges[rangeID].Locked || event.Ranges[rangeID].IsAgg {
		http.NotFound(w, r)
		return
	}

	//Calculate the score based on the shots given
	newScore := calcTotalCentres(validatedValues["shots"], grades()[event.Shooters[shooterID].Grade].ClassID)

	//Return the score to the client
	if newScore.Centres > 0 {
		fmt.Fprintf(w, "%v<sup>%v</sup>", newScore.Total, newScore.Centres)
	} else {
		fmt.Fprintf(w, "%v", newScore.Total)
	}

	go triggerScoreCalculation(newScore, rangeID, event.Shooters[shooterID], event)
}

//This function assumes all validation on input "shots" has at least been done!
//AND input "shots" is verified to contain all characters in settings[class].validShots!
func calcTotalCentres(shots string, classID int) Score {
	//TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot
	var total, centres, warning int
	var countBack string
	if classID >= 0 && classID < len(defaultClassSettings) {

		//Ignore the first sighter shots from being added to the total score. Unused sighters should be still be present in the data passed
		for _, shot := range strings.Split(shots[defaultClassSettings[classID].SightersQty:], "") {
			total += defaultClassSettings[classID].ValidShots[shot].Total
			centres += defaultClassSettings[classID].ValidShots[shot].Centres

			//Append count back in reverse order so it can be ordered by the last few shots
			countBack = defaultClassSettings[classID].ValidShots[shot].CountBack + countBack
			if shot == "-" {
				warning = legendIncompleteScore
			}
		}
	}
	if warning == 0 && isScoreHighestPossibleScore(classID, 10, total, centres) {
		warning = legendHighestPossibleScore
	}
	return Score{Total: total, Centres: centres, Shots: shots, CountBack: countBack, Warning: warning}
}

func startShootingForm(eventID, rangeID, shooterID, shots string) Form {
	temp1 := v8MaxIntegerID
	temp2 := v8MinIntegerID
	return Form{
		action: urlUpdateTotalScores,
		inputs: []Inputs{
			{
				name:      "eventid",
				html:      "hidden",
				value:     eventID,
				varType:   "string",
				maxLength: v8MaxEventID,
				minLength: v8MinEventID,
			}, {
				name:    "rangeid",
				html:    "hidden",
				value:   rangeID,
				varType: "int",
				max:     &temp1,
				min:     &temp2,
			}, {
				name:    "shooterid",
				html:    "hidden",
				value:   shooterID,
				varType: "int",
				max:     &temp1,
				min:     &temp2,
			}, {
				name:     "shots",
				html:     "number",
				required: true,
				value:    shots,
				//TODO add min and max for validation on fclass and target
				varType:   "string",
				maxLength: v8MinShots, //TODO make dynamic by getting the shooters class
				minLength: v8Minhots,
			},
		},
	}
}
