package main

func championships() Page {
	return Page{
		TemplateFile: "championships",
		Title:        "Championships",
		Theme:        TEMPLATE_HOME,
		Data: M{
			"Championship": generateForm(championshipForm()),
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
				Inner: "Add Championship",
				//AccessKey: "x",
			},
		},
	}
}
