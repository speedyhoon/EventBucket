package main

import (
	"net/http"
	"net/url"
	"fmt"
	"strings"
	"go-randomdata-master"
	"math/rand"
	"time"
)

func random_data(w http.ResponseWriter, r *http.Request) {
	event_id := "W"
	range_id := "1"
	random_grades := []string{/*"a","b","c",*/"d","e","f","g","h","i","j"}

	attributes := strings.Split(strings.Replace(r.RequestURI, "/random-data/", "", -1), "&")



	for _, request := range attributes {
		properties := strings.Split(request, "=")

		switch properties[0] {
		case "shooterQty":
			shooterQty := str_to_int(properties[1])
			counter := 0
			for counter <= shooterQty {
				counter += 1
				//make some requests for x number of shooters
//				resp, _ := http.PostForm("http://localhost/shooterInsert",
//				go http.PostForm("http://localhost/shooterInsert",
				http.PostForm("http://localhost/shooterInsert",
					url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
					"surname":   {randomdata.LastName()},
					"club":      {randomdata.State(randomdata.Large)},
					"age":       {"N"},
					"grade":     {random_grades[rand.Intn(len(random_grades)-1)]},
					"event_id":  {event_id},
				})
//				defer resp.Body.Close()
			}
		case "totalScores":
			event, _ := getEvent(event_id)
			for shooter_id, _ := range event.Shooters{
				rand.Seed(time.Now().UnixNano())
				//				rand.Seed(90)
				fmt.Println(shooter_id)
//				resp, _ := http.PostForm("http://localhost/updateTotalScores",
//				go http.PostForm("http://localhost/updateTotalScores",
//					url.Values{"first":      {randomdata.FirstName(randomdata.RandomGender)},
//					"score":			{fmt.Sprintf("%v.%v",rand.Intn(51),rand.Intn(11))},
//					"shooter_id":	{shooter_id},
//					"range_id":		{range_id},
//					"event_id":		{event_id},
//				})
//				defer resp.Body.Close()
				event_total_score_update(event_id, range_id, shooter_id, Score{
						Total: rand.Intn(51),
						Centers: rand.Intn(11),
				})
			}
		}
	}
/*
	_, ok = form["start_shooting"]

	if ok {
		//fill out every shooters range
	}*/

}
