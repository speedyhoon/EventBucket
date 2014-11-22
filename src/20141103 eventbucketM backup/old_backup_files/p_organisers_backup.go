package main

import (
	"net/http"
	"net/url"
//	"mgo"
//	"mgo/bson"
	"fmt"
)

func organisers(w http.ResponseWriter, r *http.Request) {
	html := loadHTM("organisers")
	template_data := organisers_Data()
	templator("admin", html, template_data, w)
}

func organisers_Data() map[string]interface{} {
	clubs := getCollection("club")
	club_list := []list_of_clubs{}
	name := ""
	url := ""
	for _, row := range clubs {
		name = exists(row, schemaNAME)
		url  = exists(row, "_id")
		if name != "" && url != "" {
			club_list = append(club_list, list_of_clubs{
				Name: name,
				Url: url,
			})
		}
	}
	return map[string]interface{}{
		"Title":  "Organisers",
		"Events": generateForm("eventInsert", organisers_eventForm(clubs)),
		"EventList": eventList(),
		"Clubs": generateForm("clubInsert", organisers_clubForm()),
		"ClubList": club_list,
		"Championship": generateForm("champInsert", organisers_champForm()),
	}
}

func organisers_clubForm() map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Club Name",
		},
		"submit": {
			Html:  "submit",
			Label: "Add Club",
		},
	}
}

func organisers_champForm() map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Championship Name",
		},
		"submit": {
			Html:  "submit",
			Label: "Add Championship",
		},
	}
}

func organisers_eventForm(eventsCollection []map[string]interface{}) map[string]Inputs {
	return map[string]Inputs{
		"name": {
			Html:  "text",
			Label: "Event Name",
			Required: true,
		},
		"club": {
			Html:         "select",
			SelectValues: getClubSelectBox(eventsCollection),
			Label:        "Host Club",
			Required: true,
		},
		"submit": {
			Html:  "submit",
			Label: "Add Event",
		},
	}
}
func clubInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	var newClub Club
	options := organisers_clubForm()

//	tryThis := make(map[string]interface{})
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
//				dump(array)
				if option == "name"{
					newClub.Name = fmt.Sprintf("%v",array[0])
				}
//				tryThis[option] = array[0]
			}else {
				log("\nELSE options[%v] is REQUIRED", option)
				log("OR", nil)
				log("\nELSE options[%v] not in array ", option)
			}
		}
	}


	newClub.Id = getNextId("club")
//	newClub.Rego = "ABCDEG"
//	if form["name"] != nil{
//		newClub.name = form["name"]
//	}

//	return tryThis


//	tryThis := validInsert3(r.Form, organisers_clubForm())

//	dump(newClub)
//	for index, attribute := range tryThis{
//	}


	dump(newClub)
	err := conn.C("club").Insert(newClub)
	checkErr(err)
}

type Club struct{
//	Id bson.ObjectId `bson:"_id"`
	Id string `bson:"_id"`
	Name string `bson:"n"`
//	Rego string `bson:"r,omitempty"`
	Url string `bson:"u,omitempty"`
	Ranges map[string]Ranges `bson:"R,inline,omitempty"`
}

func eventInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//	tryThis := valid8(r.Form, organisers_eventForm(getCollection("club")))
	tryThis := validInsert3(r.Form, organisers_eventForm(getCollection("club")))
	tryThis[schemaAUTOINC] = map[string]interface{}{
		schemaRANGE: 0,
		//		schemaSHOOTER: 0,
	}

	InsertDoc2(tryThis, "event")
}
type Event struct{
	Id string `bson:"_id"`
	Club string `bson:"c"`
	Name string `bson:"n"`
	Datetime string `bson:"d"`
	Ranges map[string]Ranges `bson:"R,inline,omitempty"`
	Shooters map[string]EventShooter `bson:"S,inline,omitempty"`
	TeamCat map[string]TeamCat `bson:"A,inline,omitempty"`
	Teams map[string]Team `bson:"T,inline,omitempty"`
}




func validInsert3(form url.Values, options map[string]Inputs) map[string]interface{} {
	tryThis := make(map[string]interface{})
	for option := range options {
		if options[option].Html != "submit" {
			array, ok := form[option]
			if ok && ((options[option].Required && array[0] != "") || !options[option].Required) {
//				dump(array)
				tryThis[option] = array[0]
			}else {
				log("\nELSE options[%v] is REQUIRED", option)
				log("OR", nil)
				log("\nELSE options[%v] not in array ", option)
			}
		}
	}
	return tryThis
}


func InsertDoc2(data map[string]interface{}, collection_name string){
//	log("%v",data)
//

	//	if data != false {
	//		data = [...]map[string]interface{}{"_id":getNextId(collection_name)}
	err := conn.C(collection_name).Insert(data)
	checkErr(err)
	//	}
//	return data["_id"].(string)
}


type Ranges struct{
	id int
	Name string `bson:"N"`
	Type string `bson:"T"`
	Aggs []int `bson:"a"`
	ScoreBoard bool
	Enabled bool
}
