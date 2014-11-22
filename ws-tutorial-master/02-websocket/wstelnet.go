package main

import (
	"fmt"
	"net/http"
	"os"
//	"code.google.com/p/go.net/websocket"
	"websocket"
	//"bufio"
	//"net"
    "math/rand"
	"time"
)

func main() {
	http.Handle("/script/", http.FileServer(http.Dir(".")))
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.Handle("/websocket/", websocket.Handler(ProcessSocket))
	err := http.ListenAndServe(":4000", nil)
	checkError(err)
}

//func EchoServer(ws *websocket.Conn) {
//    var msg string
//    websocket.Message.Receive(ws, &msg)
//    fmt.Println("Message Got: ", msg)
//}

//When the handler finishes the websocket connection is closed.
//If you'd like to keep the socket open you have to keep the handler running.
//eg.
//func EchoServer(ws *websocket.Conn) {
//    for {
//        var msg string
//        err := websocket.Message.Receive(ws, &msg)
//        if err != nil{
//            break
//        }
//        fmt.Println("Message Got: ", msg)
//    }
//}

func ProcessSocket(ws *websocket.Conn) {
	fmt.Println("In ProcessSocket")
	
	for {
		var msg string

		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Println("ProcessSocket: got error", err)
			_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
			return
		}
		fmt.Println("ProcessSocket: got message", msg)

		//service := msg
		//
		//tcpAddr, err := net.ResolveTCPAddr("tcp", service)
		//if err != nil {
		//	fmt.Println("Error in ResolveTCPAddr:", err)
		//	_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
		//	return
		//}
		//
		//conn, err := net.DialTCP("tcp", nil, tcpAddr)
		//if err != nil {
		//	fmt.Println("Error in DialTCP:", err)
		//	_ = websocket.Message.Send(ws, "FAIL:" + err.Error())
		//	return
		//}
		
		websocket.Message.Send(ws, "SUCC")
		
		go simulateEvent(ws, "High jump")
		
	}

	//RunTelnet(ws, conn)
}

//func simulateEvent(ws *websocket.Conn, name string, timeInSecs int) { 
func simulateEvent(ws *websocket.Conn, name string) { 

	//speed := random(0, 301)
	speed := fmt.Sprintf("%v", rand.Intn(250) )

    // sleep for a while to simulate time consumed by event
    fmt.Println("speed:", speed)
	//websocket.Message.Send(ws, "fdsa")
	websocket.Message.Send(ws, speed)
    //time.Sleep(timeInSecs * 1e9 )
	
	//time.Sleep(time.Duration(1) * time.Second)
	//time.Sleep(0.5 * time.Second)
	//amt := time.Duration(rand.Intn(250))
	amt := time.Duration(40)
	time.Sleep(time.Millisecond * amt)
    
	//fmt.Println("Finished ", name)
	
	simulateEvent(ws, "another")
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

