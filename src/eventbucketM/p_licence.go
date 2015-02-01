package main

func licence()Page{
	return Page {
		TemplateFile: "licence",
		Title: "Licence",
		Theme: TEMPLATE_HOME,
		Data: M{
			"Menu":     home_menu(URL_licence, HOME_MENU_ITEMS),
		},
	}
}

//func licence_summary()M{
//	return M{
//		"Menu":     home_menu("/licence-summary",HOME_MENU_ITEMS),
//	}
//}
