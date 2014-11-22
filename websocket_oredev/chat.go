package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"html/template"
	"websocket"
)

const listenAddr = "localhost:80"
func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)} // HL
	go match(s)
	<-s.done
}

type socket struct {
	io.ReadWriter // HL
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}


var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<script>
var input, output, websocket;

function showMessage(m) {
	var p = document.createElement("p");
	p.innerHTML = m;
	output.appendChild(p);
}

function onKey(e) {
	if (e.keyCode == 13) {
		var m = input.value;
		input.value = "";
		websocket.send(m);
		showMessage(m);
	}
}

function init() {
	input = document.getElementById("input");
	input.addEventListener("keyup", onKey, false);

	output = document.getElementById("output");

	websocket = new WebSocket("ws://{{.}}/socket");
	websocket.onmessage = function onMessage(e) { showMessage(e.data); }
	websocket.onclose = showMessage("Connection closed.");;
}

window.addEventListener("load", init, false);
</script>
</head>
<body>
<input id="input" type="text">
<div id="output"></div>
</body>
</html>
`))