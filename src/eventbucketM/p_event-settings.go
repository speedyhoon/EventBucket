package main

import (
	"fmt"
	"net/http"
	"strings"
//	"time"
)

//func updateEventName(w http.ResponseWriter, r *http.Request) {
//	validated_values := check_form(eventSettings_event_name("","").Inputs, r)
//	event_id := validated_values["event_id"]
//	event_update_name(event_id, validated_values["name"])
//	redirecter(URL_eventSettings+event_id, w, r)
//}

func rangeInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_rangeForm("").Inputs, r)
	range_agg_insert(validated_values, false)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_aggForm("", []Option{}).Inputs, r)
	range_agg_insert(validated_values, true)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
}
func range_agg_insert(validated_values map[string]string, isAgg bool) {
	var new_range Range
	new_range.Name = validated_values["name"]
	if isAgg{
//	if validated_values["agg"] != "" {
		//		new_range.Aggregate = []int64{}
		//		for _, range_name := range strings.Split(validated_values["agg"], ","){
		//			new_range.Aggregate = append(new_range.Aggregate, str_to_int64(range_name))
		//		}
		//		new_range.Aggregate = strings.Split(validated_values["agg"], ",")
		new_range.Aggregate = validated_values["agg"]
	}
	event_id := validated_values["event_id"]
	range_id, event_data := DB_event_add_range(event_id, new_range)
//	if range_id > 0 {
		go calc_new_agg_range_scores(event_id, range_id, event_data)
//	}
}
func calc_new_agg_range_scores(event_id string, range_id int, event Event) {
	ranges := []string{fmt.Sprintf("%v",range_id)}
	for shooter_id := range event.Shooters {
		event = calculate_aggs(event, shooter_id, ranges)
	}
	UpdateDoc_by_id(TBLevent, event_id, event)
}

func rangeUpdate2(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_update_range("", "").Inputs, r)
	event_id := validated_values["event_id"]
	range_id := validated_values["range_id"]
	dump("range update 2")
	dump(event_id)
	dump(range_id)
	dump(validated_values)

	unset := make(M)
	set := M{
		Dot("R",range_id,"n"): validated_values["name"],
	}

//	did_unset := false
	if validated_values["hide"] == "on"{
		set[Dot("R",range_id,"h")] = true
	}else{
		unset[Dot("R",range_id,"h")] = ""
//		did_unset = true
	}
	if validated_values["lock"] == "on"{
		set[Dot("R",range_id,"l")] = true
	}else{
		unset[Dot("R",range_id,"l")] = ""
//		did_unset = true
	}

	update := M{
		"$set": set,
	}
