package main

import "html/template"
import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("header.html", "organisers.htm")
    t.Execute(w, map[string] string {"Title": "My title", "Body": "Hi this is my body"})
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}