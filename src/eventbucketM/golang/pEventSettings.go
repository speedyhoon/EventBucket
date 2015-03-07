package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func rangeInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsAddRangeForm("").Inputs, r)
	rangeAggInsert(validatedValues, false)
	eventId := validatedValues["eventid"]
	http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsAddAggForm("", []Option{}).Inputs, r)
	rangeAggInsert(validatedValues, true)
	eventId := validatedValues["eventid"]
	http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
}
func rangeAggInsert(validatedValues map[string]string, isAgg bool) {
	new_range := Range{Name: validatedValues["name"]}
	if isAgg {
		new_range.IsAgg = true
		new_range.Aggregate = validatedValues["agg"]
	}
	eventId := validatedValues["eventid"]
	range_id, event_data := EventAddRange(eventId, new_range)
	go calc_new_agg_range_scores(eventId, range_id, event_data)
}
func calc_new_agg_range_scores(eventId string, range_id int, event Event) {
	ranges := []string{fmt.Sprintf("%v", range_id)}
	for shooter_id := range event.Shooters {
		event = eventCalculateAggs(event, shooter_id, ranges)
	}
	UpdateDocById(TBLevent, eventId, event)
}

func rangeUpdate(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsUpdateRange("", "").Inputs, r)
	eventId := validatedValues["eventid"]
	rangeId := validatedValues["rangeid"]
	update := M{
		"$set": M{
			Dot("R", rangeId, "n"): validatedValues["name"],
		},
		"$unset": M{},
	}
	thingy := "$set"
	if validatedValues["hide"] != "on" {
		thingy = "$unset"
	}
	update[thingy].(M)[Dot("R", rangeId, "h")] = true

	thingy = "$set"
	if validatedValues["lock"] != "on" {
		thingy = "$unset"
	}
	update[thingy].(M)[Dot("R", rangeId, "l")] = true

	thingy = "$set"
	if validatedValues["aggs"] == "" {
		thingy = "$unset"
	}
	update[thingy].(M)[Dot("R", rangeId, "a")] = validatedValues["aggs"]
	eventUpdateRangeData(eventId, update)
	http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
}

func eventSettingsUpdateRange(eventId, rangeId string) Form {
	return Form{
		Action: URL_updateRange,
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			}, {
				Name:  "rangeid",
				Html:  "hidden",
				Value: rangeId,
			}, {
				Name:    "hide",
				Html:    "checkbox",
				Checked: false,
			}, {
				Name:    "lock",
				Html:    "checkbox",
				Checked: false,
			}, {
				Name:        "aggs",
				Html:        "select",
				MultiSelect: true,
			}, {
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func eventSettings(eventId string) Page {
	event, _ := getEvent(eventId)
	var event_ranges []Option
	for range_id, item := range event.Ranges {
		if !item.IsAgg {
			event_ranges = append(event_ranges, Option{Value: fmt.Sprintf("%v", range_id), Display: item.Name, Selected: true})
		} else {
			var list_of_ranges = []Option{}
			agg_list := strings.Split(item.Aggregate, ",")
			for agg_id, agg := range event.Ranges {
				if !agg.IsAgg {
					list_of_ranges = append(list_of_ranges, Option{
						Value:    fmt.Sprintf("%v", agg_id),
						Display:  agg.Name,
						Selected: stringInSlice(fmt.Sprintf("%v", agg_id), agg_list),
					})
				}
			}
			var tmp = event.Ranges[range_id]
			select_options := drawOptions(Inputs{Options: list_of_ranges}, "")
			tmp.Aggregate = fmt.Sprintf("<select name=aggs form=range%v multiple size=%v required>%v</select>", range_id, len(list_of_ranges), select_options)
			event.Ranges[range_id] = tmp
		}
	}
	addAggregateForm := ""
	if len(event.Ranges) >= 2 {
		addAggregateForm = generateForm(eventSettingsAddAggForm(eventId, event_ranges))
	}
	return Page{
		TemplateFile: "eventSettings",
		Title:        "Event Settings",
		Theme:        TEMPLATE_ADMIN,
		Data: M{
			"Title":          "Event Settings",
			"EventName":      event.Name,
			"Id":             eventId,
			"AddRange":       generateForm(eventSettingsAddRangeForm(eventId)),
			"AddAgg":         addAggregateForm,
			"ListRanges":     event.Ranges,
			"ListGrades":     CLASSES,
			"isPrizemeeting": generateForm(eventSettingsIsPrizeMeet(eventId, event.IsPrizeMeet)),
			"menu":           eventMenu(eventId, event.Ranges, URL_eventSettings, event.IsPrizeMeet),
			"EventGrades":    generateForm(eventSettingsClassGrades(event)),
			"AllEventGrades": DEFAULT_CLASS_SETTINGS,
			"SortScoreboard": generateForm(eventSettingsSortScoreboard(event)),
			"FormNewEvent":   generateForm(homeFormNewEvent(getClubs(), event)),
		},
	}
}
func eventSettingsAddRangeForm(eventId string) Form {
	return Form{
		Action: URL_eventRangeInsert,
		Title:  "Add Range",
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Range Name",
				Required: true,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			}, {
				Html:  "submit",
				Value: "Create Range",
			},
		},
	}
}

