package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"encoding/json"
)

const STRLEN = 80
func main(){
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/send", send)
	http.HandleFunc("/helloThis", helloThis)
	http.HandleFunc("/img", random_image)
	http.ListenAndServe(":85", nil)
}


type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func cars(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"form": "cars",
		"make": "Ford",
		"model": "Escort",
		"wheels": 4,
	})
}

func send(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"form": "sender",
		"first": "Bill",
		"last": "Bates",
		"email": "bill@bates.com",
		"comments": "blah blah blah",
		"date": "",
	})
}

func helloThis(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, Response{
		"form": "sender",
		"first": "Bill",
		"last": "Bates",
		"email": "bill@bates.com",
		"comments": "blah blah blah",
		"date": "",
	})
}


func random_image(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	b := make([]byte, STRLEN)
	rand.Read(b)
	en := base64.StdEncoding // or URLEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	fmt.Fprint(w, Response{"img": fmt.Sprintf("src=%s\ndst=%s\n", b, d)})
}
