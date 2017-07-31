package main

import (
	"net/http"
)

func settings(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Settings",
		Data: map[string]interface{}{
			"port":  portAddr, //TODO during save, post redirect to new port address
			"theme": masterTemplate.IsDarkTheme,
		},
		//TODO Form 2: shutdown http server
	})
}
