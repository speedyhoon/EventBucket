package main

func about()M{
	return M{
		"Version":  VERSION,
		"PageName": "About",
		"Menu":     home_menu(URL_about ,HOME_MENU_ITEMS),
	}
}
