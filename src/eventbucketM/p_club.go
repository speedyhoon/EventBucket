package main

func club(clubId string) Page {
	club, err := getClub(clubId)
	if err != nil{
		//TODO return a 404 error
		return Page {
			TemplateFile: "club",
			Theme: TEMPLATE_HOME,
			Title:  "Club with id '" + clubId + "' not found",
			Data: M{
				"Menu":  home_menu(URL_club, HOME_MENU_ITEMS),
			},
		}
	}
	return Page {
		TemplateFile: "club",
		Theme: TEMPLATE_HOME,
		Title: club.Name,
		Data: M{
			"Menu": home_menu(URL_club, HOME_MENU_ITEMS),
			"ClubId": clubId,
		},
	}
}
