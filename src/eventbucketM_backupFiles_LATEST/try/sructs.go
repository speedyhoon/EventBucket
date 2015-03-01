package main

import (
	//	"encoding/json"

	"fmt"
	//	"net/http"
	//	"io/ioutil"

//	"code.google.com/p/go.net/html"
//	"strings"
	//	"log"
	"strconv"
)



func main(){
//	shooter_type := Shooter{
//		SID: 11277,
//		Surname: "Penny",
//		First_name: "Damien Ian",
//		Nickname: "Damien",
//		Club: "Canberra Rifle Club",
//	}

//	for index, value := range shooter_type{
//		fmt.Printf("%v - %v", index, value)
//	}

	this_row := []string{
		"15",
		"11277",
		"Penny",
		"Damien Ian",
		"Damien",
		"Canberra Rifle Club",
	}
	shooter_list := Shooter{}

	for i := 1; i <= 5; i++ {
		switch {
		case i == 1:
			shooter_list.SID, _ = strconv.Atoi(this_row[i])
			break
		case i == 2:
			shooter_list.Surname = this_row[i]
			break
		case i == 3:
			shooter_list.First_name = this_row[i]
			break
		case i == 4:
			shooter_list.Nickname = this_row[i]
			break
		case i == 5:
			shooter_list.Club = this_row[i]
			break
		}
	}
	fmt.Printf("%#v", shooter_list)
}
