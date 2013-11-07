package main

import(
	"fmt"
)

func main(){
//	router := controller
	router := "hello"
//	path := "controller"
	path := runit()
//	fmt.Print(path["hello"])
	exec, ok := path[router];
	if ok {
		exec()
	}else{
		fmt.Println("key not found")
	}
}

func controller(){
	fmt.Print("\nThis is the controller func")
}


func home(){
	fmt.Print("\nI am home func()")
}









func runit() map[string]func(){
	temp := map[string]func(){
		"t": home,
		"hello": controller,
	}
	return temp
}

func temp(){
	dict := map[string]int {"foo" : 1, "bar" : 2}
	//RESEARCH if using the _ is better than assigning a new variable
	//if _, ok := dict["baz"]; ok { ... }
	value, ok := dict["baz"]
	if ok {
		fmt.Println("value: ", value)
	} else {
		fmt.Println("key not found")
	}
}
