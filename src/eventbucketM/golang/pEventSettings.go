package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func rangeInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsAddRangeForm("").inputs, r)
	rangeAggInsert(validatedValues, false)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func aggInsert(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsAddAggForm("", []Option{}).inputs, r)
	rangeAggInsert(validatedValues, true)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}
func rangeAggInsert(validatedValues map[string]string, isAgg bool) {
	newRange := Range{Name: validatedValues["name"]}
	if isAgg {
		newRange.IsAgg = true
		newRange.Aggregate = validatedValues["agg"]
	}
	eventID := validatedValues["eventid"]
	rangeID, eventData := eventAddRange(eventID, newRange)
	go calcNewAggRangeScores(eventID, rangeID, eventData)
}
func calcNewAggRangeScores(eventID string, rangeID int, event Event) {
	ranges := []string{fmt.Sprintf("%v", rangeID)}
	for shooterID := range event.Shooters {
		event = eventCalculateAggs(event, shooterID, ranges)
	}
	updateDocByID(tblEvent, eventID, event)
}

func rangeUpdate(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsUpdateRange("", "").inputs, r)
	eventID := validatedValues["eventid"]
	rangeID := validatedValues["rangeid"]
	update := M{
		"$set": M{
			dot("R", rangeID, "n"): validatedValues["name"],
		},
		"$unset": M{},
	}
	mode := "$set"
	if validatedValues["hide"] != "on" {
		mode = "$unset"
	}
	update[mode].(M)[dot("R", rangeID, "h")] = true

	mode = "$set"
	if validatedValues["lock"] != "on" {
		mode = "$unset"
	}
	update[mode].(M)[dot("R", rangeID, "l")] = true

	mode = "$set"
	if validatedValues["aggs"] == "" {
		mode = "$unset"
	}
	update[mode].(M)[dot("R", rangeID, "a")] = validatedValues["aggs"]

	eventUpdateRangeData(eventID, update)
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
}

func eventSettingsUpdateRange(eventID, rangeID string) Form {
	return Form{
		action: urlUpdateRange,
		inputs: []Inputs{
			{
				name:     "name",
				html:     "text",
				label:    "Range Name",
				required: true,
				help:     "A set distance that shooters fire from and scores are recored shot by shot. If you are shooting from more than one range, use \"Create new Agg\" to calculate the totals.",
			}, {
				name:  "rangeid",
				html:  "hidden",
				value: rangeID,
			}, {
				name:    "hide",
				html:    "checkbox",
				checked: false,
			}, {
				name:    "lock",
				html:    "checkbox",
				checked: false,
			}, {
				name:        "aggs",
				html:        "select",
				multiSelect: true,
			}, {
				html:  "submit",
				inner: "Create Range",
				name:  "eventid",
				value: eventID,
			},
		},
	}
}

func eventSettings(eventID string) Page {
	event, _ := getEvent(eventID)
	var eventRanges []Option
	for rangeID, item := range event.Ranges {
		if !item.IsAgg {
			eventRanges = append(eventRanges, Option{Value: fmt.Sprintf("%v", rangeID), Display: item.Name, Selected: true})
		} else {
			var listOfRanges = []Option{}
			for aggID, agg := range event.Ranges {
				if !agg.IsAgg {
					listOfRanges = append(listOfRanges, Option{
						Value:    fmt.Sprintf("%v", aggID),
						Display:  agg.Name,
						Selected: stringInSlice(fmt.Sprintf("%v", aggID), strings.Split(item.Aggregate, ",")),
					})
				}
			}
			var tmp = event.Ranges[rangeID]
			tmp.Aggregate = fmt.Sprintf("<select name=aggs form=range%v multiple size=%v required>%v</select>", rangeID, len(listOfRanges), drawOptions(Inputs{options: listOfRanges}))
			event.Ranges[rangeID] = tmp
		}
	}
	addAggregateForm := ""
	if len(event.Ranges) >= 2 {
		addAggregateForm = generateForm(eventSettingsAddAggForm(eventID, eventRanges))
	}
	return Page{
		TemplateFile: "eventSettings",
		Title:        "Event Settings",
		Theme:        templateAdmin,
		Data: M{
			"Title":          "Event Settings",
			"EventName":      event.Name,
			"ID":             eventID,
			"AddRange":       generateForm(eventSettingsAddRangeForm(eventID)),
			"AddAgg":         addAggregateForm,
			"ListRanges":     event.Ranges,
			"ListGrades":     ClassNamesList,
			"isPrizemeeting": generateForm(eventSettingsIsPrizeMeet(eventID, event.IsPrizeMeet)),
			"menu":           eventMenu(eventID, event.Ranges, urlEventSettings, event.IsPrizeMeet),
			"EventGrades":    generateForm(eventSettingsClassGrades(event)),
			"AllEventGrades": defaultClassSettings,
			"SortScoreboard": generateForm(eventSettingsSortScoreboard(event)),
			"FormNewEvent":   generateForm(homeFormNewEvent(getClubs(), event)),
		},
	}
}
func eventSettingsAddRangeForm(eventID string) Form {
	return Form{
		action: urlEventRangeInsert,
		title:  "Add Range",
		inputs: []Inputs{
			{
				name:      "name",
				html:      "text",
				label:     "Range Name",
				autofocus: true,
				required:  true,
			}, {
				html:  "submit",
				inner: "Create Range",
				name:  "eventid",
				value: eventID,
			},
		},
	}
}