func updateSortScoreBoard(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsSortScoreboard(Event{}).Inputs, r)
	eventId := validatedValues["eventid"]
	http.Redirect(w, r, URL_scoreboard+eventId, http.StatusSeeOther)
	eventUpdateSortScoreboard(eventId, validatedValues["sort"])
}

func eventSettingsSortScoreboard(event Event) Form {
	var sortByRanges []Option
	var sortBy bool
	for index, Range := range event.Ranges {
		if fmt.Sprintf("%v", index) == event.SortScoreboard {
			sortBy = true
		}
		sortByRanges = append(sortByRanges, Option{Display: Range.Name, Value: fmt.Sprintf("%v", index), Selected: sortBy})
		sortBy = false
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
				Options:  sortByRanges,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: event.Id,
			}, {
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func eventSettingsAddAggForm(eventId string, eventRanges []Option) Form {
	return Form{
		Action: URL_eventAggInsert,
		Title:  "Add Aggregate Range",
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Aggregate Name",
				Required: true,
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: eventId,
			}, {
				Name:        "agg",
				Html:        "select",
				MultiSelect: true,
				Options:     eventRanges,
				Label:       "Sum up ranges",
				Required:    true,
			}, {
				Html:  "submit",
				Value: "Create Aggregate",
			},
		},
	}
}

func updateEventGrades(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsClassGrades(Event{}).Inputs, r)
	eventId := validatedValues["eventid"]
	grades := strings.Split(validatedValues["grades"], ",")
	var gradeIds []int
	var gradeIdInt int
	var err error
	for _, gradeId := range grades {
		gradeIdInt, err = strconv.Atoi(gradeId)
		if err == nil {
			gradeIds = append(gradeIds, gradeIdInt)
		}
	}
	http.Redirect(w, r, URL_event+eventId, http.StatusSeeOther)
	tableUpdateData(TBLevent, eventId, M{schemaGRADES: gradeIds})
}

func eventSettingsClassGrades(event Event) Form {
	return Form{
		Action: URL_updateEventGrades,
		Title:  "Availalbe Classes &amp; Grades",
		Inputs: []Inputs{
			{
				Name:        "grades",
				Html:        "select",
				MultiSelect: true,
				Options:     eventGradeOptions(event.Grades),
			}, {
				Name:  "eventid",
				Html:  "hidden",
				Value: event.Id,
			}, {
				Html:  "submit",
				Value: "Save",
			},
		},
	}
}

func eventShotsNSighters(eventId string) Page {
	//TODO Display a list of ranges and shooters scores as shots and total scores, ordered in descending order
	event, err := getEvent(eventId)
	if err != nil {
		Error.Printf("Event with id %v doesn't exist", eventId)
	}
	//TODO shooters scores not displayed for each range etc.
	Warning.Println("shooters scores not displayed for each range etc.")
	/*for _, eventRange := range event.Ranges{
		//Info.Printf("event Range: %v", eventRange.Id)
		for _, shooter := range event.Shooters{
			//TODO check if shooters has this range
			//Info.Printf("shooter: %v", shooter.Id)
		}
	}*/
	return Page{
		Theme:        TEMPLATE_EMPTY,
		Title:        "eventShotsNSighters",
		TemplateFile: "eventShotsNSighters",
		v8Url:        VURL_eventShotsNSighters,
		Data: M{
			"EventName": event.Name,
		},
	}
}

func updateIsPrizeMeet(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsIsPrizeMeet("", false).Inputs, r)
	eventId := validatedValues["eventid"]
	prizemeet := false
	if "on" == validatedValues["prizemeet"] {
		prizemeet = true
	}
	http.Redirect(w, r, URL_eventSettings+eventId, http.StatusSeeOther)
	eventUpsertData(eventId, M{"p": prizemeet})
}

func eventSettingsIsPrizeMeet(eventId string, checked bool) Form {
	return Form{
		Action: URL_updateIsPrizeMeet,
		Title:  "Prize Meeting Event",
		Inputs: []Inputs{
			{
				Name:    "prizemeet",
				Html:    "checkbox",
				Label:   "Is this Event a Prize Meeting?",
				Checked: checked,
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
