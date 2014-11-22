package main

import (
	"net/http"
	"io"
	"websocket"
)

// This is the guy that handle the WebSocket comunication.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws);
}

func main() {
	// Map the url "/echo" the your websocket EchoServer function
	http.Handle("/echo", websocket.Handler(EchoServer));
	// Start the httpServer from http package
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic(err)
	}
}