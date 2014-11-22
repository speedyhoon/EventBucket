package main
//gofmt -comments=false -tabwidth=1 -s -e pages.go
import (
	"net/http"
	"fmt"
	"mgo"
	"strings"
)

type list_of_clubs struct {
	Name, Url string
}

func obId(in interface{})string{
	return fmt.Sprintf("%v", in)[13:37]
}

func eventList() []list_of_clubs {
	events := getCollection("event")
	event_list := []list_of_clubs{}
	name := ""
	url := ""
	for _, row := range events {
		name = exists(row, schemaNAME)
		url  = exists(row, "_id")
		fmt.Println(name)
		if name != "" && url != "" {
			event_list = append(event_list, list_of_clubs{
			Name: name,
//				Name: fmt.Sprintf("%v", row["name"]),
			Url: "/event/"+url,
				//			Url: fmt.Sprintf("/event/%v", obId(row["_id"])),
				//				Url: "/event/fdafdas",
				//			Id: fmt.Sprintf("%v", fmt.Sprintf("%v", row["_id"])[13:37]),
				//			Id: bson.ObjectId(row["_id"]),
			})
		}
	}
	return event_list
}








func try(w http.ResponseWriter, r *http.Request){
	var newEvent = map[string]interface{}{
		"_id": "abc",
		"Ranges": map[string]interface{}{
			"0":	 map[string]interface{}{
				"Name": "hello",
				"Type": "agg",
			},
			"1":	 map[string]interface{}{
				"Name": "Neo",
				"Type": "range",
			},
		},
		"Shooters": map[string]interface{}{
			"0":	 map[string]interface{}{
				"Class": "Target",
				"Grade": "A",
			},
			"1":	 map[string]interface{}{
				"Class": "Fclass",
				"Grade": "FA",
			},
		},
	}
//	checkErr(conn.C("event").Upsert(newEvent))
//	newEvent := make(map[string]interface{})
	checkErr(conn.C("event").Find(map[string]interface{}{"_id":"abc"}).One(newEvent))
	fmt.Print(newEvent)

	for index, row := range newEvent{
		fmt.Print("\n\n\n")
		fmt.Print(index)
		fmt.Print(":\t")
		fmt.Print(row)
	}

//	newEvent["Ranges"]["2"] = map[string]interface{}{  //append(newEvent["Ranges"],
//		"Name": "300",
//		"Type": "agg",
//	}
//	checkErr(conn.C("event").Update(newEvent))


//	var temp map[string]interface{}
//	temp = newEvent["Ranges"].(map[string]interface{})
//	temp["2"] = map[string]interface{}{  //append(newEvent["Ranges"],
//		"Name": "300",
//		"Type": "agg",
//	}
	newEvent["Ranges"].(map[string]interface{})["3"] = map[string]interface{}{
		"Name": "nice!!",
		"Type": "range",
	}


	fmt.Print("\n\n\n")
	fmt.Print(newEvent["_id"])



//	var result map[string]interface{}
//	change := mgo.Change{
////		Update: map[string]interface{}{"$push": map[string]interface{}{"ranges":  map[string]interface{}{ "name": range_name, "type": range_type  } }},
//		Update: map[string]interface{}{"$push": map[string]interface{}{"ranges":  map[string]interface{}{ "name": range_name, "type": range_type  } }},
//		ReturnNew: true,
//	}
//	conn.C("event").Find(map[string]interface{}{"_id": newEvent["_id"]  }).Apply(change, &result)
	conn.C("event").UpsertId(newEvent["_id"], newEvent)
}






func clubInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tryThis := validInsert(r.Form, organisers_clubForm())
	InsertDoc(tryThis, "club")
}
func eventInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
//	tryThis := valid8(r.Form, organisers_eventForm(getCollection("club")))
	tryThis := validInsert(r.Form, organisers_eventForm(getCollection("club")))
	tryThis[schema("autoinc")] = map[string]interface{}{
		schemaRANGE: 0,
//		schemaSHOOTER: 0,
	}
	InsertDoc(tryThis, "event")
}
func champInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
//	tryThis := valid8(r.Form, organisers_champForm())
	tryThis := validInsert(r.Form, organisers_champForm())
	InsertDoc(tryThis, "champ")
}
func rangeInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Print("\n\n\n")
	fmt.Print(r)

