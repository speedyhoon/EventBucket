package main

import (
	"net/http"
	"fmt"
	"strings"
)
func totalScores(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_totalScores)
	templator(TEMPLATE_ADMIN, "total-scores", totalScores_Data(data, false), w)
}
func totalScoresAll(w http.ResponseWriter, r *http.Request) {
	data := get_id_from_url(r, URL_totalScoresAll)
	templator(TEMPLATE_ADMIN, "total-scores", totalScores_Data(data, true), w)
}
func totalScores_Data(data string, show_all bool) M {
	arr := strings.Split(data, "/")
	event_id := arr[0]
	rangeId, success := strToInt(arr[1])
	range_id := arr[1]
	event, eventMissing := getEvent(event_id)
	ERROR_ENTER_SCORES_IN_AGG := "<p>This range is an aggregate. Can't enter scores!</p>"
	if !success || eventMissing || event.Ranges[rangeId].Locked || event.Ranges[rangeId].Aggregate != ""{
		return M{
			"Title": "Total Scores",
			"LinkToPage": "",
			"EventId": "",
			"RangeName": "",
			"Message": ERROR_ENTER_SCORES_IN_AGG,
			"menu": "",
		}
	}

	selected_range :=  event.Ranges[rangeId]

	var totalScores_link string
	if show_all{
		totalScores_link = fmt.Sprintf("<a href=%v%v>View Incompleted Shooters</a>", URL_totalScores, data)
	}else{
		totalScores_link = fmt.Sprintf("<a href=%v%v>View All Shooters</a>", URL_totalScoresAll, data)
	}

	if len(selected_range.Aggregate) > 0{
		return M{
			"Title": "Total Scores",
			"LinkToPage": totalScores_link,
			"EventId": event_id,
			"RangeName": selected_range.Name,
			"Message": ERROR_ENTER_SCORES_IN_AGG,
			"menu": event_menu(event_id, event.Ranges, URL_totalScores, event.IsPrizeMeet),
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
	shooters_forms := make(map[string]string)
	for shooter_id,shooter_data := range event.Shooters{
		var score string
		if shooter_data.Scores[range_id].Total > 0 {
			score = fmt.Sprintf("%v", shooter_data.Scores[range_id].Total)
		}
		if shooter_data.Scores[range_id].Centers > 0 {
			score += fmt.Sprintf(".%v", shooter_data.Scores[range_id].Centers)
		}
		if show_all || (!show_all && score == "") {
			shooters_forms[shooter_id] = generateForm2(total_scores_Form(event_id, range_id, shooter_id, score))
			shooter_data.Id = shooter_id
			shooter_list = append(shooter_list, shooter_data)
		}
	}

	OrderedBy(grade, name).Sort(shooter_list)

	return M{
		"Title": "Total Scores",
		"LinkToPage": totalScores_link,
		"EventId": event_id,
		"RangeName": selected_range.Name,
		"RangeId": range_id,
		"ListRanges": event.Ranges,

//		"ListShooters": event.Shooters,
		"ListShooters": shooter_list,
		"menu": event_menu(event_id, event.Ranges,URL_totalScores, event.IsPrizeMeet),
		"FormTotalScores": shooters_forms,
		"Js": "total-scores.js",
	}
}

func updateTotalScores(w http.ResponseWriter, r *http.Request){
//	validated_values := check_form(total_scores_Form("", "", "", "").Inputs, r) //total_scores_Form(event_id, range_id, shooter_id, total, centers
	validated_values, passed := valid8(total_scores_Form("", "", "", "").Inputs, r)
	if passed {
		event_id := validated_values["event_id"].(string)
		event := validated_values["event"].(Event)
		rangeId := validated_values["range_id"].(int)
		shooter_id := validated_values["shooter_id"].(string)
		score := strings.Split(validated_values["score"].(string), ".")
		total, success := strToInt(score[0])
		if total > 0 && success {
			new_score := Score{Total: total}
			var centers int
			centers, success = strToInt(score[1])
			if len(score) > 1 && score[1] != "" && centers > 0 && success {
				centers, _ := strToInt(score[1])
				new_score.Centers = centers
			}
//			rangeId, _ := strToInt(range_id)
			shooterIds := []string{shooter_id}
			if event.Shooters[shooter_id].LinkedId != "" {
				shooterIds = append(shooterIds, event.Shooters[shooter_id].LinkedId)
			}
			go eventTotalScoreUpdate(event_id, rangeId, shooterIds, new_score)
		}
		redirecter(fmt.Sprintf("%v%v", URL_totalScores+event_id+"/", rangeId), w, r)
	}
}



func search_for_aggs(event_id, range_id string)[]string{
	var aggs_to_calculate []string
	event, _ := getEvent(event_id)
	for agg_id, range_data := range event.Ranges{
		if len(range_data.Aggregate) > 0{
			for _, this_range_id := range range_data.Aggregate{
				if string(this_range_id) == range_id{
					aggs_to_calculate = append(aggs_to_calculate, fmt.Sprintf("%v",agg_id))
				}
			}
		}
	}
	return aggs_to_calculate
}
func calculate_aggs(event Event, shooter_id string, ranges []string)Event{
//	if xx, ok := event.Shooters[shooter_id]; ok {
//		xx.count = 2
//		m["x"] = xx
//	} else {
//		panic("X isn't in the map")
//	}
//	if event.Shooters[shooter_id].Scores != nil{
//		dump("new val is not none")
//	}else{
	if event.Shooters[shooter_id].Scores == nil{
		temp_kkk := event.Shooters[shooter_id]
		temp_kkk.Scores = map[string]Score{}
		event.Shooters[shooter_id] = temp_kkk
	}
	for _, _agg_id := range ranges {
		agg_id, _ := strToInt(_agg_id)
		total := 0
		centers := 0
		count_back1 := ""
		range_id := ""
		for _, rangeId := range event.Ranges[agg_id].Aggregate {
			range_id = string(rangeId)
			total += event.Shooters[shooter_id].Scores[range_id].Total
			centers += event.Shooters[shooter_id].Scores[range_id].Centers
			count_back1 += event.Shooters[shooter_id].Scores[range_id].CountBack1
		}
		event.Shooters[shooter_id].Scores[fmt.Sprintf("%v",agg_id)] = Score{Total: total, Centers: centers, CountBack1: count_back1}
	}
//		event.Shooters[shooter_id].Scores = make([]Score{}, 1)
//		agg_total := event.Shooters[shooter_id].Scores[agg_id]
//		agg_total.Total = 0
//		agg_total.Centers = 0
//		agg_total.CountBack1 = ""
//		for _, rangeId := range event.Ranges[agg_id].Aggregate{
//			range_id := string(rangeId)
//			agg_total.Total += event.Shooters[shooter_id].Scores[range_id].Total
//			agg_total.Centers += event.Shooters[shooter_id].Scores[range_id].Centers
//			agg_total.CountBack1 += event.Shooters[shooter_id].Scores[range_id].CountBack1
//		}
//		event.Shooters[shooter_id].Scores[agg_id] = agg_total
//	}
	return event
}

func total_scores_Form(event_id, range_id, shooter_id, score string) Form {
	return Form{
		Action: URL_updateTotalScores,
		Inputs: []Inputs{
			{
				Name: "score",
				Html:      "tel",
//				Label:   "Total",
				Required: true,
				Value: score,
				//TODO add min and max for validation on fclass and taget
//				Size: 4,
//				Min: 0,
//				Step: 0.01,
//				Max: 50,
				Pattern: "[0-9]{1,2}(.[0-9]{1,2}){0,1}",
			},
//			"centers":Inputs{
//				Html:      "number",
//				Label:   "Centers",
//				Required: true,
//				Value: centers,
//				Size: 4,
//				Min: 0,
//				Max: 10,
//				//TODO add html5 validation for centers based on total.
//				//TODO add min = 0, max = parseInt(  total / max(class_valid_shots) )
//			},
			{
				Name: "shooter_id",
				Html: "hidden",
				Value: shooter_id,
			},
			{
				Name: "range_id",
				Html: "hidden",
				Value: range_id,
			},
			{
				Name: "event_id",
				Html: "hidden",
				Value: event_id,
			},
			{
				Html:    "submit",
				Value:   "Save",
			},
		},
	}
}
