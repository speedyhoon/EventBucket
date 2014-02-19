package main

import (
    "websocket"
    "fmt"
    "log"
    "time"
    )

func main() {
    origin := "http://localhost:8080/"
    url := "ws://localhost:8080/ws"

    var err error
    var ws *websocket.Conn
    for {
        ws, err = websocket.Dial(url, "", origin)
        if err != nil {
            fmt.Println("Connection fails, is being re-connection")
            time.Sleep(1*time.Second)
            continue
        }
        break
    }
    if _, err := ws.Write([]byte("something")); err != nil {
        log.Fatal(err)
    }

}