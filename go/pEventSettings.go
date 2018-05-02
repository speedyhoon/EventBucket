package main

import (
	"fmt"
	"net/http"

	"github.com/speedyhoon/forms"
	"github.com/speedyhoon/session"
)

func eventSettings(w http.ResponseWriter, r *http.Request, event Event) {
	//Retrieve any submitted form that failed to save.
	fs, action := session.Forms(w, r, getFields, eventEdit, eventRangeNew, eventAggNew, eventRangeEdit, eventAggEdit, eventAvailableGrades)
	if action != eventEdit {
		fs[action].Fields[0].Value = event.Club.Name
		fs[action].Fields[1].Value = event.Name
		fs[action].Fields[2].Value = event.DateTime
		fs[action].Fields[3].Value = event.Closed
		fs[action].Fields[4].Value = event.ID
	}
	fs[eventRangeNew].Fields[1].Value = event.ID

	fs[eventAggNew].Fields[1].Options = dataListRanges(event.Ranges, true)
	fs[eventAggNew].Fields[2].Value = event.ID

	fs[eventAvailableGrades].Fields[0].Options = availableGrades(event.Grades)
	fs[eventAvailableGrades].Fields[1].Value = event.ID

	render(w, page{
		Title:   "Event Settings",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Ranges":               dataListRanges(event.Ranges, false),
			"Event":                event,
			"eventEdit":            fs[eventEdit],
			"eventRangeNew":        fs[eventRangeNew],
			"eventAggNew":          fs[eventAggNew],
			"RangeDataList":        event.Club.Mounds,
			"eventRangeUpdate":     fs[eventRangeEdit],
			"eventAggUpdate":       fs[eventAggEdit],
			"eventAvailableGrades": fs[eventAvailableGrades],
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
	eventID := f.Fields[4].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Event{
			ClubID:   f.Fields[0].Str(),
			Name:     f.Fields[1].Str(),
			DateTime: f.Fields[2].Time(),
			Closed:   f.Fields[3].Checked(),
		}, &Event{}, updateEventDetails)
}

func eventRangeInsert(f forms.Form) (string, error) {
	eventID := f.Fields[1].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{Name: f.Fields[0].Str()}, &Event{}, eventAddRange)
}

func eventRangeUpdate(f forms.Form) (string, error) {
	eventID := f.Fields[0].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:     f.Fields[1].Uint(),
			Name:   f.Fields[2].Str(),
			Locked: f.Fields[3].Checked(),
			Order:  f.Fields[4].Uint(),
		}, &Event{}, editRange)
}

func eventAggUpdate(f forms.Form) (string, error) {
	eventID := f.Fields[0].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			ID:    f.Fields[1].Uint(),
			Name:  f.Fields[2].Str(),
			Aggs:  f.Fields[3].UintSlice(),
			Order: f.Fields[4].Uint(),
		}, &Event{}, editRange)
}

func eventAggInsert(f forms.Form) (string, error) {
	eventID := f.Fields[2].Str()
	return urlEventSettings + eventID,
		updateDocument(tblEvent, eventID, &Range{
			Name:  f.Fields[0].Str(),
			Aggs:  f.Fields[1].UintSlice(),
			IsAgg: true,
		}, &Event{}, eventAddAgg)
}
