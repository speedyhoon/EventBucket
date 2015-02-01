package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

func templatePage(url string, data Page, w http.ResponseWriter) {
	templator(data.Theme, url, data.Data, w)
}

func templator(main_template string, content_template string, data M, w http.ResponseWriter) {
	//Ajax responses should not use this function! Instead use "generator(w, body, data)"
	source := loadHTM(main_template)
	remove_chars := map[string][]byte{
		"^^DIR_JS^^":       []byte(DIR_JS),
		"^^DIR_CSS^^":      []byte(DIR_CSS),
		"^^DIR_PNG^^":     []byte(DIR_PNG),
		"^^FAVICON^^":      []byte(FAVICON),
		"^^CURRENT_YEAR^^": []byte(fmt.Sprintf("%v", time.Now().Year())),
		"^^BODY^^":         loadHTM(content_template),
	}
	for search, replace := range remove_chars {
		source = bytes.Replace(source, []byte(search), replace, -1)
	}
	generator(w, string(source), data)
}

func loadHTM(pageName string) []byte {
	//TODO add all html sources to []byte constant in a new go file
	pageName = strings.Replace(pageName, "/", "", -1)
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_MINIFIED, pageName))
	checkErr(err)
	return dev_mode_loadHTM(pageName, bytes)
}

func generator(w http.ResponseWriter, fillin string, data M) {
	my_html := template.New("my_template").Funcs(template.FuncMap{
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
	checkErr(t.Execute(w, data))
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
//	Menu{
//		Name: "Close Menu",
//		Link: "#",
//	},
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

var HOME_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
	Menu{
		Name: "Archive",
		Link: URL_archive,
	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
	Menu{
		Name: "About",
		Link: URL_about,
	},
}

/*var ORGANISERS_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
	Menu{
		Name: "Archive",
		Link: URL_archive,
	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
}*/

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
