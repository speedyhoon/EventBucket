package main

import "net/http"

func home(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Home",
		Data: M{
			"Stuff": "Hommmer page!",
		},
	}, w)
}

func event(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Event",
		Data: M{
			"Stuff": "EVENT page!",
		},
	}, w)
}

func events(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Events",
		Data: M{
			"Stuff": "EVENTS page!",
		},
	}, w)
}

func eventArchive(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Archive",
		Data: M{
			"Stuff": "Archive page!",
		},
	}, w)
}

func clubs(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Clubs",
		Data: M{
			"Stuff": "CLUBS page!",
		},
	}, w)
}

func shooters(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Shooters",
		Data: M{
			"Stuff": "SHOOTERS page!",
		},
	}, w)
}

func about(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "About",
		Data: M{
			"Stuff": "About page!",
		},
	}, w)
}

func licence(w http.ResponseWriter, r *http.Request) {
	templater(Page{
		Title: "Licence",
		Data: M{
			"Stuff": "Licence page!",
		},
	}, w)
}
