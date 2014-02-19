package main

import "fmt"


func main(){
	fmt.Print("start\n")
	getmap()
	fmt.Print("end\n")
}
func getmap()map[int]interface{}{
	return map[int]interface{}{
		0: "event name",
		1: "club_id",
		2: "datetime",

		3: "description",
		4: "shooters",
		5: "settings",
		6: map[int]interface{}{
			0: 1234,
			1: "class",
			2: "grade",
			//scores
			3: map[int]interface{}{

			},
		},
		//teams
		7: "",
		8: "",
		9: "",
		10:"",
	}
}
