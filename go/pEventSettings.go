package main

import (
	"fmt"
	"net/http"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/forms"
)

func eventSettings(w http.ResponseWriter, r *http.Request, event Event) {
	//Retrieve any submitted form that failed to save.
	action, forms := session.Forms(w, r, getForm, eventEdit, eventRangeNew, eventAggNew, eventRangeEdit, eventAggEdit, eventAvailableGrades)
	if action != eventEdit {
		forms[0].Fields[0].Value = event.Club.Name
		forms[0].Fields[1].Value = event.Name
		forms[0].Fields[2].Value = event.Date
		forms[0].Fields[3].Value = event.Time
		forms[0].Fields[4].Checked = event.Closed
		forms[0].Fields[5].Value = event.ID
	}
	forms[1].Fields[1].Value = event.ID

	forms[2].Fields[1].Options = dataListRanges(event.Ranges, true)
	forms[2].Fields[2].Value = event.ID

	//AvailableGrades
	forms[5].Fields[0].Options = availableGrades(event.Grades)
	forms[5].Fields[1].Value = event.ID

	render(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":               dataListRanges(event.Ranges, false),
			"Event":                event,
			"eventEdit":            forms[0],
			"eventRangeNew":        forms[1],
			"eventAggNew":          forms[2],
			"RangeDataList":        event.Club.Mounds,
			"eventRangeUpdate":     forms[3],
			"eventAggUpdate":       forms[4],
			"eventAvailableGrades": forms[5],
		},
	})
}

func dataListRanges(ranges []Range, selected bool) (options []forms.Option) {
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, forms.Option{Label: r.Name, Value: fmt.Sprintf("%d", r.ID), Selected: selected})
		}
	}
	return options
}

func eventDetailsUpsert(f forms.Form) (string, error) {
	eventID := f.Fields[5].Value
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Event{
			ClubID: f.Fields[0].Value,
			Name:   f.Fields[1].Value,
			Date:   f.Fields[2].Value,
			Time:   f.Fields[3].Value,
			Closed: f.Fields[4].Checked,
		}, &Event{}, updateEventDetails)
}

func eventRangeInsert(f forms.Form) (string, error) {
	eventID := f.Fields[1].Value
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{Name: f.Fields[0].Value}, &Event{}, eventAddRange)
}

func eventRangeUpdate(f forms.Form) (string, error) {
	eventID := f.Fields[0].Value
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:     f.Fields[1].ValueUint,
			Name:   f.Fields[2].Value,
			Locked: f.Fields[3].Checked,
			Order:  f.Fields[4].ValueUint,
		}, &Event{}, editRange)
}

func eventAggUpdate(f forms.Form) (string, error) {
	eventID := f.Fields[0].Value
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:    f.Fields[1].ValueUint,
			Name:  f.Fields[2].Value,
			Aggs:  f.Fields[3].ValueUintSlice,
			Order: f.Fields[4].ValueUint,
		}, &Event{}, editRange)
}

func eventAggInsert(f forms.Form) (string, error) {
	eventID := f.Fields[2].Value
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			Name:  f.Fields[0].Value,
			Aggs:  f.Fields[1].ValueUintSlice,
			IsAgg: true,
		}, &Event{}, eventAddAgg)
}
