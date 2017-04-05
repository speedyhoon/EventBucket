package main

import (
	"net/http"
)

func settings(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Settings",

		Data: map[string]interface{}{
			//Form 1: update port and/or theme
			"Port":  portAddr + "--BuildDate GoVersion f fdsasds5", //during save, change post redirect to new port address
			"Theme": masterTemplate.Theme,
		},
		//TODO Form 2: shutdown http server
	})
}
