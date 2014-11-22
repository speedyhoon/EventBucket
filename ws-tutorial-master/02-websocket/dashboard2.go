package main

import (
	"fmt"
	"net/http"
	"os"
	"websocket"
    "math/rand"
	"time"
)

func main() {
	//http.Handle("/script/", http.FileServer(http.Dir(".")))
	//http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.Handle("/websocket/", websocket.Handler(ProcessSocket))
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
var on = false
var hibeam = false
var left = false
var right = false

func ProcessSocket(ws *websocket.Conn) {
	//fmt.Println("In ProcessSocket")
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Println("ProcessSocket: got error", err)
			_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
			return
		}
		//fmt.Println("ProcessSocket: got message", msg)
		toggle(msg)
		
		if on == false{
			on = true
			go simulateVarry(ws, "speed", 149, 40, time.Millisecond, 0)
			go simulateVarry(ws, "revs", 180, 40, time.Millisecond, 0)
			go simulateVarry(ws, "odom", 100, 1, time.Second, 0)
			go simulateVarry(ws, "fuel", 100, 1, time.Second, 0)
			go simulateVarry(ws, "mpg", 30, 3, time.Second, 0)
			go simulateVarry(ws, "lp100", 10, 3, time.Second, 0)
			//go simulateVarry(ws, "hibeam", 2, 1, time.Second, 0)
			//go simulateVarry(ws, "left", 2, 1, time.Second, 0)
			//go simulateVarry(ws, "right", 2, 1, time.Second, 0)
			go simulateVarry(ws, "oil", 2, 1, time.Second, 0)
			go simulateVarry(ws, "handbrake", 2, 1, time.Second, 0)
			go simulateVarry(ws, "water", 2, 1, time.Second, 0)
		}
	}
}

func simulateVarry(ws *websocket.Conn, attribute string, amount int, delay int, units time.Duration, offset int){
	value := fmt.Sprintf("{\"%v\":%v}", attribute, (rand.Intn(amount) + offset) )
	websocket.Message.Send(ws, value)
	time.Sleep(units * time.Duration(delay))
	simulateVarry(ws, attribute, amount, delay, units, offset)
}

func toggle(message string){
	if message == "hibeam" {
		if hibeam{
			hibeam = false
		}else{
			hibeam = true
		}
		fmt.Println(message, hibeam)
	}
	if message == "left" {
		if left{
			left = false
		}else{
			left = true
		}
		fmt.Println(message, left)
	}
	if message == "right" {
		if right{
			right = false
		}else{
			right = true
		}
		fmt.Println(message, right)
	}
}

//1 litre per 100km = 282.48105314960625 mpg