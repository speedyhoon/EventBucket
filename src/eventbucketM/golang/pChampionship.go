package main

func championships() Page {
	return Page{
		TemplateFile: "championships",
		Title:        "Championships",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Championship": generateForm(championshipForm()),
			"Menu":         homeMenu(URL_championships, HOME_MENU_ITEMS),
		},
	}
}

func championshipForm() Form {
	return Form{
		Action: URL_champInsert,
		Title:  "Create Championship",
		Inputs: []Inputs{
			{
				Name:     "name",
				Html:     "text",
				Label:    "Championship Name",
				Required: true,
			},
			{
				Html:  "submit",
				Value: "Add Championship",
			},
		},
	}
}
