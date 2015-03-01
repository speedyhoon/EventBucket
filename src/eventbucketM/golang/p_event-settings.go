package main

import (
	"fmt"
	"net/http"
	"strings"
)

func rangeInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_rangeForm("").Inputs, r)
	range_agg_insert(validated_values, false)
	event_id := validated_values["event_id"]
	http.Redirect(w, r, URL_eventSettings+event_id, http.StatusSeeOther)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_add_aggForm("", []Option{}).Inputs, r)
	range_agg_insert(validated_values, true)
	event_id := validated_values["event_id"]
	http.Redirect(w, r, URL_eventSettings+event_id, http.StatusSeeOther)
}
func range_agg_insert(validated_values map[string]string, isAgg bool) {
	var new_range Range
	new_range.Name = validated_values["name"]
	if isAgg {
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
	go calc_new_agg_range_scores(event_id, range_id, event_data)
}
func calc_new_agg_range_scores(event_id string, range_id int, event Event) {
	ranges := []string{fmt.Sprintf("%v", range_id)}
	for shooter_id := range event.Shooters {
		event = calculate_aggs(event, shooter_id, ranges)
	}
	UpdateDoc_by_id(TBLevent, event_id, event)
}

func rangeUpdate2(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_update_range("", "").Inputs, r)
	eventId := validated_values["event_id"]
	rangeId := validated_values["range_id"]
	update := M{
		"$set": M{
			Dot("R", rangeId, "n"): validated_values["name"],
		},
	}
	if validated_values["hide"] == "on" {
		update["$set"].(M)[Dot("R", rangeId, "h")] = true
	} else {
		update["$unset"].(M)[Dot("R", rangeId, "h")] = false
	}

	if validated_values["lock"] == "on" {
		update["$set"].(M)[Dot("R", rangeId, "l")] = true
	} else {
		update["$unset"].(M)[Dot("R", rangeId, "l")] = false
	}
	event_update_range_data(eventId, update)
	http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
}
func eventSettings_update_range(eventId, rangeId string) Form {
	return Form{
		Action: URL_updateRange,
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			{
				Name:  "event_id",
				Html:  "hidden",
				Value: eventId,
			},
			{
				Name:  "range_id",
				Html:  "hidden",
				Value: rangeId,
			},
			{
				Name:    "hide",
				Html:    "checkbox",
				Checked: false,
			},
			{
				Name:    "lock",
				Html:    "checkbox",
				Checked: false,
			},
			//"aggs": Inputs{
			//	Html:        "select",
			//	MultiSelect: true,
			//},
			{
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func eventSettings(event_id string) Page {
	event, _ := getEvent(event_id)
	var event_ranges []Option
	for range_id, item := range event.Ranges {
		if item.Aggregate == "" {
			event_ranges = append(event_ranges, Option{Value: fmt.Sprintf("%v", range_id), Display: item.Name})
		} else {
			var list_of_ranges = []Option{}
			agg_list := strings.Split(item.Aggregate, ",")
			for agg_id, agg := range event.Ranges {
				if agg.Aggregate == "" {
					ok := stringInSlice(fmt.Sprintf("%v", agg_id), agg_list)
					list_of_ranges = append(list_of_ranges, Option{
						Value:    fmt.Sprintf("%v", agg_id),
						Display:  agg.Name,
						Selected: ok,
					})
				}
			}
			var tmp = event.Ranges[range_id]
			select_options := draw_options(Inputs{Options: list_of_ranges}, "")
			tmp.Aggregate = fmt.Sprintf("<select name=aggs form=range%v multiple size=%v>%v</select>", range_id, len(list_of_ranges), select_options)
			event.Ranges[range_id] = tmp
		}
	}
	var add_agg string
	if len(event.Ranges) >= 2 {
		add_agg = generateForm2(eventSettings_add_aggForm(event_id, event_ranges))
	}
	return Page{
		TemplateFile: "eventSettings",
		Theme:        TEMPLATE_ADMIN,
		Title:        "Event Settings",
		Data: M{
			"Title":          "Event Settings",
			"EventName":      event.Name,
			"Id":             event_id,
			"AddRange":       generateForm2(eventSettings_add_rangeForm(event_id)),
			"AddAgg":         add_agg,
			"ListRanges":     event.Ranges,
			"ListGrades":     CLASSES,
			"isPrizemeeting": generateForm2(eventSettings_isPrizeMeet(event_id, event.IsPrizeMeet)),
			//		"AddDate":        generateForm2(eventSettings_add_dateForm(event_id, event.Date, event.Time)),
			"menu":        event_menu(event_id, event.Ranges, URL_eventSettings, event.IsPrizeMeet),
			"EventGrades": generateForm2(eventSettingsClassGrades(event.Id, event.Grades)),
			//		"ChangeName":     generateForm2(eventSettings_event_name(event.Name, event_id)),
			"AllEventGrades": DEFAULT_CLASS_SETTINGS,
			"SortScoreboard": generateForm2(eventSettings_sort_scoreboard(event_id, event.SortScoreboard, event.Ranges)),
			"FormNewEvent":   generateForm2(home_form_new_event(getClubs(), event.Name, event.Club, event.Date, event.Time, false)),
		},
	}
}
func eventSettings_add_rangeForm(event_id string) Form {
	return Form{
		Action: URL_eventRangeInsert,
		Title:  "Add Range",
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			},
			{
				Name:  "event_id",
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
	validated_values := check_form(eventSettings_sort_scoreboard("", "", make([]Range, 0)).Inputs, r)
	event_id := validated_values["event_id"]
	http.Redirect(w, r, URL_scoreboard+event_id, http.StatusSeeOther)
	event_update_sort_scoreboard(event_id, validated_values["sort"])
}

func eventSettings_sort_scoreboard(event_id string, existing_sort string, ranges []Range) Form {
	var sort_by_ranges []Option
	var sort_by bool
	for index, Range := range ranges {
		sort_by = false
		if fmt.Sprintf("%v", index) == existing_sort {
			sort_by = true
		}
		sort_by_ranges = append(sort_by_ranges, Option{Display: Range.Name, Value: fmt.Sprintf("%v", index), Selected: sort_by})
	}
	return Form{
		Action: URL_updateSortScoreBoard,
		Title:  "Sort Scoreboard",
		Inputs: []Inputs{
			{
				Name:     "sort",
				Html:     "select",
				Label:    "Sort Scoreboard by Range",
				Required: true,
				Options:  sort_by_ranges,
			},
			{
				Name:  "event_id",
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

func eventSettings_add_aggForm(event_id string, event_ranges []Option) Form {
	return Form{
		Action: URL_eventAggInsert,
		Title:  "Add Aggregate Range",
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Aggregate Name",
				Required: true,
			},
			{
				Name:  "event_id",
				Html:  "hidden",
				Value: event_id,
			},
			{
				Name:        "agg",
				Html:        "select",
				MultiSelect: true,
				Options:     event_ranges,
				Label:       "Sum up ranges",
				Required:    true,
			},
			{
				Html:  "submit",
				Value: "Create Aggregate",
			},
		},
	}
}

func updateEventGrades(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettingsClassGrades("", []int{}).Inputs, r)
	event_id := validated_values["event_id"]
	http.Redirect(w, r, URL_event+event_id, http.StatusSeeOther)
	event_upsert_data(event_id, M{schemaGRADES: validated_values["grades"]})
}

func slice_to_map_bool(input []string) map[string]bool {
	output := make(map[string]bool)
	for _, value := range input {
		output[value] = true
	}
	return output
}

func eventSettingsClassGrades(eventId string, grades []int) Form {
	return Form{
		Action: URL_updateEventGrades,
		Title:  "Classes &amp; Grades",
		Inputs: []Inputs{
			{
				Name: "grades",
				Html: "select",
				//				Label:          "select Classes &amp; Grades in this event",
				MultiSelect: true,
				Options:     eventGradeOptions(grades),
			},
			{
				Name:  "event_id",
				Html:  "hidden",
				Value: eventId,
			},
			{
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func eventShotsNSighters(eventId string) Page {
	//TODO Display a list of ranges and shooters scores as shots and total scores, ordered in descending order
	//Info.Printf("eventid = %v",eventId)
	event, err := getEvent(eventId)
	if err == nil {
		//TODO shooters scores not displayed for each range etc.
		Warning.Println("shooters scores not displayed for each range etc.")
		/*for _, eventRange := range event.Ranges{
			//Info.Printf("event Range: %v", eventRange.Id)
			for _, shooter := range event.Shooters{
				//TODO check if shooters has this range
				//Info.Printf("shooter: %v", shooter.Id)
			}
		}*/
	}
	return Page{
		Theme:        TEMPLATE_EMPTY,
		Title:        "eventShotsNSighters",
		TemplateFile: "eventShotsNSighters",
		Data: M{
			"EventName": event.Name,
		},
		v8Url: VURL_eventShotsNSighters,
	}
}

func updateIsPrizeMeet(w http.ResponseWriter, r *http.Request) {
	validated_values := check_form(eventSettings_isPrizeMeet("", false).Inputs, r)
	event_id := validated_values["event_id"]
	prizemeet := false
	if "on" == validated_values["prizemeet"] {
		prizemeet = true
	}
	http.Redirect(w, r, URL_eventSettings+event_id, http.StatusSeeOther)
	event_upsert_data(event_id, M{"p": prizemeet})
}

func eventSettings_isPrizeMeet(event_id string, checked bool) Form {
	return Form{
		Action: URL_updateIsPrizeMeet,
		Title:  "Prize Meeting Event",
		Inputs: []Inputs{
			{
				Name:    "prizemeet",
				Html:    "checkbox",
				Label:   "Is this Event a Prize Meeting?",
				Checked: checked,
			},
			{
				Name:  "event_id",
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
