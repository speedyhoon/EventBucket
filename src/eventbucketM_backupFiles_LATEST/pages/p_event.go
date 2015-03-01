package main

import (
	"net/http"
	"fmt"
)

func event(w http.ResponseWriter, r *http.Request) {
	page_url := "/event/"
	//	var validID = regexp.MustCompile(`\A`+page_url+`[0-9a-f]{24}\z`)
	url := fmt.Sprintf("%v", r.URL)

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
	var this_event map[string]interface{}
	for _, row := range event{
//		if obId(row["_id"]) == event_id{
		if row["_id"] == event_id{
			this_event = row
		}
	}

	fmt.Println(event_id)
	fmt.Println(event)
	fmt.Printf("%v", event)
	return map[string]interface{}{
		"Title": this_event[schemaNAME],
		"Id": event_id,
		"ListRanges": this_event[schemaRANGE],
	}
}

func event_HTML() string {
	return loadHTM("event")
}
