package main

func licence() Page {
	return Page{
		TemplateFile: "licence",
		Title:        "Licence",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Menu": homeMenu(URL_licence, HOME_MENU_ITEMS),
		},
	}
}
