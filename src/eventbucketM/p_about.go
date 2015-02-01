package main

func about()Page {
	return Page {
		Name: "About",
		Theme: TEMPLATE_HOME,
		Data: M{
			"Version":  VERSION,
			"PageName": "About",
			"Menu":     home_menu(URL_about, HOME_MENU_ITEMS),
		},
	}
}
