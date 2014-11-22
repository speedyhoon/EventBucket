package main

import (
	"encoding/json"
	"fmt"
)

const (
	hello = `json:"S,omitempty"`
)


type AutoInc struct{
	Mound int									`json:"M,omitempty"`
	Event int									`json:"E,omitempty"`
	Club int										`json:"C,omitempty"`
	Range int									`json:"R,omitempty"`
	Shooter int									`json:"S,omitempty"`
	Scores map[string]Score					`bson:"omitempty,inline"`
}

type Score struct{

}

func main(){
	temp := AutoInc{
		Mound: 1,
		Event: 2,
		Club: 3,
		Range: 4,
		Shooter: 5,
	}

	b, _ := json.Marshal(temp)
	c, _ := json.Marshal(AutoInc{Mound:2})

	fmt.Printf("%v\n", string(b))
	fmt.Printf("%v\n", string(c))




	fmt.Printf(hello)
}
