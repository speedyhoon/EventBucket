package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func totalScores(data string) Page {
	return totalScoresData(data, false)
}

func totalScoresAll(data string) Page {
	return totalScoresData(data, true)
}

func totalScoresData(data string, showAll bool) Page {
	arr := strings.Split(data, "/")
	eventID := arr[0]
	var titleAll string
	emptyPage := Page{
		TemplateFile: "total-scores",
		Title:        "Total Scores" + titleAll,
		Theme:        templateAdmin,
		Data: M{
			"Title":      "Total Scores",
			"LinkToPage": "",
			"EventID":    "",
			"RangeName":  "",
			"Message":    errorEnterScoresInAgg,
			"menu":       "",
		},
	}
	if !(len(arr) >= 2) {
		return emptyPage
	}
	rangeID := arr[1]
	strRangeID, err := strconv.Atoi(rangeID)
	event, eventMissing := getEvent(eventID)
	if event.Ranges == nil {
		return emptyPage
	}
	currentRange := event.Ranges[strRangeID]
	if showAll {
		titleAll = " Show All"
	}
	if err != nil || eventMissing != nil || currentRange.Locked || currentRange.Aggregate != "" {
		return emptyPage
	}
	var totalScoresLink string
	if showAll {
		totalScoresLink = fmt.Sprintf("<a href=%v%v>View Incompleted Shooters</a>", urlTotalScores, data)
	} else {
		totalScoresLink = fmt.Sprintf("<a href=%v%v>View All Shooters</a>", urlTotalScoresAll, data)
	}

	if len(currentRange.Aggregate) > 0 {
		return Page{
			TemplateFile: "total-scores",
			Title:        "Total Scores" + titleAll,
			Theme:        templateAdmin,
			Data: M{
				"LinkToPage": totalScoresLink,
				"EventID":    eventID,
				"RangeName":  currentRange.Name,
				"Message":    errorEnterScoresInAgg,
				"menu":       eventMenu(eventID, event.Ranges, urlTotalScores, event.IsPrizeMeet),
			},
		}
	}

	//Sort the list of shooters by grade only
	grade := func(c1, c2 *EventShooter) bool {
		return c1.Grade < c2.Grade
	}
	name := func(c1, c2 *EventShooter) bool {
		return c1.FirstName < c2.FirstName
	}

	var shooterList2 []EventShooter
	var shootersForms2 []string
	for shooterID, shooterData2 := range event.Shooters {
		var score string
		if shooterData2.Scores[rangeID].Total > 0 {
			score = fmt.Sprintf("%v", shooterData2.Scores[rangeID].Total)
		}
		if shooterData2.Scores[rangeID].Centres > 0 {
			score += fmt.Sprintf(".%v", shooterData2.Scores[rangeID].Centres)
		}
		if showAll || (!showAll && score == "") {
			shootersForms2 = append(shootersForms2, generateForm(totalScoresForm(eventID, rangeID, shooterID, score)))
			shooterData2.ID = shooterID
			shooterList2 = append(shooterList2, shooterData2)
		}
	}

	orderedBy(grade, name).Sort(shooterList2)

	return Page{
		TemplateFile: "total-scores",
		Title:        "Total Scores" + titleAll,
		Theme:        templateAdmin,
		Data: M{
			"LinkToPage": totalScoresLink,
			"EventID":    eventID,
			"RangeName":  currentRange.Name,
			"RangeID":    rangeID,
			"ListRanges": event.Ranges,
			//"ListShooters": event.Shooters,
			"ListShooters":    shooterList2,
			"menu":            eventMenu(eventID, event.Ranges, urlTotalScores, event.IsPrizeMeet),
			"FormTotalScores": shootersForms2,
			"Js":              "total-scores.js",
		},
	}
}

