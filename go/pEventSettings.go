package main

import (
	"fmt"
	"net/http"
	"github.com/speedyhoon/session"
	"github.com/speedyhoon/forms"
)

func eventSettings(w http.ResponseWriter, r *http.Request, event Event) {
	//Retrieve any submitted form that failed to save.
	f, submitted := session.Forms(w, r, getFields, eventEdit, eventRangeNew, eventAggNew, eventRangeEdit, eventAggEdit, eventAvailableGrades)
	if submitted.Action != eventEdit {
		submitted.Fields[0].Value = event.Club.Name
		submitted.Fields[1].Value = event.Name
		submitted.Fields[2].Value = event.Date
		submitted.Fields[3].Value = event.Time
		submitted.Fields[4].Checked = event.Closed
		submitted.Fields[5].Value = event.ID
	}
	f[1].Fields[1].Value = event.ID

	f[2].Fields[1].Options = dataListRanges(event.Ranges, true)
	f[2].Fields[2].Value = event.ID

	//AvailableGrades
	f[5].Fields[0].Options = availableGrades(event.Grades)
	f[5].Fields[1].Value = event.ID

	render(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":               dataListRanges(event.Ranges, false),
			"Event":                event,
			"eventEdit":            f[0],
			"eventRangeNew":        f[1],
			"eventAggNew":          f[2],
			"RangeDataList":        event.Club.Mounds,
			"eventRangeUpdate":     f[3],
			"eventAggUpdate":       f[4],
			"eventAvailableGrades": f[5],
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
