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

func all(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "_All",
		Data: M{
			"Event": M{
				"Stuff":   "EVENT page!",
				"EventId": "eventId",
			},
			"About": M{
				"Stuff": "EVENTS page!",
			},
			"Clubs": M{
				"Stuff": "CLUBS page!",
			},
			"Shooters": M{
				"Stuff": "SHOOTERS page!",
				"Fds": []field{
					{
						Error: "i caused an error!@",
						Options: []option{
							{Label: "label", Value: "2 3"},
							{Label: "text", Value: `"t`},
							{Label: "search", Value: ">s"},
						},
					},
					{
						Options: []option{
							{Label: "warrack", Value: "r23"},
							{Label: "horsham", Value: "t52"},
							{Label: "stawell", Value: "s82"},
						},
					},
				},
			},
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

/*
type form struct {
	Action string
	Title  string
	Field  []field
	Help   string
	Table  bool
	Id     string
}

type field struct {
	Name, Html, Label, Help, Pattern, Placeholder, AutoComplete string //AutoComplete values can be: "on" or "off"
	Checked, MultiSelect, Required                              bool
	Min, Max                                                    *int
	Size                                                        int
	//	Options                                                     []option
	Step                 float64
	VarType              string //the type of variable to return
	MaxLength, MinLength int    //the length of variable to return
	Error                string
	Snippet              interface{}
	Autofocus            bool
	Action               string //Way to switch the forms action to a different purpose
	Value                interface{}
	AccessKey, Inner, Id string
	DataList             bool
	Class                string
}*/

/*
var globalForm = map[string]func(string) form3{
	"eventSettingsAddRangeForm": func(eventID string) form3 {
		return form3{
			action: "/EventRangeInsert",
			title:  "Add Range",
			fields: []Field2{
				searchbox{
					name:      "name",
					label:     "Range Name",
					error:     "whoops this seems to be an unexpected error :(",
					autoFocus: true,
					required:  true,
				}, submit{
					label: "Create Range",
					name:  "eventid",
					Value: eventID,
				},
			},
		}
	},
}*/

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