//	if did_unset{
	if len(unset) >= 1{
		update["$unset"] = unset
	}

	dump("update Data:::")
	dump(update)
	event_update_range_data(event_id, update)
	redirecter(URL_eventSettings+event_id, w, r)
}
func eventSettings_update_range(event_id, range_id string) Form {
	return Form{
		Action: URL_updateRange,
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Name: "range_id",
				Html:  "hidden",
				Value: range_id,
			},
			{
				Name: "hide",
				Html:    "checkbox",
				Checked: false,
			},
			{
				Name: "lock",
				Html:    "checkbox",
				Checked: false,
			},
//			"aggs": Inputs{
//				Html:        "select",
//				MultiSelect: true,
//			},
			{
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func eventSettings(event_id string) M {
	event, _ := getEvent(event_id)
	var event_ranges []Option
	for range_id, item := range event.Ranges {
		if item.Aggregate == "" {
			event_ranges = append(event_ranges, Option{Value:fmt.Sprintf("%v",range_id), Display:item.Name})
		} else {
			var list_of_ranges = []Option{}
			agg_list := strings.Split(item.Aggregate, ",")
			for agg_id, agg := range event.Ranges {
				if agg.Aggregate == "" {
					ok := stringInSlice(fmt.Sprintf("%v",agg_id), agg_list)
					list_of_ranges = append(list_of_ranges, Option{
						Value:    fmt.Sprintf("%v",agg_id),
						Display:  agg.Name,
						Selected: ok,
					})
				}
			}
			var tmp = event.Ranges[range_id]
			select_options := draw_options(Inputs{Options:list_of_ranges},"")
			tmp.Aggregate = fmt.Sprintf("<select name=aggs form=range%v multiple size=%v>%v</select>", range_id, len(list_of_ranges), select_options)
			event.Ranges[range_id] = tmp
		}
	}
	var add_agg string
	if len(event.Ranges) >= 2 {
		add_agg = generateForm2(eventSettings_add_aggForm(event_id, event_ranges))
	}
	return M{
		"Title":          "Event Settings",
		"EventName":      event.Name,
		"Id":             event_id,
		"AddRange":       generateForm2(eventSettings_add_rangeForm(event_id)),
		"AddAgg":         add_agg,
		"ListRanges":     event.Ranges,
		"ListGrades":     CLASSES,
		"isPrizemeeting":	generateForm2(eventSettings_isPrizeMeet(event_id, event.IsPrizeMeet)),
//		"AddDate":        generateForm2(eventSettings_add_dateForm(event_id, event.Date, event.Time)),
		"menu":           event_menu(event_id, event.Ranges, URL_eventSettings, event.IsPrizeMeet),
		"EventGrades":    generateForm2(eventSettings_class_grades(event)),
//		"ChangeName":     generateForm2(eventSettings_event_name(event.Name, event_id)),
		"AllEventGrades": DEFAULT_CLASS_SETTINGS,
		"SortScoreboard": generateForm2(eventSettings_sort_scoreboard(event_id, event.SortScoreboard, event.Ranges)),
		"FormNewEvent": 	generateForm2(home_form_new_event(getClubs(), event.Name,event.Club,event.Date,event.Time,false)),
	}
}
func eventSettings_add_rangeForm(event_id string) Form {
	return Form{
		Action: URL_eventRangeInsert,
		Title:  "Add Range",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func updateSortScoreBoard(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_sort_scoreboard("", "", make([]Range,0)).Inputs, r)
	event_id := validated_values["event_id"]
//	redirecter(URL_eventSettings+event_id, w, r)
	redirecter(URL_scoreboard+event_id, w, r)
	event_update_sort_scoreboard(event_id, validated_values["sort"])
}

func eventSettings_sort_scoreboard(event_id string, existing_sort string, ranges []Range) Form {
	var sort_by_ranges []Option
	var sort_by bool
	for index, Range := range ranges {
		sort_by = false
		if fmt.Sprintf("%v",index )== existing_sort {
			sort_by = true
		}
		sort_by_ranges = append(sort_by_ranges, Option{Display: Range.Name, Value: fmt.Sprintf("%v",index), Selected: sort_by})
	}
	//	export(sort_by_ranges)
	return Form{
		Action: URL_updateSortScoreBoard,
		Title:  "Sort Scoreboard",
		Inputs: []Inputs{
			{
				Name: "sort",
				Html:           "select",
				Label:          "Sort Scoreboard by Range",
				Required:       true,
				Options: sort_by_ranges,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

/*func dateUpdate(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_dateForm("", "", "").Inputs, r)
	event_id := validated_values["event_id"]
	redirecter(URL_eventSettings+event_id, w, r)
	event_update_date(event_id, validated_values["date"], validated_values["time"])
}*/
/*func eventSettings_add_dateForm(event_id, date, hour_minute string) Form {
	if date == "" {
		date = time.Now().Format("2006-02-01")
	}
	if hour_minute == "" {
		hour_minute = time.Now().Format("15:04")
	}
	return Form{
		Action: URL_dateUpdate,
		Title:  "Date &amp; Time",
		Inputs: []Inputs{
			{
				Name: "date",
				Html:     "date",
				Label:    "Date",
				Required: true,
				Value:    date,
			},
			{
				Name: "time",
				Html:  "time",
				Label: "Time",
				Value: hour_minute,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Save Date",
			},
		},
	}
}*/
func eventSettings_add_aggForm(event_id string, event_ranges []Option) Form {
	return Form{
		Action: URL_eventAggInsert,
		Title:  "Add Aggregate Range",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:     "text",
				Label:    "Aggregate Name",
				Required: true,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Name: "agg",
				Html:         "select",
				MultiSelect:  true,
				Options: event_ranges,
				Label:        "Sum up ranges",
				Required: true,
			},
			{
				Html:  "submit",
				Value: "Create Aggregate",
			},
		},
	}
}

func updateEventGrades(w http.ResponseWriter, r *http.Request) {
	var event Event
	validated_values := check_form(eventSettings_class_grades(event).Inputs, r)
	event_id := validated_values["event_id"]
//	redirecter(URL_eventSettings+event_id, w, r)
	redirecter(URL_event+event_id, w, r)
	event_upsert_data(event_id, M{schemaGRADES: validated_values["grades"]})
}

func slice_to_map_bool(input []string) map[string]bool {
	output := make(map[string]bool)
	for _, value := range input {
		output[value] = true
	}
	return output
}

func eventSettings_class_grades(event Event) Form {
	allGrades := grades()
	var grades []Option
	selected := false
	var grade_list map[string]bool
	selected_grades := strings.Split(event.Grades, ",")
	no_grades_selected := event.Grades == ""
	if !no_grades_selected {
		grade_list = slice_to_map_bool(selected_grades)
	}

	for _, class_settings := range DEFAULT_CLASS_SETTINGS{
		for _, grade_id := range class_settings.Grades {
			selected = false
			gradeId := fmt.Sprintf("%v",grade_id)
			if grade_list[gradeId] || no_grades_selected {
				selected = true
			}
			grades = append(grades, Option{
				Value:    gradeId,
				Display:  allGrades[grade_id].LongName,
				Selected: selected,
			})
		}
	}
	var event_id string
	if event.Id != "" {
		event_id = event.Id
	}
	return Form{
		Action: URL_updateEventGrades,
		Title:  "Classes &amp; Grades",
		Inputs: []Inputs{
			{
				Name: "grades",
				Html:           "select",
//				Label:          "select Classes &amp; Grades in this event",
				MultiSelect:    true,
				Options: grades,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

/*func eventSettings_event_name(event_name, event_id string) Form {
	return Form{
		Action: URL_updateEventName,
		Title:  "Event name",
		Inputs: []Inputs{
			{
				Name: "name",
				Html:        "text",
				Label:       "Change event name",
				Value:       event_name,
				Placeholder: event_name,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}*/

//func totalScores_update(event_id, shooter_id, range_id string)  Form {
//	return Form{
//		Action: URL_updateTotalScores,
//		Inputs: map[string]Inputs{
//			schemaTOTAL:Inputs{
//				Html:      "number",
//				Label:   "Total",
//				Required: true,
//				Min: 0,
//				Max: 60,
//			},
//			schemaCENTER:Inputs{
//				Html:      "number",
//				Label:   "Centers",
//				Required: true,
//				Min: 0,
//				Max: 60,
//			},
//			"event_id":Inputs{
//				Html: "hidden",
//				Value: event_id,
//				Required: true,
//			},
//			"shooter_id":Inputs{
//				Html: "hidden",
//				Value: shooter_id,
//				Required: true,
//			},
//			"range_id":Inputs{
//				Html: "hidden",
//				Value: range_id,
//				Required: true,
//			},
//			"submit":Inputs{
//				Html:      "submit",
//				Value:   "Save",
//			},
//		},
//	}
//}




func eventShotsNSighters(w http.ResponseWriter, r *http.Request) {
//	var event Event
	r.ParseForm()
	form := r.Form
//	fmt.Println("event_id:::", r.Form["event_id"])
//	fmt.Println("shots:::", r.Form["shots"])
//	fmt.Println("sight:::", r.Form["sight"])
//	fmt.Println("form:")
//	export(form)

	if event_id, ok := form["event_id"]; ok && len(event_id) > 0{
		fmt.Println("event_id=",event_id)
		if shots, ok := form["shots"]; ok {
			fmt.Println("shots...")
			for range_id, range_data := range shots {
				for class_id, shot_value := range range_data {
					fmt.Println("range=",range_id," class=",class_id," value=",shot_value)
				}
			}
		}else{
			fmt.Println("shots not found")
		}
	}

//	validated_values := check_form(eventSettings_class_grades(event).Inputs, r)

//	dump(form)

//	event_id := validated_values["event_id"]
//	redirecter(URL_eventSettings+event_id, w, r)
//	event_upsert_data(event_id, M{schemaGRADES: validated_values["grades"]})
}





func updateIsPrizeMeet(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_isPrizeMeet("", false).Inputs, r)
	event_id := validated_values["event_id"]
	prizemeet := false
	if "on" ==validated_values["prizemeet"]{
		prizemeet = true
	}
	redirecter(URL_eventSettings+event_id, w, r)
	event_upsert_data(event_id, M{"p": prizemeet})
}

func eventSettings_isPrizeMeet(event_id string, checked bool) Form {
	return Form{
		Action: URL_updateIsPrizeMeet,
		Title:  "Prize Meeting Event",
		Inputs: []Inputs{
			{
				Name: "prizemeet",
				Html:     "checkbox",
				Label:    "Is this Event a Prize Meeting?",
				Checked:   checked,
			},
			{
				Name: "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}
