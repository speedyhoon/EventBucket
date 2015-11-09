package main

import "net/http"

func home(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		title: "Home",
		data: M{
			"Stuff": "Hommmer page!",
		},
	}, w)
}

func event(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		title: "Event",
		data: M{
			"Stuff": "EVENT page!",
		},
	}, w)
}
