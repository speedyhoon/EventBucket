package main

import (
	"net/http"
	"fmt"
	//	"os"

	//	"html/template"
	"regexp"
	"mgo/bson"
//	"mgo"
)

func organisers(w http.ResponseWriter, r *http.Request) {
	templator("admin", organisers_HTML(), organisers_Data(), w)
}

type list_of_clubs struct {
	Name, Id string
}

func organisers_Data() map[string]interface{} {
	clubs := getCollection("club")
	club_list := []list_of_clubs{}
	for _, row := range clubs {
		club_list = append(club_list, list_of_clubs{
			Name: fmt.Sprintf("%v", row["name"]),
			Id: fmt.Sprintf("%v", fmt.Sprintf("%v", row["_id"])[11:34]),
		})
	}
	return map[string]interface{}{
		"Title": "Organisers",
		"Events": generateForm("eventInsert", organisers_eventForm(clubs)),
		"EventList": eventList(getCollection("event")),
		"Clubs": generateForm("clubInsert", organisers_clubForm()),
		"ClubList": club_list,
		"Championship": generateForm("champInsert", organisers_champForm()),
	}
}

func eventList(events []map[string]interface{})[]list_of_clubs{
	event_list := []list_of_clubs{}
	for _, row := range events {
		event_list = append(event_list, list_of_clubs{
			Name: fmt.Sprintf("%v", row["name"]),
//			Id: fmt.Sprintf("%v", fmt.Sprintf("%v", row["_id"])[13:37]),


//			Id: bson.ObjectId(row["_id"]),
		})
	}
	return event_list
}
//ObjectId("52a9a1ffff7f0c7aacacbe09")
func organisers_HTML() string {
	return `<h1>{{ .Title}}</h1>
	{{if .Events}}` +
		panel("Events", pane("Create Event", `{{XTC .Events}}`) +
		`{{if .EventList}}
			<table><tr><th>Name</th></tr>
				{{with .EventList}}
					{{range .}}
						<tr><td><a href=`+addQuotes(`{{.Url}}`)+`>{{.Name}}</a></td></tr>
					{{end}}
				{{end}}
			</table>
		{{end}}`) +
	`{{end}}
	{{if .Clubs}}` +
		panel("Clubs", pane("Create Club", `{{XTC .Clubs}}`) +
			`{{if .ClubList}}` +
				pane("Create Club",
					`<table><tr><th>Name</th></tr>
						{{with .ClubList}}
							{{range .}}
								<tr><td><a href=`+addQuotes(`{{.Url}}`)+`>{{.Name}}</a></td></tr>
							{{end}}
						{{end}}
					</table>`)+`{{end}}`) +
	`{{end}}
	{{if .Championship}}` +
		panel("Championships", pane("Create Championship", `{{XTC .Championship}}`)) +
	`{{end}}`
}

func organisers_clubForm() map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Club Name",
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Add Club",
		},
	}
}

func organisers_champForm() map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Championship Name",
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Add Championship",
		},
	}
}

func organisers_eventForm(eventsCollection []map[string]interface{}) map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Event Name",
		},
		"club":Inputs{
			Html:      "select",
			SelectValues:   getClubSelectBox(eventsCollection),
			Label:   "Event Name",
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Add Event",
		},
	}
}


func tempTry(w http.ResponseWriter, r *http.Request){
//	r.ParseForm()
//	tryThis := valid8(r.Form, organisers_clubForm())


//	temp3 := []interface{}{bson.M{"_id": "club"}, bson.M{"next":1}	}




	conn.C("club").Insert(bson.M{"_id":"autoInc(`club`)", "name":"one"})
//	conn.C("club").Insert(bson.M{"_id":autoInc("club"), "name":"two"})
//	conn.C("club").Insert(bson.M{"_id":autoInc("club"), "name":"three"})



	var result map[string]interface{}
	checkErr(conn.C("autoinc").Find( bson.M{"_id": "club", "next":1} ).One(&result))

	for index, row := range result{
		fmt.Printf("%v",index)
		fmt.Printf("\t\t")
		fmt.Printf("%v",row)
		fmt.Printf("\n")

	}


//
//
//
//	function counter(name) {
//	var ret = db.counters.findAndModify({query:{_id:name}, update:{$inc : {next:1}}, "new":true, upsert:true});
//// ret == { "_id" : "users", "next" : 1 }
//return ret.next;
//}
//
//db.users.insert({_id:counter("users"), name:"Sarah C."}) // _id : 1
//db.users.insert({_id:counter("users"), name:"Bob D."}) // _id : 2


	//	tryThis["_id"] = result["next"] + 1

}


func clubInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tryThis := valid8(r.Form, organisers_clubForm())
	InsertDoc(tryThis, "club")
}
func eventInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tryThis := valid8(r.Form, organisers_eventForm(getCollection("club")))
	InsertDoc(tryThis, "event")
}
func champInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tryThis := valid8(r.Form, organisers_champForm())
	InsertDoc(tryThis, "champ")
}
func rangeInsert(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	tryThis := valid8(r.Form, eventSettings_add_rangesForm())
	//select event by id
	//append new range to event
	//save changes
	InsertDoc(tryThis, "event")
}


func events(w http.ResponseWriter, r *http.Request){
	templator("admin", events_HTML(), events_Data(), w)
}
func events_Data() map[string]interface{} {
	return map[string]interface{}{
		"Title": "Event List",
		"EventList": eventList(getCollection("event")),
	}
}

func events_HTML() string {
	return `<h1>{{ .Title}}</h1>
	{{if .EventList}}
		<table><tr><th>Name</th></tr>
			{{with .EventList}}
				{{range .}}
					<tr><td><a href=`+addQuotes(`/event/{{.Id}}`)+`>{{.Name}}</a></td></tr>
				{{end}}
			{{end}}
		</table>
	{{end}}`
}







func event(w http.ResponseWriter, r *http.Request) {
	page_url := "/event/"
//	var validID = regexp.MustCompile(`\A`+page_url+`[0-9a-f]{24}\z`)
	url := fmt.Sprintf("%v",r.URL)

	fmt.Println(url)

//	if validID.MatchString(url){

		templator("admin", event_HTML(), event_Data(url[len(page_url):]), w)
//	}else{
//		redirectPermanent("/events")
//		fmt.Println("redirected user "+url)
//	}
}

func event_Data(event_id string) map[string]interface{} {
//	event := getDocument("event",event_id)

//	event := getShit("event",event_id)

	event := getCollection("event")
	fmt.Println(event_id)
	fmt.Println(event)
	fmt.Printf("%v",event)
	return map[string]interface{}{
		"Title": `event["name"]`,
		"Id": event_id,
	}
}

func event_HTML() string {
	return `<h1>{{ .Title}}</h1>
	<a href=`+addQuotes(`/eventSettings/{{.Id}}`)+`>Event Settings</a>`
}


func eventSettings(w http.ResponseWriter, r *http.Request) {
	page_url := "/eventSettings/"
	var validID = regexp.MustCompile(`\A`+page_url+`[0-9a-f]{24}\z`)
	url := fmt.Sprintf("%v",r.URL)
	if validID.MatchString(url){
		templator("admin", eventSettings_HTML(), eventSettings_Data(url[len(page_url):]), w)
	}else{
		redirectPermanent("/events")
		fmt.Println("redirected user "+url)
	}
}
func eventSettings_Data(event_id string) map[string]interface{} {
	return map[string]interface{}{
		"Title": "Event Settings",
		"Id": event_id,
		"AddRange": generateForm("rangeInsert", eventSettings_add_rangesForm()),
	}
}
func eventSettings_HTML() string {
	return `<h1>{{ .Title}}</h1>
	<a href=`+addQuotes(`/eventSettings/{{ .Id}}`)+`>Event Settings</a>
	`+panel("Ranges",pane("Add Ranges",`{{XTC .AddRange}}`))
}
func eventSettings_add_rangesForm() map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Range Name",
		},
//		"club":Inputs{
//			Html:      "select",
//			SelectValues:   getClubSelectBox(eventsCollection),
//			Label:   "Event Name",
//		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Add Range",
		},
	}
}
