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
		go simulateVarry(ws, "speed", 149, 40, time.Millisecond)
		go simulateVarry(ws, "revs", 400, 40, time.Millisecond)
		go simulateVarry(ws, "odom", 100, 1, time.Second)
		go simulateVarry(ws, "fuel", 100, 1, time.Second)
		go simulateVarry(ws, "mpg", 30, 3, time.Second)
		go simulateVarry(ws, "hibeam", 2, 1, time.Second)
		go simulateVarry(ws, "left", 2, 1, time.Second)
		go simulateVarry(ws, "right", 2, 1, time.Second)
		go simulateVarry(ws, "oil", 2, 1, time.Second)
		go simulateVarry(ws, "handbrake", 2, 1, time.Second)
		go simulateVarry(ws, "water", 2, 1, time.Second)
	}
}

func simulateVarry(ws *websocket.Conn, attribute string, amount int, delay int, units time.Duration){
	value := fmt.Sprintf("{\"%v\":%v}", attribute, rand.Intn(amount) )
	websocket.Message.Send(ws, value)
	time.Sleep(units * time.Duration(delay))
	simulateVarry(ws, attribute, amount, delay, units)
}

//func simulateSpeed(ws *websocket.Conn){
//	speed := fmt.Sprintf("{\"speed\":%v}", rand.Intn(149) )
//	websocket.Message.Send(ws, speed)
//	amt := time.Duration(40)
//	time.Sleep(time.Millisecond * amt)
//	simulateSpeed(ws)
//}
//
//func simulateRevs(ws *websocket.Conn){
//	speed := fmt.Sprintf("{\"revs\":%v}", rand.Intn(7000) )
//	websocket.Message.Send(ws, speed)
//	amt := time.Duration(40)
//	time.Sleep(time.Millisecond * amt)
//	simulateRevs(ws)
//}

//func simulateFuel(ws *websocket.Conn){
//	fuel := fmt.Sprintf("{\"fuel\":%v}", rand.Intn(50) )
//	websocket.Message.Send(ws, fuel)
//	amt := time.Duration(rand.Intn(5))
//	time.Sleep(time.Second * amt)
//	simulateFuel(ws)
//}

//func simulateAny(ws *websocket.Conn, attribute string){
//	value := fmt.Sprintf("{\"%v\":%v}", attribute, rand.Intn(2) )
//	//fmt.Println(attribute, " is: ", value)
//	websocket.Message.Send(ws, value)
//	amt := time.Duration(1)
//	time.Sleep(time.Second * amt)
//	simulateAny(ws, attribute)
//}