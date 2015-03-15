package main

import (
	"fmt"
	"net/http"
	"strings"
)

func totalScores(data string) Page {
	return totalScores_Data(data, false)
}

func totalScoresAll(data string) Page {
	return totalScores_Data(data, true)
}

func totalScores_Data(data string, showAll bool) Page {
	arr := strings.Split(data, "/")
	eventId := arr[0]
	range_Id := arr[1]
	rangeId, err := strToInt(range_Id)
	event, eventMissing := getEvent(eventId)
	currentRange := event.Ranges[rangeId]
	ERROR_ENTER_SCORES_IN_AGG := "<p>This range is an aggregate. Can't enter scores!</p>"
	if err != nil || eventMissing != nil || currentRange.Locked || currentRange.Aggregate != "" {
		return Page{
			TemplateFile: "total-scores",
			Title:        "Total Scores",
			Theme:        TEMPLATE_ADMIN,
			Data: M{
				"Title":      "Total Scores",
				"LinkToPage": "",
				"EventId":    "",
				"RangeName":  "",
				"Message":    ERROR_ENTER_SCORES_IN_AGG,
				"menu":       "",
			},
		}
	}
	var totalScores_link string
	if showAll {
		totalScores_link = fmt.Sprintf("<a href=%v%v>View Incompleted Shooters</a>", URL_totalScores, data)
	} else {
		totalScores_link = fmt.Sprintf("<a href=%v%v>View All Shooters</a>", URL_totalScoresAll, data)
	}

	if len(currentRange.Aggregate) > 0 {
		return Page{
			TemplateFile: "total-scores",
			Title:        "Total Scores",
			Theme:        TEMPLATE_ADMIN,
			Data: M{
				"LinkToPage": totalScores_link,
				"EventId":    eventId,
				"RangeName":  currentRange.Name,
				"Message":    ERROR_ENTER_SCORES_IN_AGG,
				"menu":       eventMenu(eventId, event.Ranges, URL_totalScores, event.IsPrizeMeet),
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

	var shooter_list []EventShooter
	var shooters_forms []string
	for shooterId, shooter_data := range event.Shooters {
		var score string
		if shooter_data.Scores[range_Id].Total > 0 {
			score = fmt.Sprintf("%v", shooter_data.Scores[range_Id].Total)
		}
		if shooter_data.Scores[range_Id].Centers > 0 {
			score += fmt.Sprintf(".%v", shooter_data.Scores[range_Id].Centers)
		}
		if showAll || (!showAll && score == "") {
			shooters_forms = append(shooters_forms, generateForm(totalScoresForm(eventId, range_Id, shooterId, score)))
			shooter_data.Id = shooterId
			shooter_list = append(shooter_list, shooter_data)
		}
	}

	OrderedBy(grade, name).Sort(shooter_list)

	return Page{
		TemplateFile: "total-scores",
		Title:        "Total Scores",
		Theme:        TEMPLATE_ADMIN,
		Data: M{
			"LinkToPage": totalScores_link,
			"EventId":    eventId,
			"RangeName":  currentRange.Name,
			"RangeId":    rangeId,
			"ListRanges": event.Ranges,
			//"ListShooters": event.Shooters,
			"ListShooters":    shooter_list,
			"menu":            eventMenu(eventId, event.Ranges, URL_totalScores, event.IsPrizeMeet),
			"FormTotalScores": shooters_forms,
			"Js":              "total-scores.js",
		},
	}
}

func updateTotalScores(w http.ResponseWriter, r *http.Request) {
	//	validatedValues := check_form(totalScoresForm("", "", "", "").Inputs, r) //totalScoresForm(eventId, rangeId, shooterId, total, centers
	validatedValues, passed := valid8(totalScoresForm("", "", -1, "").Inputs, r)
	if passed {
		eventId := validatedValues["eventid"].(string)
		event := validatedValues["event"].(Event)
		rangeId := validatedValues["rangeid"].(int)
		shooterId := validatedValues["shooterid"].(int)
		score := strings.Split(validatedValues["score"].(string), ".")
		total, err := strToInt(score[0])
		if total > 0 && err == nil {
			new_score := Score{Total: total}
			var centers int
			centers, err = strToInt(score[1])
			if len(score) > 1 && score[1] != "" && centers > 0 && err == nil {
				centers, _ := strToInt(score[1])
				new_score.Centers = centers
			}
			shooterIds := []int{shooterId}
			if event.Shooters[shooterId].LinkedId != nil {
				shooterIds = append(shooterIds, *event.Shooters[shooterId].LinkedId)
			}
			go eventTotalScoreUpdate(eventId, rangeId, shooterIds, new_score)
		}
		http.Redirect(w, r, fmt.Sprintf("%v%v", URL_totalScores+eventId+"/", rangeId), http.StatusSeeOther)
	}
}

func searchForAggs(ranges []Range, rangeId int) []int {
	var aggFound []int
	var num int
	var err error
	for _, rangeObj := range ranges {
		if len(rangeObj.Aggregate) > 0 {
			for _, thisRangeId := range rangeObj.Aggregate {
				num, err = strToInt(thisRangeId)
				if err == nil && num == rangeId {
					aggFound = append(aggFound, *rangeObj.Id)
				}
			}
		}
	}
	return aggFound
}
func eventSearchForAggs(eventId, rangeId string) []string {
	var aggs_to_calculate []string
	event, _ := getEvent(eventId)
	for agg_id, range_data := range event.Ranges {
		if len(range_data.Aggregate) > 0 {
			for _, this_rangeId := range range_data.Aggregate {
				if string(this_rangeId) == rangeId {
					aggs_to_calculate = append(aggs_to_calculate, fmt.Sprintf("%v", agg_id))
				}
			}
		}
	}
	return aggs_to_calculate
}

func calculateAggs(shooterScores map[string]Score, ranges []int, shooterIds []int, eventRanges []Range) M {
	if shooterScores == nil {
		//TODO maybe it's best not to proceed when no scores exist
		shooterScores = make(map[string]Score)
	}
	var total, centers int
	var countBack, strRangeId string
	updateBson := make(M)
	for _, aggId := range ranges {
		total = 0
		centers = 0
		for _, rangeId := range eventRanges[aggId].Aggregate {
			strRangeId = fmt.Sprintf("%v", rangeId)
			total += shooterScores[strRangeId].Total
			centers += shooterScores[strRangeId].Centers
			countBack = shooterScores[strRangeId].CountBack1
		}
		for _, shooterId := range shooterIds {
			updateBson[Dot(schemaSHOOTER, shooterId, aggId)] = Score{Total: total, Centers: centers, CountBack1: countBack}
		}
	}
	return updateBson
}

func eventCalculateAggs(event Event, shooterId int, ranges []string) Event {
	if event.Shooters[shooterId].Scores == nil {
		event.Shooters[shooterId].Scores = make(map[string]Score)
	}
	for _, _agg_id := range ranges {
		agg_id, _ := strToInt(_agg_id)
		total := 0
		centers := 0
		count_back1 := ""
		rangeId := ""
		for _, range_Id := range event.Ranges[agg_id].Aggregate {
			rangeId = string(range_Id)
			total += event.Shooters[shooterId].Scores[rangeId].Total
			centers += event.Shooters[shooterId].Scores[rangeId].Centers
			count_back1 += event.Shooters[shooterId].Scores[rangeId].CountBack1
		}
		event.Shooters[shooterId].Scores[fmt.Sprintf("%v", agg_id)] = Score{Total: total, Centers: centers, CountBack1: count_back1}
	}
	return event
}

func totalScoresForm(eventId, rangeId string, shooterId int, score string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: []Inputs{
			{
				Name: "score",
				Html: "tel",
				//Label:   "Total",
				Required: true,
				Value:    score,
				//TODO add min and max for validation on fclass and target
				//Size: 4,
				//Min: 0,
				//Step: 0.01,
				//Max: 50,
				Pattern: "[0-9]{1,2}(.[0-9]{1,2}){0,1}",
				//},
				//"centers":Inputs{
				//	Html:      "number",
				//	Label:   "Centers",
				//	Required: true,
				//	Value: centers,
				//	Size: 4,
				//	Min: 0,
				//	Max: 10,
				//	//TODO add html5 validation for centers based on total.
				//	//TODO add min = 0, max = parseInt(  total / max(class_valid_shots) )
			}, {
				Name: "shooterid",
				Html: "hidden",
				//TODO maybe change this to an interface so it can accept different types
				Value: fmt.Sprintf("%v", shooterId),
			}, {
				Name:  "rangeid",
				Html:  "hidden",
				Value: rangeId,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			}, {
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}