func updateSortScoreBoard(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsSortScoreboard(Event{}).inputs, r)
	eventID := validatedValues["eventid"]
	http.Redirect(w, r, urlScoreboard+eventID, http.StatusSeeOther)
	eventUpdateSortScoreboard(eventID, validatedValues["sort"])
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
		action: urlUpdateSortScoreBoard,
		title:  "Sort Scoreboard",
		inputs: []Inputs{
			{
				name:     "sort",
				html:     "select",
				label:    "Sort Scoreboard by Range",
				required: true,
				options:  sortByRanges,
			}, {
				html:  "submit",
				inner: "Save",
				name:  "eventid",
				value: event.ID,
				help:  "Sort the Scoreboard is normally set to an aggregate range to sort each shooters scores. This is generally a \"Grand Total\" range. Daily aggregate ranges can also be used for different days. Just change it here each day. By default the scoreboard is not sorted at all.",
			},
		},
	}
}

func eventSettingsAddAggForm(eventID string, eventRanges []Option) Form {
	return Form{
		action: urlEventAggInsert,
		title:  "Add Aggregate Range",
		inputs: []Inputs{
			{
				name:     "name",
				html:     "text",
				label:    "Aggregate Name",
				required: true,
				help:     "An Aggregate column (or agg for short) is used to sum up each shooters range scores. These are great as Total columns, Day Aggs and Prizemeeting Grand Aggregates. You can select which ranges and/or other aggs to add together with the \"Agg Ranges\" options. Championships can be setup by using several aggs and disabling previous ranges on the scoreboard that you no longer wish to display.",
			}, {
				name:        "agg",
				html:        "select",
				multiSelect: true,
				options:     eventRanges,
				label:       "Sum up ranges",
				required:    true,
			}, {
				html:  "submit",
				inner: "Create Aggregate",
				name:  "eventid",
				value: eventID,
			},
		},
	}
}

func updateEventGrades(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsClassGrades(Event{}).inputs, r)
	eventID := validatedValues["eventid"]
	grades := strings.Split(validatedValues["grades"], ",")
	var gradeIDs []int
	var gradeIDInt int
	var err error
	for _, gradeID := range grades {
		gradeIDInt, err = strconv.Atoi(gradeID)
		if err == nil {
			gradeIDs = append(gradeIDs, gradeIDInt)
		}
	}
	http.Redirect(w, r, urlEvent+eventID, http.StatusSeeOther)
	tableUpdateData(tblEvent, eventID, M{schemaGRADES: gradeIDs})
}

func eventSettingsClassGrades(event Event) Form {
	return Form{
		action: urlUpdateEventGrades,
		title:  "Available Classes &amp; Grades",
		inputs: []Inputs{
			{
				name:        "grades",
				html:        "select",
				multiSelect: true,
				options:     eventGradeOptions(event.Grades),
			}, {
				html:  "submit",
				inner: "Save",
				name:  "eventid",
				value: event.ID,
				help:  "Select the available classes &amp; grades for this event. This limits the classes that shooters are able to enter.",
				class: "nm",
			},
		},
	}
}

func eventShotsNSighters(eventID string) Page {
	//TODO Display a list of ranges and shooters scores as shots and total scores, ordered in descending order
	event, err := getEvent(eventID)
	if err != nil {
		warning.Printf("Event with id %v doesn't exist", eventID)
	}
	//TODO shooters scores not displayed for each range etc.
	warning.Println("shooters scores not displayed for each range etc.")
	/*for _, eventRange := range event.Ranges{
		//info.Printf("event Range: %v", eventRange.ID)
		for _, shooter := range event.Shooters{
			//TODO check if shooters has this range
			//info.Printf("shooter: %v", shooter.ID)
		}
	}*/
	return Page{
		Theme:        templateEmpty,
		Title:        "eventShotsNSighters",
		TemplateFile: "eventShotsNSighters",
		v8Url:        vURLEventShotsNSighters,
		Data: M{
			"EventName": event.Name,
		},
	}
}

func updateIsPrizeMeet(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(eventSettingsIsPrizeMeet("", false).inputs, r)
	eventID := validatedValues["eventid"]
	prizemeet := false
	if "on" == validatedValues["prizemeet"] {
		prizemeet = true
	}
	http.Redirect(w, r, urlEventSettings+eventID, http.StatusSeeOther)
	tableUpdateData(tblEvent, eventID, M{"p": prizemeet})
}

func eventSettingsIsPrizeMeet(eventID string, checked bool) Form {
	return Form{
		action: urlUpdateIsPrizeMeet,
		title:  "Prize Meeting Event",
		inputs: []Inputs{
			{
				name:    "prizemeet",
				html:    "checkbox",
				label:   "Disable Total Scores page",
				checked: checked,
				class:   "il", //change checkbox to initial style
			}, {
				html:  "submit",
				inner: "Save",
				name:  "eventid",
				value: eventID,
				help:  "Disabling Total Scores page can help with calculating countbacks and discovering incorrect scorecards",
				class: "nm",
			},
		},
	}
}