//	tryThis := valid8(r.Form, eventSettings_add_rangesForm("fdsafas"))
	tryThis := validInsert(r.Form, eventSettings_add_rangesForm("fdsafas"))



//	fmt.Print("\n\n\n")
//	fmt.Print(tryThis)

//	fmt.Print("\n\n\n")
//	fmt.Print(tryThis[schemaID])

//	url := "/eventSettings/"+tryThis[schemaID].(string)
	url := fmt.Sprintf("/eventSettings/%v",tryThis[schemaID])


	redirecter(url, w, r)
	appendRange(tryThis[schemaID].(string), tryThis[schemaNAME].(string), "")
}
func shooterInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tryThis := validInsert(r.Form, eventSettings_add_shooterForm("e"))

	tryThis = convert_class_grade(tryThis)

	event_id := tryThis[schemaID].(string)
//	url := fmt.Sprintf("/eventSettings/%v",event_id)
	redirecter(fmt.Sprintf("/eventSettings/%v",event_id), w, r)

	tryThis[schemaID] = InsertDoc(tryThis, "shooter")

	appendShooter(event_id, tryThis)
}

func convert_class_grade(event map[string]interface{})map[string]interface{}{
	/*converts map key and attribute from:
		"classgrade": "target,B",
	to
		schemaCLASS: "target",
		schemaGRADE: "B",
	*/
	classgrade := strings.Split(event["classgrade"].(string), ",")
	event[schemaCLASS] = classgrade[0]
	event[schemaGRADE] = classgrade[1]
	delete(event, "classgrade")
	return event
}



func clubRangeInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//	tryThis := valid8(r.Form, organisers_champForm())
	tryThis := validInsert(r.Form, clubRangeInsertForm())
	dump(tryThis)
	club_id := tryThis["clubid"].(string)
	dump(club_id)
	redirecter(fmt.Sprintf("/club/%v",club_id), w, r)
	var club map[string]interface{}
	change := mgo.Change{
		Update: map[string]interface{}{
			"$inc": map[string]interface{}{fmt.Sprintf("%v.%v",schemaAUTOINC,schemaRANGE): 1},
		},
		ReturnNew: true,
	}
	conn.C("club").FindId(club_id).Apply(change, &club)
	dump("00000000000000000000000000000000000")
	range_id := club[schemaAUTOINC].(map[string]interface{})[schemaRANGE]

	delete(tryThis, "clubid")

	change = mgo.Change{
		Update: map[string]interface{}{
			"$set": map[string]interface{}{fmt.Sprintf("%v.%v",schemaRANGE,range_id): tryThis},
		},
//		ReturnNew: true,
	}
	conn.C("club").FindId(club_id).Apply(change, &club)
}


func updateTotalScores(w http.ResponseWriter, r *http.Request){
	temp := map[string]interface{}{
		"total": "T",
		"center": "C",
	}
	dump(temp["total"])









//
//
//
//	r.ParseForm()
//
//	tryThis := validInsert(r.Form, totalScores_update("e","e"))
//	separator_setting := "."
//	dump("UPDATE")
//	dump(tryThis)
//
//	scores := strings.Split(tryThis[schema("total")].(string), separator_setting)
//	dump("fdsfdsafdsafdsafsdafsa")
//	dump(scores)
//	dump("fdsfdsafdsafdsafsdafsa")
//
//
////	tryThis = tryThis.(map[string]interface{})
//	event_id := tryThis["eventId"]//.(string)
//	shooter_id := tryThis["shooterId"]//.(string)
//	range_id := tryThis["rangeId"]//.(string)
//	total := scores[0]//tryThis[schema("total")]//.(string)
//	center := scores[1]//tryThis[schema("center")]//.(string)
//
//	dump(event_id)
//	dump(shooter_id)
//	dump(range_id)
//	dump(total)
//	dump(center)
//
//	redirecter(fmt.Sprintf("/totalScores/%v/%v", event_id, range_id), w, r)
//
//	change := mgo.Change{
//		Update: map[string]interface{}{
//			"$set": map[string]interface{}{
//				fmt.Sprintf("%v.%v.%v",schema("shooter"),shooter_id,range_id): map[string]interface{}{
//					schema("total"): total,
//					schema("center"): center,
//				},
//			},
//		},
//		ReturnNew: true,
//	}
//	var event map[string]interface{}
//	conn.C("event").FindId(event_id).Apply(change, &event)
}
