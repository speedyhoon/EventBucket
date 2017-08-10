package main

import (
	"net/http"
)

func settings(w http.ResponseWriter, r *http.Request) {
	_, forms := sessionForms(w, r, 22)
	templater(w, page{
		Title: "Settings",
		Data: map[string]interface{}{
			"settings": forms[0],
		},
	})
}

func settingsUpdate(f form) (string, error) {
	masterTemplate.IsDarkTheme = f.Fields[0].Value == "Dark"
	return "", nil
}