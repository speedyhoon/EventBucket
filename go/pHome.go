package main

import (
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	formatYMD  = "2006-01-02"
	formatTime = "15:04"
	formatGMT  = "Mon, 02 Jan 2006 15:04:05 GMT"
)

func home(w http.ResponseWriter, r *http.Request) {

	myForm := getSession(w, r)
	existingFields := myForm.fields

	if len(existingFields) <= 0 {
		warn.Println("no fields present")
		var noEvents bool
		noEvents = true
		existingFields = []field{
			{},
			{Required: noEvents},
			{Value: time.Now().Format(formatYMD)},
			{Value: time.Now().Format(formatTime)},
			/*{
				Error: "This is error on a search bar.",
				Options: []option{
					{Label: "label", Value: "2 3"},
					{Label: "text", Value: `"T`},
					{Label: "Search", Value: ">S"},
				},
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
			},*/
		}
	} else {
		info.Println("found something")
	}

	//	var myForm []input
	//	info.Println("sessionForm", sessionForm)
	//	cookies, err := r.Cookie("z")
	//	info.Println("r.Cookie=", cookies, err)

	//	w.Header().Del("z=fdsa")

	//	w.Header().Set("Set-Cookie", fmt.Sprintf("z=; path=/; expires=%v", time.Now().UTC().Add(-5*time.Minute).Format(gmtFormat)))

	/*if err != nil {



		warn.Println("cookiw error", err)

		myForm = []input{
			{}, {}, {}, {},
		}
	} else {
		info.Println("cookies", cookies.Name)
		info.Println("cookies", cookies.Value)

		myForm = []input{
			{
				Error: "1",
			}, {
				Error: "2",
			}, {
				Error: "3",
			}, {
				Error: "4",
			},
		}
	}*/
	//	var errorForm []input

	//	var ok bool

	//	for index, cookie := range cookies {
	//		cookie = strings.TrimPrefix(cookie, "z=")
	//		if cookie != cookies[index] {
	//			errorForm, ok = globalSessions[cookie]
	//			if !ok {
	//				warn.Println("oops con't find that one :(")
	//			}
	//			return
	//		}
	//	}
	//

	templater(w, page{
		Title: "Home",
		Data: M{
			//			"MyForm": errorForm,
			//"MyForm": myForm,
			"MyForm": existingFields,
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

func clubs(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Clubs",
		Data: M{
			"Default": []field{
				{
					Error:   "hey hey!",
					Value:   "true",
					Options: []option{},
				}, {
					Error: "I caused an error!@",
					Value: "fds",
					Options: []option{
						{Label: "label", Value: "2 3"},
						{Label: "text", Value: `"T`},
						{Label: "search", Value: ">S"},
					},
				}, {
					Value: "AbC",
				},
			},
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
					Error: "I caused an error!@",
					Options: []option{
						{Label: "label", Value: "2 3"},
						{Label: "text", Value: `"T`},
						{Label: "search", Value: ">S"},
					},
				},
				{
					Options: []option{
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
	//All EventBucket page urls and ids are lowercase
	lowerURL := strings.ToLower(strings.TrimSuffix(r.URL.Path, "/"))

	//prevents a redirect loop if url is already in lowercase letters.
	if r.URL.Path != lowerURL {

		//check if the request matches any of the pages that don't require parameters
		if strings.Count(lowerURL, "/") >= 2 {
			for _, page := range []string{urlAbout, urlArchive, urlClubs /*urlEvent,*/, urlEvents, urlLicence, urlShooters} {
				if strings.HasPrefix(lowerURL, page) {
					//redirect to page without parameters
					http.Redirect(w, r, page, http.StatusSeeOther)
					return
				}
			}
		}
		http.Redirect(w, r, lowerURL, http.StatusSeeOther)
		return
	}
	w.WriteHeader(status)
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
