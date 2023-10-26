package main

import (
	"fmt"
	"net/http"

	"github.com/speedyhoon/frm"
	"github.com/speedyhoon/session"
)

func eventSettings(w http.ResponseWriter, r *http.Request, event Event) {
	// Retrieve any submitted form that failed to save.
	fs, action := session.Get(w, r, frmEventEdit, frmEventRangeNew, frmEventAggNew, frmEventRangeEdit, frmEventAggEdit, frmEventAvailableGrades)
	if action != frmEventEdit {
		fs[action].Fields[0].Value = event.Club.Name
		fs[action].Fields[1].Value = event.Name
		fs[action].Fields[2].Value = event.DateTime
		fs[action].Fields[3].Value = event.Closed
		fs[action].Fields[4].Value = event.ID
	}
	fs[frmEventRangeNew].Fields[1].Value = event.ID

	fs[frmEventAggNew].Fields[1].Options = dataListRanges(event.Ranges, true)
	fs[frmEventAggNew].Fields[2].Value = event.ID

	fs[frmEventAvailableGrades].Fields[0].Options = availableGrades(event.Grades)
	fs[frmEventAvailableGrades].Fields[1].Value = event.ID

	render(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":               dataListRanges(event.Ranges, false),
			"Event":                event,
			"eventEdit":            fs[frmEventEdit],
			"eventRangeNew":        fs[frmEventRangeNew],
			"eventAggNew":          fs[frmEventAggNew],
			"RangeDataList":        event.Club.Mounds,
			"eventRangeUpdate":     fs[frmEventRangeEdit],
			"eventAggUpdate":       fs[frmEventAggEdit],
			"eventAvailableGrades": fs[frmEventAvailableGrades],
		},
	})
}

func dataListRanges(ranges []Range, selected bool) (options []frm.Option) {
	for _, r := range ranges {
		if !r.IsAgg {
			options = append(options, frm.Option{Label: r.Name, Value: fmt.Sprintf("%d", r.ID), Selected: selected})
		}
	}
	return options
}

func eventDetailsUpsert(f frm.Form) (string, error) {
	eventID := f.Fields[4].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Event{
			ClubID:   f.Fields[0].Str(),
			Name:     f.Fields[1].Str(),
			DateTime: f.Fields[2].Time(),
			Closed:   f.Fields[3].Checked(),
		}, &Event{}, updateEventDetails)
}

func eventRangeInsert(f frm.Form) (string, error) {
	eventID := f.Fields[1].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{Name: f.Fields[0].Str()}, &Event{}, eventAddRange)
}

func eventRangeUpdate(f frm.Form) (string, error) {
	eventID := f.Fields[0].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:     f.Fields[1].Uint(),
			Name:   f.Fields[2].Str(),
			Locked: f.Fields[3].Checked(),
			Order:  f.Fields[4].Uint(),
		}, &Event{}, editRange)
}

func eventAggUpdate(f frm.Form) (string, error) {
	eventID := f.Fields[0].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:    f.Fields[1].Uint(),
			Name:  f.Fields[2].Str(),
			Aggs:  f.Fields[3].Uints(),
			Order: f.Fields[4].Uint(),
		}, &Event{}, editRange)
}

func eventAggInsert(f frm.Form) (string, error) {
	eventID := f.Fields[2].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			Name:  f.Fields[0].Str(),
			Aggs:  f.Fields[1].Uints(),
			IsAgg: true,
		}, &Event{}, eventAddAgg)
}
