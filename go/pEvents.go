package main

import (
	"net/http"
	"sort"
	"strings"
	"time"
)

func events(w http.ResponseWriter, r *http.Request) {
	events, err := getCalendarEvents()
	if err == nil {
		//Sort list of events by date then by name
		orderedByEvent(sortByDate, sortByName).Sort(events)
	}
	_, forms := sessionForms(w, r, eventNew)

	templater(w, page{
		Title: "Events",
		Error: err,
		Data: map[string]interface{}{
			"eventNew": forms[0],
			"Events":   events,
			"Network":  localIPs(),
		},
	})
}

//CalendarEvent is the same as Event struct without Shooters and their scores.
type CalendarEvent struct {
	ID     string `json:"I"`
	Name   string `json:"n"`
	ClubID string `json:"C,omitempty"`
	Club   string `json:"c,omitempty"`
	Date   string `json:"d,omitempty"`
	Time   string `json:"t,omitempty"`
	ISO    time.Time
	Ranges []Range `json:"R,omitempty"`
	Closed bool    `json:"z,omitempty"`
}

//TODO change sort form true/false to 1/0/-1
type lessFunc2 func(p1, p2 *CalendarEvent) bool

type multiSorter2 struct {
	changes []CalendarEvent
	less    []lessFunc2
}

func (ms *multiSorter2) Sort(changes []CalendarEvent) {
	ms.changes = changes
	sort.Sort(ms)
}

func (ms *multiSorter2) Len() int {
	return len(ms.changes)
}

func (ms *multiSorter2) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

func (ms *multiSorter2) Less(i, j int) bool {
	p, q := &ms.changes[i], &ms.changes[j]
	//Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}

func orderedByEvent(less ...lessFunc2) *multiSorter2 {
	return &multiSorter2{
		less: less,
	}
}

func sortByDate(c1, c2 *CalendarEvent) bool {
	return c1.ISO.After(c2.ISO)
}
func sortByName(c1, c2 *CalendarEvent) bool {
	return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
}
