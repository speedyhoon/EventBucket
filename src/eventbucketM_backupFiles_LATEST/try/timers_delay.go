package main

import (
	"time"
	"net/http"
	"fmt"
)

//var itemsUpdate = "Boo"

func main() {
	itemsUpdate := "Hi there"
	timer := time.NewTimer(time.Second * 5)
	go func() {
		<- timer.C
		println("Timer expired")
		fmt.Printf(itemsUpdate)
	}()
//	stop := timer.Stop()
//	println("Timer cancelled:", stop)

	itemsUpdate = "Boo Radlyet"
	time.AfterFunc(5 * time.Second, func() {
		println("expired after 5 seconds")
		fmt.Printf(itemsUpdate)
	})

	http.HandleFunc("/", runnit)

	err := http.ListenAndServe(":801", nil)
	if err != nil{
		fmt.Printf("ListenAndServe: %v", err)
	}
}

func runnit(w http.ResponseWriter, r *http.Request){
	fmt.Printf("boo :)")
}
