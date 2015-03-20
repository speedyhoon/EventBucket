package main

func licence() Page {
	return Page{
		TemplateFile: "licence",
		Title:        "Licence",
		Theme:        TEMPLATE_HOME,
		Data:         M{},
	}
}
