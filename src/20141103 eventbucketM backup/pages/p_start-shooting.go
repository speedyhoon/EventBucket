package main

import (
	"net/http"
	"fmt"
	"strings"
//	"strconv"
//	"reflect"
)
func startShooting(w http.ResponseWriter, r *http.Request) {
	page_url := "/startShooting/"
	url := fmt.Sprintf("%v", r.URL)

	slice := strings.Split(url[len(page_url):], "/")
//	dump(reflect.TypeOf(slice))
//	dump(len(slice))
	if len(slice) == 2{
//		dump("slice type is correct")
		templator("admin", startShooting_HTML(), startShooting_Data(slice[0], slice[1]), w)
	}else{
		dump("startShooting event_id and range_id wasn't found")
	}
}
func startShooting_Data(event_id, range_id string) map[string]interface{} {
//	arr := strings.Split(url, "/")
//	dump(url)
//	var array string = arr[1]

//	index, err := strconv.Atoi(array)
//	checkErr(err)

//	index++

//	fmt.Println("\n\n====")
//	fmt.Printf(  fmt.Sprintf("%v", arr[1]  ))
//	fmt.Println("\n\n")

	event := getDoc_by_id("event", event_id)

	currentRange := range_id // event[schemaRANGE].([]map[string]interface{})[index]

	dump(event[schemaRANGE])

//	if ok {
//		dump("current Range was found ok!!")
//	}

//	if event[schemaRANGE].(map[string]interface{})[arr[1]] {
//		currentRange = event[schemaRANGE].(map[string]interface{})[arr[1]]
//	}


	classShots := map[string][]string{
		"fclass": []string{
			"S1","S2","S3","1","2","3","4","5","6","7","8","9","10","11","12","13","14","15",
		},
		"target": []string{
			"S1","S2","1","2","3","4","5","6","7","8","9","10",
		},
		"match": []string{
			"S1","S2","1","2","3","4","5","6","7","8","9","10","11","12","13","14","15","16","17","18","19","20",
		},
	}

//	fmt.Println("\n\n=fds===")
//	fmt.Printf(   fmt.Sprintf("%v", event[schemaRANGE].(map[string]interface{})[arr[1]].(map[string]interface{})[schemaNAME]   ))
//	fmt.Println("\n\n")


	output := map[string]interface{}{
		"EventId": event_id,
//		"Title": fmt.Sprintf("Start Shooting: %v", event[schemaRANGE].(map[string]interface{})["1"].(map[string]interface{})[schemaNAME]),
		"RangeName": currentRange,//  event[schemaRANGE].(map[string]interface{})[arr[1]].(map[string]interface{})[schemaNAME],

		"classShots": classShots,

//		"Id": event_id,
//		"AddRange": generateForm("rangeInsert", eventSettings_add_rangesForm(event_id)),
//		"RangeList": event[schemaRANGE],
//		"RangeList":/ event_obj.Ranges,
		//		"RangeList": ranges,
//		"AddShooter": generateForm("shooterInsert", eventSettings_add_shooterForm(event_id)),
		"ShooterList": event[schemaSHOOTER],
	}
	if event_ranges, ok := event[schemaRANGE]; ok {
		output["RangeList"] = event_ranges
	}
	return output
}
func startShooting_HTML() string {
	return loadHTM("start-shooting")
}
