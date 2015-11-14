package main

import (
	"html/template"
	"net/http"
	"strings"
)

const titleSeperator = " - "
const masterTemplate = "../htm/master.htm"

var (
	templates     = make(map[string]*template.Template)
	homeMenuItems = []menu{
		{
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
		},
	}
)

type menu struct {
	Name, Link string
}

type Page struct {
	Title       string
	Data        M
	Menu        []menu
	CurrentYear string
}

func templater(page Page, w http.ResponseWriter) {
	pageName := strings.Split(strings.ToLower(page.Title), titleSeperator)[0]
	//	page.data["Menu"] = homeMenuItems
	page.Menu = homeMenuItems
	page.CurrentYear = currentYear
	stuff, ok := templates[pageName]
	if !ok {
		templates[pageName] = template.Must(template.ParseFiles("../htm/"+pageName+".htm", masterTemplate))
		stuff = templates[pageName]
	}

	//TODO remove below line to dynamically reload templates
	if debug {
		templates[pageName] = template.Must(template.ParseFiles("../htm/"+pageName+".htm", masterTemplate))
	}

	err := stuff.ExecuteTemplate(w, "master", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
