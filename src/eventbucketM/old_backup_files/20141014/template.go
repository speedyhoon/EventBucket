package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
//	"strings"
)

type ErrorMsg struct {
	Title, Message string
	Info           bool
}

var error_queue []ErrorMsg

func error_message(info bool, index, title, message string) {
	for _, queue_item := range error_queue {
		if queue_item.Title == title {
			return
		}
	}
	error_queue = append(error_queue, ErrorMsg{
		//	error_queue[index] = ErrorMsg{
		Title:   title,
		Message: message,
		Info:    info,
	})
	//	}
}
func remove_error(title string) {
	for index, queue_item := range error_queue {
		if queue_item.Title == title {
			error_queue[index] = ErrorMsg{}
		}
	}
}

func templator(main_template string, content_template string, data map[string]interface{}, w http.ResponseWriter) {
	//Ajax responses should not use this function! Instead use "generator(w, body, data)"
	source := loadHTM(main_template)
	error_list := []byte(render_errors())
	remove_chars := map[string][]byte{
		"^^DIR_JS^^":       []byte(DIR_JS),
		"^^DIR_CSS^^":      []byte(DIR_CSS),
		"^^DIR_ICON^^":     []byte(DIR_ICON),
		"^^FAVICON^^":      []byte(FAVICON),
		"^^CURRENT_YEAR^^": []byte(fmt.Sprintf("%v", time.Now().Year())),
		"^^ERROR^^":        error_list,
		"^^BODY^^":         loadHTM(content_template),
	}
	for search, replace := range remove_chars {
		source = bytes.Replace(source, []byte(search), replace, -1)
	}
	generator(w, string(source), data)
}

func render_errors() string {
	var output string
	if len(error_queue) >= 1 {
		for _, error := range error_queue {
			if error.Title != "" && error.Message != "" {
				var class = "error"
				if error.Info {
					class = "info"
				}
				output += fmt.Sprintf("<div class=%v><h2>%v:</h2>%v</div>", class, error.Title, error.Message)
			}
		}
		error_queue = []ErrorMsg{}
	}
	return output
}

func loadHTM(page_name string) []byte {
	//TODO add all html sources to []byte constant in a new go file
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_MINIFIED, page_name))
	checkErr(err)
	return dev_mode_loadHTM(page_name, bytes)
}

func generator(w http.ResponseWriter, fillin string, data map[string]interface{}) {
	my_html := template.New("my_template").Funcs(template.FuncMap{
		"HTM": func(x string) template.HTML {
			return template.HTML(x)
		},
		"CLASS": func(class string) string {
			return class_translation(class)
		},
		"CLASSLONG": func(class string) string {
			return class_long_translation(class)
		},
		"JSCLASS": func(class string) string {
			return js_class_translation(class)
		},
		"GRADE": func(grade string) string {
			return grade_translation(grade)
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
				return fmt.Sprintf("%s", ShotsToValue[string(score.Shots[shot_index])])
			}
			return ""
		},
		"NOSHOOTERS": func() template.HTML {
			return template.HTML(ERROR_NO_SHOOTERS)
		},
		"ERROR_NO_EVENTS": func() template.HTML {
			return template.HTML(ERROR_NO_EVENTS)
		},
		"POSITION": func(position int) template.HTMLAttr {
//			return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
						if position > 0{
			//				if position <= 3 {
								return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
			//				}
						}
						return template.HTMLAttr("")
		},
//		"START_SHOOTING_SHOTS": func(score Score) template.HTML {
//			var output string
//			for _, shot := range strings.Split(score.Shots, "") {
//				output += fmt.Sprintf("<td>%v</td>", ShotsToValue[shot])
//			}
//			return template.HTML(output)
//		},
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
//	Menu{
//		Name: "Organisers",
//		Link: URL_organisers,
//	},
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

func event_menu(event_id string, event_ranges map[string]Range, page_url string) string {
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


			//The a tag is needed for my ipad
			menu += fmt.Sprintf("<li class=%v><a href=#>%v</a><ul>", class, menu_item.Name)
			for range_id, range_item := range event_ranges {
				if len(range_item.Aggregate) == 0 && !range_item.Hidden {
					menu += fmt.Sprintf("<li><a href=%v%v/%v>%v</a></li>", menu_item.Link, event_id, range_id, range_item.Name)
				}
			}
			menu += "</ul></li>"
		} else {
			if menu_item.Link == "/"{
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, menu_item.Link, menu_item.Name)
			}else if menu_item.Link[len(menu_item.Link)-1:] == "/" {
				menu += fmt.Sprintf("<li%v><a href=%v%v>%v</a></li>", selected, menu_item.Link, event_id, menu_item.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, menu_item.Link, menu_item.Name)
			}
		}
		selected = ""
	}
	menu += "</ul>"
	return menu
}

var HOME_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
//	Menu{
//		Name: "Features",
//		Link: "/features",
//	},
//	Menu{
//		Name: "Events",
//		Link: "/events",
//	},
//	Menu{
//		Name: "Clubs",
//		Link: URL_clubs,
//	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
//	Menu{
//		Name: "Event Archive",
//		Link: URL_archive,
//	},
	Menu{
		Name: "About",
		Link: URL_about,
	},
}

var ORGANISERS_MENU_ITEMS = []Menu{
	Menu{
		Name: "Home",
		Link: "/",
	},
	Menu{
		Name: "Clubs",
		Link: URL_clubs,
	},
	Menu{
		Name: "Events",
		Link: URL_events,
	},
	Menu{
		Name: "Event Archive",
		Link: URL_archive,
	},
	Menu{
		Name: "Organisers",
		Link: URL_organisers,
	},
}

func standard_menu(menu_items []Menu) string {
	menu := "<ul>"
	for _, menu_item := range menu_items {
		menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
	}
	menu += "</ul>"
	return menu
}

func home_menu(page string, menu_items []Menu) string {
	menu := "<ul id=menu>"
	for _, menu_item := range menu_items {
		if page != menu_item.Link {
			menu += fmt.Sprintf("<li><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
		}else{
			menu += fmt.Sprintf("<li class=v><a href=%v>%v</a></li>", menu_item.Link, menu_item.Name)
		}
	}
	menu += "</ul>"
	return menu
}
