package main

import (
	"net/http"
	"io"
	"websocket"
)

func EchoServer(ws *websocket.Conn){
	io.Copy(ws, ws)
}

func main(){
	http.Handle("/", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
