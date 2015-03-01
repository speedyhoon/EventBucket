package main

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	templator("home", loadHTM("home"), homeData(getCollection("event")), w)
}
type HomeCalendar struct {
	Id, Name, Club, ClubId, Day, Date, Month, Time string
}
func homeData(event []map[string]interface{}) map[string]interface{} {
	all := []HomeCalendar{}
	for _, row := range event {
		datetime := exists(row, "datatime")
		all = append(all, HomeCalendar{
//			Name: exists(row, "name"),
			Id: exists(row, schemaID),
			Name: exists(row, schemaNAME),
			Club: "Club Name",
			ClubId: exists(row, schemaCLUB),
//			Club: exists(row, "clubId"),
			Day: datetime,
			Date: datetime,
			Month: datetime,
			Time: datetime,
		})
	}
	return map[string]interface{}{
		"Events": all,
		"PageName": "Calendar",
		"Menu": map[string]string{
			"All Events": "",
			"Features": "features",
			"Clubs": "clubs",
			"Organisers": "organisers",
			"Archive": "archive",
		},
	}

}
