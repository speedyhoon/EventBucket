package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"strings"
)

func templator(viewController Page, w http.ResponseWriter, r *http.Request) {
	if viewController.v8Url != nil {
		none := viewController.v8Url.FindStringSubmatch(r.URL.Path)
		if none == nil {
			http.NotFound(w, r)
			return
		}
	}

	viewController.Data["DirCss"]		= DIR_CSS	//TODO remove & replace with folder name via build script directly into the html files
	viewController.Data["DirJpeg"]	= DIR_JPEG	//TODO remove & replace with folder name via build script directly into the html files
	viewController.Data["DirJs"]		= DIR_JS		//TODO remove & replace with folder name via build script directly into the html files
	viewController.Data["DirPng"]		= DIR_PNG	//TODO remove & replace with folder name via build script directly into the html files
	viewController.Data["DirSvg"]		= DIR_SVG	//TODO remove & replace with folder name via build script directly into the html files
	//	viewController.Data["DirWebp"]= DIR_WEBP	//TODO remove & replace with folder name via build script directly into the html files
	viewController.Data["Favicon"]	= "/p/a"		//TODO remove & replace with hashed filename via build script directly into the html files
	viewController.Data["Title"] = viewController.Title
	viewController.Data["CurrentYear"] = time.Now().Year()	//TODO remove & replace with current year via build script directly into the html files
	viewController.Data["NewRelic"] = NEWRELIC	//TODO replace with the NEWRELIC html template via build script directly into the html files

	//Search in Theme html file & replace "^^BODY^^" with TemplateFile
	source := bytes.Replace(loadHTM(viewController.Theme), []byte("^^BODY^^"), loadHTM(viewController.TemplateFile), -1)
	source = bytes.Replace(source, []byte("^^NetworkAdaptor^^"), loadHTM("NetworkAdaptor"), -1)
	generator(w, string(source), viewController)
}

func generator(w http.ResponseWriter, fillin string, viewController Page) {
	my_html := template.New(viewController.TemplateFile+"Template").Funcs(template.FuncMap{
		"HTM": func(x string) template.HTML {
			return template.HTML(x)
		},
		"HTMattr": func(value string) template.HTMLAttr {
			return template.HTMLAttr(value)
		},
		"JS": func(x string) template.JS {
			return template.JS(x)
		},
		"CLASS": func(grade int) string {
			return grades()[grade].ClassName
		},
		"CLASSLONG": func(grade int) string {
			return grades()[grade].LongName
		},
		"JSCLASS": func(grade int) string {
			return fmt.Sprintf("%v", grades()[grade].ClassId)
		},
		"GRADE": func(grade int) string {
			return grades()[grade].Name
		},
		"Fieldset": func(title string) template.HTML {
			return template.HTML(field_set(title))
		},
		"EndFieldset": func() template.HTML {
			return template.HTML("</fieldset>")
		},
		"COLSPAN": func(longest_shots []string, short_shots int) template.HTMLAttr {
			if len(longest_shots) > short_shots {
				return template.HTMLAttr(fmt.Sprintf(" colspan=%v", len(longest_shots)-short_shots+1))
			}
			return template.HTMLAttr("")
		},
		"CSSclass": func(class_name1, class_name2 interface{}) template.HTMLAttr {
			if class_name1 != "" && class_name2 != "" {
				return template.HTMLAttr(fmt.Sprintf(" class=%v%v", class_name1, class_name2))
			}
			return template.HTMLAttr("")
		},
		"DisplayShot": func(shot_index int, score Score) string {
			if len(score.Shots) > shot_index {
//				return fmt.Sprintf("%s", ShotsToValue(string(score.Shots[shot_index])))
				return ShotsToValue(string(score.Shots[shot_index]))
			}
			return ""
		},
		"START_SHOOTING_SHOTS": func(score Score) template.HTML {
			var output string
			for _, shot := range strings.Split(score.Shots, "") {
				output += fmt.Sprintf("<td>%v</td>", ShotsToValue(shot))
			}
			return template.HTML(output)
		},
		"NOSHOOTERS": func() template.HTML {
			return template.HTML("<p>No Shooters entered in this event.</p>")
		},
		"VAR2STR": func(input interface{}) string {
			return fmt.Sprintf("%v", input)
		},
		"POSITION": func(position int) template.HTMLAttr {
			if position > 0{
				return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
			}
			return template.HTMLAttr("")
		},
	})
	t := template.Must(my_html.Parse(fillin))
	err := t.Execute(w, viewController.Data)
	if err != nil {
		Warning.Println(err)
	}
}

