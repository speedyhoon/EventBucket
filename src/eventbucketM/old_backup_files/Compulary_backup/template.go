package main

import (
	"net/http"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
//	"os"
)

func generator(w http.ResponseWriter, fillin string, data map[string]interface{}) {
	my_html := template.New("my_template").Funcs(template.FuncMap {
		"ESC": func(x string) template.HTML {
//			dump(x)
//			if x != "" {
				return template.HTML(x)
//			}else {
//				return ""
//			}
		},

		//TODO I don't think these is used at all
//		"SCHEMA": func(attribute, x string) string {
//			return schema(attribute)
//		},
//		"EMPTY": func(x string) string {
//			if strings.len(x) > 0{
//				return x
//			}
//			return ""
//		},


	})
	t := template.Must(my_html.Parse(fillin))
//	t, err := my_html.ParseFiles(`htm\organisers.htm`)
//	checkErr(err)
	err := t.Execute(w, data)
	if err != nil {
		fmt.Print("executing template:", err)
	}
}

func templator(template string, fill_in string, data map[string]interface{}, w http.ResponseWriter) {
	switch (template){
	case "home":
		generator(w, minifyHtml("templateHome", templateHome(fill_in)), data)
		break
	case "admin":
		generator(w, minifyHtml("templateAdmin", new_template_admin2(fill_in)), data)
//		generator(w, minifyHtml("templateAdmin", templateAdmin(fill_in)), data)
		break
	case "ajax":
		generator(w, fill_in, data)
	}
}

func loadHTM(page_name string)string {
//	bytes, err := ioutil.ReadFile("htm/" + page_name + ".htm")
//	if DEV {
//		bytes, err = ioutil.ReadFile("html/" + page_name + ".html")
//	}
//	file := string(bytes)
//	if err != nil{	//File does not exist
		bytes, err := ioutil.ReadFile("html/"+page_name+".html")
		checkErr(err)
		file := minifyHtml(page_name, string(bytes))
		remove_chars := map[string]string{
			"~~CSSEXT~~": css_extension,
			"~~ICONEXT~~": icon_extension,
		}
		for search, replace := range remove_chars {
			file = strings.Replace(file, search, replace, -1)
		}
		err = ioutil.WriteFile("htm/"+page_name+".htm", []byte(file), 0777)
		checkErr(err)
//	}
	return file
}

func minifyHtml(page_name, html string)string{
	if MINIFY {
		minify := html
		remove_chars := map[string]string {
			"	": "",	//Tab
			"\n": "",	//new line
			"\r": "",	//carriage return
		}
		replace_chars := map[string]string {	//TODO remove spaces between block elements like: </dov> <div> but keep between inline elements like </span> <span>
			"  ": " ",	//double spaces
			"type=text": "",
			"type=\"text\"": "",
			"type='text'": "",
			" >": ">",
//			" <": "<",
			"< ": "<",
			">  <": "> <",
			" />": "/>",
			"/ >": "/>",
			"<br/>": "<br>",
			"</br>": "<br>",
			"<br />": "<br>",
			//"\"": "",
			//"'": "",
		}
		for search, replace := range remove_chars {
			minify = strings.Replace(minify, search, replace, -1)
		}

		repeat := true
		for repeat {
			backup := minify
			for search, replace := range replace_chars {
				length := len(minify)
				minify = strings.Replace(minify, search, replace, -1)
				if length != len(minify){
					log("A dodgy character (%v) was found in the source! Please replace with (%v).", search, replace)
				}
			}
			if minify == backup {
				break
			}
		}
		minify = strings.Replace(minify, "~~~", " ", -1)

		minify_len := len(minify)
		html_len := len(html)
		if minify_len != html_len {
			log("Page '%v' had %v bytes removed (%v percent), total: %v, from: %v", page_name, html_len-minify_len, minify_len*100/html_len, minify_len, html_len)
		}else{
			log("Page '%v' OK", page_name)
		}
		return minify
	}
	return html
}

func templateHome(body string) string {
	current_year := fmt.Sprintf("%v", time.Now().Year())
	return `<!doctype html>
<html>
<head>
	<link rel=stylesheet href=/css/{{if .Css}}{{.Css}}{{else}}home{{end}}`+css_extension+`>
	<link rel=icon type=image/png href=/i/{{if .Ico}}{{.Ico}}{{else}}icon{{end}}`+icon_extension+`>
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<div id=a>
	<h1><a href=/>EventBucket</a></h1>
	<span>Making scoring easier for rifle clubs!</span>
	<ul id=menu>
		{{range $link, $href := .Menu}}
			<li><a href=/{{$href}}>{{$link}}</a></li>
		{{end}}
	</ul>
	<p>EventBucket is designed to assist with collecting scores during shooting events. Has your rifle club ever had difficulty keeping up with scoring for weekly matches, pennants or even prize meetings? EventBucket will solve these issues. It does not require electronic targets and can complement existing scoring and check scoring using cards and blackboards.</p>
	<h2>{{.PageName}}</h2>
	`+body+`
	<p id=b>Copyright &copy; 2011-`+current_year+` EventBucket</p>
</div>
</body></html>`
}

func new_template_admin2(body string)string{
	source := loadHTM("admin_template")
	return strings.Replace(source, "~~BODY~~", body, -1)
}

func templateAdmin(body string) string {
	return `<!doctype html>
<html>
<head>
	<link rel=stylesheet href=/css/{{if .Css}}{{.Css}}{{else}}main{{end}}`+css_extension+`>
	<link rel=icon type=image/png href=/i/{{if .Ico}}{{.Ico}}{{else}}icon{{end}}`+icon_extension+`>
	<title>EventBucket{{if .Title}} - {{.Title}}{{end}}</title>
</head>
<body>
	` + body + `
	{{if .Js}}<script src={{.Js}}></script>{{end}}
</body>
</html>`
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
		output := ""

		if inputData.Html != "submit" && inputData.Label != "" {
			output += " "+inputData.Label+": "
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
	output := fmt.Sprintf("<form action=%v method=post>", addQuotes("/"+action))
	for inputName, inputData := range formData {
		if inputData.Html != "submit" && inputData.Label != "" {
			output += " "+inputData.Label+": "
		}
		if inputData.Html == "select" {
			output += "<select name="+inputName

			if inputData.MultiSelect {
				output += " multiple"
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Help != "" {
				output += fmt.Sprintf("title=%v", addQuotes(inputData.Help))
			}
			output += ">"
			//TODO <option value="" disabled selected>Select your option</option>
			if inputData.Placeholder != "" {
				output+= fmt.Sprintf("<option value disabled selected>%v</option>", inputData.Label)
			}
			for _, option := range inputData.Select {
				output += fmt.Sprintf("<option>%v</option>", option)
			}
			for value, option := range inputData.SelectValues {
				output += fmt.Sprintf("<option value=%v>%v</option>", addQuotes(value), option)
			}
			output += "</select>"
		}else if inputData.Html == "submit" {
			output += "<button>"+inputData.Label+"</button>"
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
			if inputData.Placeholder != "" {
				output += " placeholder="+inputData.Placeholder
			}
			if inputData.Min > 0 {
				output += " min="+echo(inputData.Min)
			}
			if inputData.Required {
				output += " required"
			}
			if inputData.Help != "" {
				output += fmt.Sprintf("title=%v", addQuotes(inputData.Help))
			}
			if inputData.Html == "submit" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Label))
			}
			if inputData.Value != "" {
				output += fmt.Sprintf(" value=%v", addQuotes(inputData.Value))
			}
			output += ">"
		}
	}
	output += "</form>"
	return output
}
