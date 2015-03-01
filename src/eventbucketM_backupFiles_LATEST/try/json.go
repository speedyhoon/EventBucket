package main

import (
		"encoding/json"

	"fmt"
	//	"net/http"
	//	"io/ioutil"
//	"io"

	//	"code.google.com/p/go.net/html"
//		"strings"
	//	"log"
//	"strconv"
)



func main(){
//	const jsonStream = `[{"Discipline":{"id":"1","name":"Target Rifle"},"Grade":{"id":"3","name":"C","threshold":"0.00"},"Shooter":{"id":"2796","name":"Stephen Nitschke"},"avg_score":"90.99329","number_of_shoots":"2"}]`
	const jsonStream = `[{"Discipline":{"id":"1"}}]`

	type Discipline struct{
		id string
//		name string
	}
	type N_Grade struct{
		id string
		name string
		threshold string
	}
	type N_Shooter struct{
		id string
		name string
	}
	type NRAA_Grades struct {
		Discipline Discipline
//		Discipline string
//		Grade N_Grade
//		Shooter N_Shooter
//		avg_score string
//		number_of_shoots string
	}


	test := NRAA_Grades{
		Discipline: Discipline{
			id: "hi!",
		},
	}



//	type Message struct {
//		Name, Text string
//	}
//	dec := json.NewDecoder(strings.NewReader(jsonStream))

//	for {
//		x := make([]NRAA_Grades, 0)
//		if err := dec.Decode(&x); err == io.EOF {
//			break
//		} else if err != nil {
//			fmt.Printf("%#v",err)
//		}
//		fmt.Printf("%s: %s\n", x.Name, x.Text)
//		fmt.Printf("%+v \n", x)
//	}

//	str := `{"page": 1, "fruits": ["apple", "peach"]}`
//	res := []NRAA_Grades{}
//	json.Unmarshal([]byte(jsonStream), &res)


	strB, _ := json.Marshal(test)
//	fmt.Printf(string(strB))

	fmt.Printf("%v \n", string(strB))
	fmt.Printf("\n\n\n\n")

	var mike NRAA_Grades

	tree := json.Unmarshal(strB, &mike)
	fmt.Printf("%+v \n", tree)
	fmt.Printf("%+v \n", mike.Discipline.id)
}
