package main

import (
	"net/http"
	"fmt"
	"strings"
//	"strconv"
)
func scoreboard(w http.ResponseWriter, r *http.Request) {
	page_url := "/scoreboard/"
	url := fmt.Sprintf("%v", r.URL)
	templator("admin", scoreboard_HTML(), scoreboard_Data(url[len(page_url):]), w)
}
func scoreboard_Data(url string) map[string]interface{} {
	arr := strings.Split(url, "/")

	event_id := arr[0]

//	dump("gggg")
//	dump(url)
//	dump(arr)

//	var array string = arr[1]
//
//	index, err := strconv.Atoi(array)
//	checkErr(err)
//
//	index++

//	fmt.Println("\n\n====")
//	fmt.Printf(  fmt.Sprintf("%v", arr[1]  ))
//	fmt.Println("\n\n")

	event := getDoc_by_id("event", event_id)

	dump(event)

//	currentRange :=  event[schemaRANGE].(map[string]interface{})[range_id]

//	dump(event[schemaRANGE])

//	if ok {
//		dump("current Range was found ok!!")
//	}

//	if event[schemaRANGE].(map[string]interface{})[arr[1]] {
//		currentRange = event[schemaRANGE].(map[string]interface{})[arr[1]]
//	}


//	classShots := map[string][]string{
//		"fclass": []string{
//			"S1","S2","S3","1","2","3","4","5","6","7","8","9","10","11","12","13","14","15",
//		},
//		"target": []string{
//			"S1","S2","1","2","3","4","5","6","7","8","9","10",
//		},
//		"match": []string{
//			"S1","S2","1","2","3","4","5","6","7","8","9","10","11","12","13","14","15","16","17","18","19","20",
//		},
//	}

//	fmt.Println("\n\n=fds===")
//	fmt.Printf(   fmt.Sprintf("%v", event[schemaRANGE].(map[string]interface{})[arr[1]].(map[string]interface{})[schemaNAME]   ))
//	fmt.Println("\n\n")


	output := map[string]interface{}{
		"EventId": arr[0],
		"EventName": event[schemaNAME],
//		"Title": fmt.Sprintf("Start Shooting: %v", event[schemaRANGE].(map[string]interface{})["1"].(map[string]interface{})[schemaNAME]),
//		"RangeName": currentRange.(map[string]interface{})[schemaNAME],//  event[schemaRANGE].(map[string]interface{})[arr[1]].(map[string]interface{})[schemaNAME],
//		"RangeId": range_id,
//		"classShots": classShots,

//		"Id": event_id,
//		"AddRange": generateForm("rangeInsert", eventSettings_add_rangesForm(event_id)),
//		"RangeList": event[schemaRANGE],
//		"RangeList":/ event_obj.Ranges,
		//		"RangeList": ranges,
//		"AddShooter": generateForm("shooterInsert", eventSettings_add_shooterForm(event_id)),
		"ShooterList": event[schemaSHOOTER],
		"RangeList": event[schemaRANGE],
		"Css": "scoreboard",
	}
	if event_ranges, ok := event[schemaRANGE]; ok {
		output["RangeList"] = event_ranges
	}
	return output
}
func scoreboard_HTML() string {
	return loadHTM("scoreboard")
}
