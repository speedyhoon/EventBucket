package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
	"html/template"
    "websocket"
)

const listenAddr = "localhost:4000"

func main() {
    http.HandleFunc("/", rootHandler)
    http.Handle("/socket", websocket.Handler(socketHandler))
    err := http.ListenAndServe(listenAddr, nil)
    if err != nil {
        log.Fatal(err)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var rootTemplate = template.Must(template.New("root").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Simple Websockets</title>
	</head>
	<body>
	<script>
	websocket = new WebSocket("ws://{{.}}/socket");
	websocket.onmessage = onMessage;
	websocket.onclose = onClose;
	</script>
	</body>
	</html>
	`))
    rootTemplate.Execute(w, listenAddr)
}

type socket struct {
    io.ReadWriter
    done chan bool
}

func (s socket) Read(b []byte) (int, error)  { return s.conn.Read(b) }
func (s socket) Write(b []byte) (int, error) { return s.conn.Write(b) }

func (s socket) Close() error {
    s.done <- true
    return nil
}

func socketHandler(ws *websocket.Conn) {
    s := socket{ws, make(chan bool)}
    go match(s)
    <-s.done
}