func updateTotalScores(w http.ResponseWriter, r *http.Request) {
	//	validatedValues := checkForm(totalScoresForm("", "", "", "").inputs, r) //totalScoresForm(eventID, rangeID, shooterID, total, centres
	validatedValues, passed := valid8(totalScoresForm("", "", -1, "").inputs, r)
	if passed {
		eventID := validatedValues["eventid"].(string)
		event := validatedValues["event"].(Event)
		rangeID := validatedValues["rangeid"].(int)
		shooterID := validatedValues["shooterid"].(int)
		score := strings.Split(validatedValues["score"].(string), ".")
		total, err := strconv.Atoi(score[0])
		if total > 0 && err == nil {
			newScore := Score{Total: total}
			var centres int
			centres, err = strconv.Atoi(score[1])
			if len(score) > 1 && score[1] != "" && centres > 0 && err == nil {
				centres, _ = strconv.Atoi(score[1])
				newScore.Centres = centres
			}
			shooterIDs := []int{shooterID}
			if event.Shooters[shooterID].LinkedID != nil {
				shooterIDs = append(shooterIDs, *event.Shooters[shooterID].LinkedID)
			}
			go eventTotalScoreUpdate(eventID, rangeID, shooterIDs, newScore)
		}
		http.Redirect(w, r, fmt.Sprintf("%v%v", urlTotalScores+eventID+"/", rangeID), http.StatusSeeOther)
	}
}

func searchForAggs(ranges []Range, rangeID int) []int {
	var aggsFound []int
	var foundRangeID int
	var err error
	for indexRangeID, rangeObj := range ranges {
		if len(rangeObj.Aggregate) > 0 {
			for _, strRangeID := range strings.Split(rangeObj.Aggregate, ",") {
				foundRangeID, err = strconv.Atoi(fmt.Sprintf("%v", strRangeID))
				if err == nil && foundRangeID == rangeID {
					aggsFound = append(aggsFound, indexRangeID)
				}
			}
		}
	}
	return aggsFound
}

func calculateAggs(shooterScores map[string]Score, ranges []int, shooterIDs []int, eventRanges []Range) M {
	if shooterScores == nil {
		return M{}
	}
	var total, centres int
	var countBack string
	updateBson := make(M)
	for _, aggID := range ranges {
		total = 0
		centres = 0
		for _, rangeID := range strings.Split(eventRanges[aggID].Aggregate, ",") {
			total += shooterScores[rangeID].Total
			centres += shooterScores[rangeID].Centres
			countBack = shooterScores[rangeID].CountBack
		}
		for _, shooterID := range shooterIDs {
			updateBson[dot(schemaSHOOTER, shooterID, aggID)] = Score{Total: total, Centres: centres, CountBack: countBack}
		}
	}
	return updateBson
}

func eventCalculateAggs(event Event, shooterID int, ranges []string) Event {
	if event.Shooters[shooterID].Scores == nil {
		event.Shooters[shooterID].Scores = make(map[string]Score)
	}
	var aggID, total, centres int
	var countBack, strRangeID string
	var err error
	for _, strAggID := range ranges {
		aggID, err = strconv.Atoi(strAggID)
		if err == nil {
			total = 0
			centres = 0
			countBack = ""
			for _, rangeID := range event.Ranges[aggID].Aggregate {
				strRangeID = string(rangeID)
				total += event.Shooters[shooterID].Scores[strRangeID].Total
				centres += event.Shooters[shooterID].Scores[strRangeID].Centres
				countBack += event.Shooters[shooterID].Scores[strRangeID].CountBack
			}
			event.Shooters[shooterID].Scores[string(aggID)] = Score{Total: total, Centres: centres, CountBack: countBack}
		}
	}
	return event
}

func totalScoresForm(eventID, rangeID string, shooterID int, score string) Form {
	return Form{
		action: urlUpdateTotalScores,
		inputs: []Inputs{
			{
				name: "score",
				html: "tel",
				//label:   "Total",
				required: true,
				value:    score,
				//TODO add min and max for validation on fclass and target
				//size: 4,
				//min: 0,
				//step: 0.01,
				//max: 50,
				pattern: "[0-9]{1,2}(.[0-9]{1,2}){0,1}",
				//},
				//"centres":Inputs{
				//	html:      "number",
				//	label:   "Centres",
				//	required: true,
				//	value: centres,
				//	size: 4,
				//	min: 0,
				//	max: 10,
				//	//TODO add html5 validation for centres based on total.
				//	//TODO add min = 0, max = parseInt(  total / max(class_valid_shots) )
			}, {
				name: "shooterid",
				html: "hidden",
				//TODO maybe change this to an interface so it can accept different types
				value: fmt.Sprintf("%v", shooterID),
			}, {
				name:  "rangeid",
				html:  "hidden",
				value: rangeID,
			}, {
				html:  "submit",
				inner: "Save",
				name:  "eventid",
				value: eventID,
			},
		},
	}
}
