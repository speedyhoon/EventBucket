package main

import (
	"net/http"
//	"regexp"
	"fmt"
)

func eventSettings(w http.ResponseWriter, r *http.Request) {
	page_url := "/eventSettings/"
//	var validID = regexp.MustCompile(`\A` + page_url + `[0-9a-f]{24}\z`)
	url := fmt.Sprintf("%v", r.URL)
//	if validID.MatchString(url) {
//		templator("admin", eventSettings_HTML(), eventSettings_Data(url[len(page_url):]), w)
//	}else {
//		redirectPermanent("/events")
//		fmt.Println("redirected user " + url)
//	}
	templator("admin", eventSettings_HTML(), eventSettings_Data(url[len(page_url):]), w)
}
func eventSettings_Data(event_id string) map[string]interface{} {
	event := getDoc_by_id("event", event_id)
//	dump(event[schemaRANGE])
//	for index, ranger := range event[schemaRANGE].(map[string]Ranges) {
//		if event[schemaRANGE].(map[string]Ranges)[index].Aggs && event[schemaRANGE].(map[string]Ranges)[index].Aggs == "~"{
//			event[schemaRANGE].(map[string]Ranges)[index].Aggs = ""
//		}
//	}

	output := map[string]interface{}{
		"Title": "Event Settings",
		"Id": event_id,
		"AddRange": generateForm("rangeInsert", eventSettings_add_rangeForm(event_id)),
		"AddAgg": generateForm("aggInsert", eventSettings_add_aggForm(event_id)),
		"ListRanges": event[schemaRANGE],
//		"RangeList":/ event_obj.Ranges,
		//		"RangeList": ranges,
		"AddShooter": generateForm("shooterInsert", eventSettings_add_shooterForm(event_id)),
		"ShooterList": event[schemaSHOOTER],
	}
	if event_ranges, ok := event[schemaRANGE]; ok {
		output["RangeList"] = event_ranges
	}
	return output

//	{{if.RangetList}} {{end}}

//
//	fmt.Println("\n\n-")
//
//	//	var event_obj Event
//	event_obj := getEvent(event_id)
//	//
//	fmt.Print(fmt.Sprintf("%v",event_obj.Ranges))
//
//
//
//	//	var event_obj map[string]interface{}
//	//	event_obj = getDocument("event", event_id)
//	//	var ranges map[string]interface{}
//	//	ranges = event_obj["ranges"]
//
//	//	fmt.Print(event_obj["Ranges"])
//
//
//	//	vardump(event_obj)
//
//	var ranges map[string]Ranges
//	for index, row := range event_obj.Ranges{
//		ranges[index] = row //.(string)
//
//	}
//
//	return map[string]interface{}{
//		"Title": "Event Settings",
//		"Id": event_id,
//		"AddRange": generateForm("rangeInsert", eventSettings_add_rangesForm(event_id)),
//		//		"RangeList": event_obj["Ranges"],
//		"RangeList": event_obj.Ranges,
//		//		"RangeList": ranges,
//	}
}
func eventSettings_HTML() string {
	return loadHTM("event-settings")
}
func eventSettings_add_rangeForm(event_id string) map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Range Name",
			Required: true,
		},
		"id":Inputs{
			Html: "hidden",
			Value: event_id,
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Create Range",
		},
	}
}
func eventSettings_add_aggForm(event_id string) map[string]Inputs {
	return map[string]Inputs{
		"name":Inputs{
			Html:      "text",
			Label:   "Aggregate Name",
			Required: true,
		},
		"id":Inputs{
			Html: "hidden",
			Value: event_id,
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Create Aggregate",
		},
	}
}
func eventSettings_add_shooterForm(event_id string) map[string]Inputs {
	return map[string]Inputs{
		"firstname":Inputs{
			Html:      "text",
			Label:   "First Name",
			Required: true,
		},
		"surname":Inputs{
			Html:      "text",
			Label:   "Surname",
			Required: true,
		},
		"club":Inputs{
			Html:      "text",
			Label:   "Club",
			Required: true,
		},
		"id":Inputs{
			Html: "hidden",
			Value: event_id,
			Required: true,
		},
		"classgrade":Inputs{
			Html:		"select",
			Label: "Class & Grade",
			Required: true,
			//TODO make these dynamic from club settings
			SelectValues: map[string]string{
				"target,A": "Target A",
				"target,B": "Target B",
				"target,C": "Target C",
				"fclass,A": "F Class A",
				"fclass,B": "F Class B",
				"match,A": "Match A",
				"match,B": "Match B",
			},
		},
		//		"club":Inputs{
		//			Html:      "select",
		//			SelectValues:   getClubSelectBox(eventsCollection),
		//			Label:   "Event Name",
		//		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Add Shooter",
		},
	}
}



func totalScores_update(event_id, shooter_id, range_id string) map[string]Inputs {
	return map[string]Inputs{
		schemaTOTAL:Inputs{
			Html:      "number",
			Label:   "Total",
			Required: true,
			Min: 0,
			Max: 60,
		},
		schemaCENTER:Inputs{
			Html:      "number",
			Label:   "Centers",
			Required: true,
			Min: 0,
			Max: 60,
		},
		"event_id":Inputs{
			Html: "hidden",
			Value: event_id,
			Required: true,
		},
		"shooter_id":Inputs{
			Html: "hidden",
			Value: shooter_id,
			Required: true,
		},
		"range_id":Inputs{
			Html: "hidden",
			Value: range_id,
			Required: true,
		},
		"submit":Inputs{
			Html:      "submit",
			Label:   "Save",
		},
	}
}
