package main

import (
	"net/http"
)

func settings(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Settings",
		Data: map[string]interface{}{
			"Port":  portAddr, //TODO during save, redirect after post request to the new port address
			"Theme": masterTemplate.Theme,
		},
		//TODO Add form 2 for shutting down the HTTP server
	})
}
