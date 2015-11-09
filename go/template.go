package main

import (
	"html/template"
	"net/http"
	"strings"
)

const titleSeperator = " - "
const masterTemplate = "../htm/master.htm"

var templates = make(map[string]*template.Template)

type Page struct {
	title string
	data  M
}

func templater(page Page, w http.ResponseWriter) {
	pageName := strings.Split(strings.ToLower(page.title), titleSeperator)[0]

	stuff, ok := templates[pageName]
	if !ok {
		templates[pageName] = template.Must(template.ParseFiles("../htm/"+pageName+".htm", masterTemplate))
		stuff = templates[pageName]
	}

	//TODO remove below line to dynamically reload templates
	if debug {
		templates[pageName] = template.Must(template.ParseFiles("../htm/"+pageName+".htm", masterTemplate))
	}

	err := stuff.ExecuteTemplate(w, "master", page.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
