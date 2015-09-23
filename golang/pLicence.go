package main

func licence() Page {
	return Page{
		TemplateFile: "licence",
		Title:        "Licence",
		Theme:        templateHome,
		Data:         M{},
	}
}
