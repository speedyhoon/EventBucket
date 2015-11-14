package main

import "net/http"

func home(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Home",
		Data: M{
			"Stuff": "Hommmer page!",
		},
	})
}

func event(w http.ResponseWriter, r *http.Request, eventId string) {
	templater(w, page{
		Title: "Event",
		Data: M{
			"Stuff":   "EVENT page!",
			"EventId": eventId,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Events",
		Data: M{
			"Stuff": "EVENTS page!",
		},
	})
}

func eventArchive(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Archive",
		Data: M{
			"Stuff": "Archive page!",
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Clubs",
		Data: M{
			"Stuff": "CLUBS page!",
		},
	})
}

func shooters(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Shooters",
		Data: M{
			"Stuff": "SHOOTERS page!",
		},
	})
}

func about(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "About",
		Data: M{
			"Stuff": "About page!",
		},
	})
}

func licence(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Licence",
		Data: M{
			"Stuff": "Licence page!",
		},
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	//if status == http.StatusNotFound {
	//	fmt.Fprint(w, "custom 404")
	//}
	templater(w, page{
		Title: "Error",
		Data: M{
			//		"Status": "404 Page Not Found",
			"Status": status,
		},
	})
}
