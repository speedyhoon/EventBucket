package main

import (
	"net/http"
	"sort"
	"strings"

	"github.com/speedyhoon/session"
)

func events(w http.ResponseWriter, r *http.Request) {
	events, err := getEvents(onlyOpen)
	if err == nil {
		//Sort list of events by date then by name
		eventOrderedBy(sortByDate, sortByName).Sort(events)
	}
	fs, _ := session.Forms(w, r, getFields, eventNew)

	render(w, page{
		Title: "Events",
		Error: err,
		Data: map[string]interface{}{
			"eventNew": fs[eventNew],
			"Events":   events,
			"Network":  localIPs(),
		},
	})
}

func onlyOpen(event Event) bool {
	return !event.Closed
}

//TODO change sort form true/false to 1/0/-1
type compareEvent func(p1, p2 *Event) bool

type eventSort struct {
	events []Event
	less   []compareEvent
}

func (ms *eventSort) Sort(events []Event) {
	ms.events = events
	sort.Sort(ms)
}

func (ms *eventSort) Len() int {
	return len(ms.events)
}

func (ms *eventSort) Swap(i, j int) {
	ms.events[i], ms.events[j] = ms.events[j], ms.events[i]
}

func (ms *eventSort) Less(i, j int) bool {
	p, q := &ms.events[i], &ms.events[j]
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

func eventOrderedBy(less ...compareEvent) *eventSort {
	return &eventSort{
		less: less,
	}
}

func sortByDate(c1, c2 *Event) bool {
	return c1.DateTime.After(c2.DateTime)
}
func sortByName(c1, c2 *Event) bool {
	return strings.ToLower(c1.Name) < strings.ToLower(c2.Name)
}
