package main

import (
	"net/http"
	"fmt"
	"html/template"
)

const table_style = ".table{display:table}.tr{display:table-row}.td,.th{display:table-cell}.th{font-weight:900;text-align:center}"

func generator(w http.ResponseWriter, fillin string, data map[string]interface{}) {
	my_html := template.New("fillin").Funcs(template.FuncMap {
	"XTC": func(x string) template.HTML {
		return template.HTML(x)
	}})
	t := template.Must(my_html.Parse(fillin))
	err := t.Execute(w, data)
	if err != nil {
		fmt.Print("executing template:", err)
	}
}

func templator(template string, fill_in string, data map[string]interface{}, w http.ResponseWriter) {
	switch (template){
	case "home":
		generator(w, templateHome(fill_in), data)
		break
	case "admin":
		generator(w, templateAdmin(fill_in), data)
		break
	case "ajax":
		generator(w, fill_in, data)
	}
}
func templateHome(body string) string {
	return `<!doctype html>
<html>
<head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
<style>body{background:#000;color:#777}` + table_style + `</style>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
}
func templateAdmin(body string) string {
	return `<!doctype html><html><head>
	{{if .Css}}<link rel=stylesheet href={{.Css}}>{{end}}
	{{if .Ico}}<link rel=icon type={{.IcoType}} href={{.Ico}}>{{end}}
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
	<style>body{background:#000}*{color:#777}` + table_style + `</style>
</head><body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body></html>`
}
func panel(title, source string) string {
	return `<h2>` + title + `</h2><div>` + source + `</div>`
}
func pane(title, body string) string {
	return `<div><h4>` + title + `</h4>` + body + `</div>`
}
func generateTableForm(formData map[string]InputTypers) string {
	form := "<form"
	if form_data, ok := formData["form"]; ok {
		form += fmt.Sprintf(" method=%v action=%v", addQuotes(form_data.method), addQuotes(form_data.action))
		if form_data.table == true {
			form += " class=tr"
		}
	}else {
		panic(fmt.Sprintf("A form element is not set in the form object: %v", formData))
	}
	display_as_table := formData["form"].table

	//	for attribute, value := range form_attributes{
	//		if attribute != "table"{
	//			form += fmt.Sprintf(" %v=%v",attribute, addQuotes(value))
	//		}
	//	}
	form += ">"
	//	output := "<form method=post action="+addQuotes(action)+">"
	for inputName, inputData := range formData {
		output := "\n"

		if inputData.Html != "submit" && inputData.Label != "" {
			output += " "+inputData.Label+":"
		}
		if inputData.Html == "text" || inputData.Html == "submit" || inputData.Html == "number" || inputData.Html == "url" || inputData.Html == "datetime" || inputData.Html == "checkbox" || inputData.Html == "radio" {
			output += "<input"
			if inputData.Html == "submit" || inputData.Html == "number" || inputData.Html == "url" || inputData.Html == "datetime" || inputData.Html == "checkbox" || inputData.Html == "radio" {
				output += " type="+inputData.Html
			}
			if inputData.Html != "submit" {
				output += " name="+inputName
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Disabled {
				output += " disabled"
			}
			if inputData.Checked {
				output += " checked"
			}
			//			if inputData.AutoCorrect == false {
			//				output += " autocorrect=off"
			//			}
			//			if inputData.AutoCorrect == true {
			//				output += " autocorrect=on"
			//			}
			//			if inputData.AutoCapitalize == false {
			//				output += " autocapitalize=off"
			//			}
			//			if inputData.AutoCapitalize == true {
			//				output += " autocapitalize=on"
			//			}

			if inputData.PlaceHolder != "" {
				output += " placeHolder="+addQuotes(inputData.PlaceHolder)
			}
			if inputData.RangeMin > 0 {
				output += fmt.Sprintf(" min=%d", inputData.Min)
			}
			if inputData.RangeMax > 0 {
				output += fmt.Sprintf(" max=%d", inputData.Max)
			}
			//			if  inputData.Html == "submit"{
			//				output += fmt.Sprintf(" value=%v",addQuotes(inputData.Label))
			//			}
			if inputData.Value != "" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Value))
			}
			output += ">"
		}else if inputData.Html == "select" {
			output += "<select name="+inputName
			if inputData.Required {
				output += " required"
			}
			if inputData.MultiSelect {
				output += " multiple"
			}
			if inputData.PlaceHolder != "" {
				output += " placeHolder="+addQuotes(inputData.PlaceHolder)
			}
			if inputData.AutoCorrect == false {
				output += " autocorrect=off"
			}
			if inputData.AutoCorrect == true {
				output += " autocorrect=on"
			}
			output += ">"
			for _, option := range inputData.Select {
				output += fmt.Sprintf("\n\t<option>%v</option>", option)
			}
			for value, option := range inputData.SelectValues {
				output += fmt.Sprintf("\n\t<option value=%v>%v</option>", addQuotes(value), option)
			}
			output += "</select>"
		}
		if inputData.Help != "" {
			output += fmt.Sprintf("<abbr title=%v>?</abbr>", addQuotes(inputData.Help))
		}

		if display_as_table {
			form += fmt.Sprintf("<span class=td>%v</span>", output)
		}else {
			form += output
		}
	}
	form += "</form>"
	return form
}

func generateForm(action string, formData map[string]Inputs) string {
	output := fmt.Sprintf("<form action=%v method=post>", addQuotes(action))
	//	output := "<form action="+addQuotes(action)+" method=post>"
	for inputName, inputData := range formData {
		output += "\n"

		if inputData.Html != "submit" && inputData.Label != "" {
			output += " "+inputData.Label+":"
		}
		if inputData.Html == "select" {
			output += "<select name="+inputName
			if inputData.MultiSelect {
				output += " multiple"
			}
			output += ">"
			for _, option := range inputData.Select {
				output += fmt.Sprintf("\n\t<option>%v</option>", option)
			}
			for value, option := range inputData.SelectValues {
				output += fmt.Sprintf("\n\t<option value=%v>%v</option>", addQuotes(value), option)
			}
			output += "</select>"
		}else {
			output += "<input"
			if inputData.Html != "text" {
				output += " type="+inputData.Html
			}
			if inputData.Html != "submit" {
				output += " name="+inputName
			}
			if inputData.Checked {
				output += " checked"
			}
			if inputData.Html == "submit" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Label))
			}
			if inputData.Value != "" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Value))
			}
			output += ">"
		}
		if inputData.Help != "" {
			output += fmt.Sprintf("<abbr title=%v>?</abbr>", addQuotes(inputData.Help))
		}
	}
	output += "</form>"
	return output
}
