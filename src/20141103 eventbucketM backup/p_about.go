package main

import (
	"net/http"
)

func about(w http.ResponseWriter, r *http.Request) {
	templator(TEMPLATE_HOME, "about", aboutData(), w)
}

func aboutData() map[string]interface{} {
	return map[string]interface{}{
		"Version":  VERSION,
		"PageName": "About",
		"Menu":     home_menu("/about",HOME_MENU_ITEMS),
	}
}
