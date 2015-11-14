package main

import (
	"html/template"
	"net/http"
	"strings"
)

const (
	titleSeparator     = " - "
	masterTemplatePath = "../htm/master.htm"
)

var (
	templates      = make(map[string]*template.Template)
	masterTemplate = Template{
		CSS:         "dirCss",
		CurrentYear: currentYear,
		JS:          "dirJS",
		Menu: []menu{{
			Name: "Home",
			Link: urlHome,
		}, {
			Name: "Archive",
			Link: urlArchive,
		}, {
			Name: "Clubs",
			Link: urlClubs,
		}, {
			Name: "Events",
			Link: urlEvents,
		}, {
			Name: "About",
			Link: urlAbout,
		}, {
			Name: "Shooters",
			Link: urlShooters,
		}, {
			Name: "Licence",
			Link: urlLicence,
		}},
		PNG: "dirPNG",
	}
)

type menu struct {
	Name, Link string
}

type page struct {
	Title string
	Data  M
}

type Template struct {
	JS, CurrentYear, CSS, PNG string
	Page                      page
	Menu                      []menu
}

func templater(w http.ResponseWriter, page page) {
	pageName := strings.Split(strings.ToLower(page.Title), titleSeparator)[0]
	masterTemplate.Page = page
	html, ok := templates[pageName]
	if !ok || debug { //debug is for dynamically reloading templates on every request
		templates[pageName] = template.Must(template.ParseFiles("../htm/"+pageName+".htm", masterTemplatePath))
		html = templates[pageName]
	}

	err := html.ExecuteTemplate(w, "master", masterTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
