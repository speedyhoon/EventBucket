package main

import (
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r)
	templater(w, page{
		Title: "Home",
		Data: M{
			"NewEvent": eventNewDefaultValues(sessionForm),
		},
	})
}

func report(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Report",
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

func about(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	templater(w, page{
		Title: "About",
		Data: M{
			"Hostname":    hostname,
			"IpAddresses": localIPs(),
		},
	})
}

func licence(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Licence",
	})
}

/*
Adding forms to a page:
	in a HTML form there are 3 main areas:
		the data being displayed - textbox values and select box options
		validation - attributes to hinder user from submitting invalid form data - required, min, max etc...
		presentation & functionality
		error handling - bypassing validation or complex validation not implemented in HTML (required checkbox group)


	it's difficult to maintain standardised forms in all areas without making a standard form builder
	building a form during runtime is quite slow string & slice concatenation
	passing the form to validation has a lot of fields that aren't needed

create the HTML
add validation struct
add population data - map[string]interface OR anonymous struct
*/
