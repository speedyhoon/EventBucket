package main

import (
	"net/http"
)

func licence(w http.ResponseWriter, r *http.Request) {
	test := map[string]interface{}{
		"Menu":     home_menu("/licence",HOME_MENU_ITEMS),
	}
	templator(TEMPLATE_HOME, "licence", test, w)
}

func licence_summary(w http.ResponseWriter, r *http.Request) {
	test := map[string]interface{}{
		"Menu":     home_menu("/licence-summary",HOME_MENU_ITEMS),
	}
	templator(TEMPLATE_HOME, "licence-summary", test, w)
}
