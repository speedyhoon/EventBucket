package main

import (
	"net/http"
	"os"
)

func insertEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		/*405 Method Not Allowed
		A request was made of a resource using a request method not supported by that resource; for example,
		using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
		//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
	} else {
		r.ParseForm()
		if len(r.Form) != 4 {
			info.Println("invalid number of form items")

			w.Header().Add("Set-Cookie", sessionError("invalid number of form items"))

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/event/1a", http.StatusSeeOther)
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	info.Println("globalSessions", globalSessions)
	cookies := r.Cookies()
	info.Println("cookies", cookies)

	templater(w, page{
		Title: "Home",
		Data: M{
			"Stuff": "Hommmer page!",
			"MyForm": []input{
				{
					Error: "This is error on a search bar.",
					//					Options: []Option{
					//						{Label: "label", Value: "2 3"},
					//						{Label: "text", Value: `"T`},
					//						{Label: "Search", Value: ">S"},
					//					},
				},
				{
					Error: "Another error on club input.",
					Options: []option{
						{Label: "Warrack", Value: "R23"},
						{Label: "Horsham", Value: "T52"},
						{Label: "Stawell", Value: "S82"},
					},
				},
				{
					Error: "Error on date field!",
				},
				{
					Error: "Error on time field!",
				},
			},
		},
	})
}

func event(w http.ResponseWriter, r *http.Request, eventId string) {

	pageURL := "/event/"
	if r.URL.Path[len(pageURL):] != "1a" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	eventId = "1A"
	templater(w, page{
		Title: "Event",
		Data: M{
			"Stuff":   "EVENT page!",
			"EventId": eventId,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, "/events", http.StatusSeeOther)
	}
	templater(w, page{
		Title: "Events",
		Data: M{
			"Stuff": "EVENTS page!",
		},
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
				"Fds": []input{
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
			"Fds": []field{
				{
					error: "I caused an error!@",
					options: []option{
						{Label: "label", Value: "2 3"},
						{Label: "text", Value: `"T`},
						{Label: "search", Value: ">S"},
					},
				},
				{
					options: []option{
						{Label: "Warrack", Value: "R23"},
						{Label: "Horsham", Value: "T52"},
						{Label: "Stawell", Value: "S82"},
					},
				},
				//			Date{},
				//			Time{},
				//			Check{},
				//			Hidden{},
			},
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