type Menu struct {
	Name, Link string
	Ranges     bool
}

var EventMenuItems = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
	Menu{
		Name: "Event",
		Link: URL_event,
	},
	Menu{
		Name: "Event Settings",
		Link: URL_eventSettings,
	},
	Menu{
		Name: "Scoreboard",
		Link: URL_scoreboard,
	},
	Menu{
		Name:   "Total Scores",
		Link:   URL_totalScores,
		Ranges: true,
	},
	Menu{
		Name:   "Start Shooting",
		Link:   URL_startShooting,
		Ranges: true,
	},
	/*Menu{
		Name: "Close Menu",
		Link: "#",
	},*/
}

func event_menu(event_id string, event_ranges []Range, page_url string, isPrizeMeet bool) string {
	menu := "<ul>"
	selected := ""
	for _, menu_item := range EventMenuItems {
		if menu_item.Link == page_url{
			selected = " class=v"
		}

		if menu_item.Ranges {
			class := "m"
			if menu_item.Link == page_url{
				class = `"v m"`
			}
			if (isPrizeMeet && menu_item.Link != URL_totalScores) || !isPrizeMeet{
				//The a tag is needed for my ipad
				menu += fmt.Sprintf("<li class=%v><a href=#>%v</a><ul>", class, menu_item.Name)
				for range_id, range_item := range event_ranges {
					if len(range_item.Aggregate) == 0 && !range_item.Hidden {
						menu += fmt.Sprintf("<li><a href=%v%v/%v>%v - %v</a></li>", menu_item.Link, event_id, range_id, range_id, range_item.Name)
					}
				}
				menu += "</ul></li>"
			}
		} else {
			if menu_item.Link == "/"{
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menu_item.Link), menu_item.Name)
			}else if menu_item.Link[len(menu_item.Link)-1:] == "/" {
				menu += fmt.Sprintf("<li%v><a href=%v%v>%v</a></li>", selected, addQuotes(menu_item.Link), event_id, menu_item.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menu_item.Link), menu_item.Name)
			}
		}
		selected = ""
	}
	menu += "</ul>"
	return menu
}

func scoreboard_menu(event_id string, event_ranges []Range, page_url string, isPrizeMeet bool) string {
	menu := "<ul>"
	selected := ""
	for _, menu_item := range EventMenuItems {
		if menu_item.Link == page_url{
			selected = " class=v"
		}

		if menu_item.Ranges {
			class := "m"
			if menu_item.Link == page_url{
				class = `"v m"`
			}
			if (isPrizeMeet && menu_item.Link != URL_totalScores) || !isPrizeMeet{
				//The a tag is needed for my ipad
				menu += fmt.Sprintf("<li class=%v><a href=#>%v</a><ul>", class, menu_item.Name)
				for range_id, range_item := range event_ranges {
					if len(range_item.Aggregate) == 0 && !range_item.Hidden {
						menu += fmt.Sprintf("<li><a href=%v%v/%v>%v - %v</a></li>", menu_item.Link, event_id, range_id, range_id, range_item.Name)
					}
				}
				menu += "</ul></li>"
			}
		} else {
			if menu_item.Link == "/"{
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menu_item.Link), menu_item.Name)
			}else if menu_item.Link[len(menu_item.Link)-1:] == "/" {
				menu += fmt.Sprintf("<li%v><a href=%v%v>%v</a></li>", selected, addQuotes(menu_item.Link), event_id, menu_item.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menu_item.Link), menu_item.Name)
			}
		}
		selected = ""
	}
	menu += "<li><a id=scoreSettings href=#scoreboard_settings onclick=\"var d=document.getElementById('scoreboard_settings');d.style.display=(d.style.display?'':'block')\">&nbsp;</a></li>"
	menu += "</ul>"
	return menu
}

func standard_menu(menu_items []Menu) string {
	menu := "<ul>"
	for _, menu_item := range menu_items {
		menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", addQuotes(menu_item.Link), menu_item.Name)
	}
	menu += "</ul>"
	return menu
}

func home_menu(page string, menu_items []Menu) string {
	menu := "<ul id=menu>"
	for _, menu_item := range menu_items {
		if page != menu_item.Link {
			menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", addQuotes(menu_item.Link), menu_item.Name)
		}else{
			menu += fmt.Sprintf("<li class=v><a href=%v>%v</a></li>", addQuotes(menu_item.Link), menu_item.Name)
		}
	}
	menu += "</ul>"
	return menu
}
