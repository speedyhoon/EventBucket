package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func templator(viewController Page, w http.ResponseWriter, r *http.Request) {
	if viewController.v8Url != nil {
		none := viewController.v8Url.FindStringSubmatch(r.URL.Path)
		if none == nil {
			http.NotFound(w, r)
			return
		}
	}
	viewController.Data["Title"] = viewController.Title
	viewController.Data["CurrentYear"] = time.Now().Year()
	//Search in Theme html file & replace "^^BODY^^" with TemplateFile
	source := bytes.Replace(loadHTM(viewController.Theme), []byte("^^BODY^^"), loadHTM(viewController.TemplateFile), -1)
	generator(w, string(source), viewController)
}

func generator(w http.ResponseWriter, fillin string, viewController Page) {
	html := template.New(viewController.TemplateFile + "Template").Funcs(template.FuncMap{
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
			return template.HTML(fieldSet(title))
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
			if position > 0 {
				return template.HTMLAttr(fmt.Sprintf(" class=p%v", position))
			}
			return template.HTMLAttr("")
		},
	})
	t := template.Must(html.Parse(fillin))
	err := t.Execute(w, viewController.Data)
	if err != nil {
		Error.Println(err)
	}
}

var EventMenuItems = []Menu{
	{
		Name: "Home",
		Link: "/",
	}, {
		Name: "Event",
		Link: URL_event,
	}, {
		Name: "Event Settings",
		Link: URL_eventSettings,
	}, {
		Name: "Scoreboard",
		Link: URL_scoreboard,
	}, {
		Name:   "Start Shooting",
		Link:   URL_startShooting,
		Ranges: true,
	}, {
		Name:   "Total Scores",
		Link:   URL_totalScores,
		Ranges: true,
	}, {
		Name: "Close Menu",
		Link: "#",
	},
}

func eventMenu(eventId string, eventRanges []Range, pageUrl string, isPrizeMeet bool) string {
	menu := "<ul>"
	selected := ""
	for _, menuItem := range EventMenuItems {
		if menuItem.Link == pageUrl {
			selected = " class=v"
		}
		if menuItem.Ranges {
			if (isPrizeMeet && menuItem.Link != URL_totalScores) || !isPrizeMeet {
				//The a tag is needed for my ipad
				if len(eventRanges) >= 1 {
					menu += fmt.Sprintf("<li%v><a href=#>%v</a><ul>", selected, menuItem.Name)
					for rangeId, range_item := range eventRanges {
						if !range_item.IsAgg && !range_item.Hidden {
							menu += fmt.Sprintf("<li><a href=%v%v/%v>%v - %v</a></li>", menuItem.Link, eventId, rangeId, rangeId, range_item.Name)
						}
					}
					menu += "</ul></li>"
				}
			}
		} else if menuItem.Name == "Close Menu" {
			//Don't show the close menu item when there are no ranges available
			if len(eventRanges) >= 1 {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menuItem.Link), menuItem.Name)
			}
		} else {
			if menuItem.Link[len(menuItem.Link)-1:] == "/" && menuItem.Link != "/" {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menuItem.Link+eventId), menuItem.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a></li>", selected, addQuotes(menuItem.Link), menuItem.Name)
			}
		}
		selected = ""
	}
	if pageUrl == URL_scoreboard {
		menu += "<li><a id=scoreSettings href=#scoreboard_settings onclick=\"var d=document.getElementById('scoreboard_settings');d.style.display=(d.style.display?'':'block')\">&nbsp;</a></li>"
	}
	return menu + "</ul>"
}

func homeMenu(page string, menuItems []Menu) string {
	menu := "<ul id=menu>"
	var attributes string
	for _, menuItem := range menuItems {
		if page == menuItem.Link {
			attributes = " class=v"
		} else {
			attributes = " href=" + addQuotes(menuItem.Link)
		}
		menu += fmt.Sprintf("<li><a%v>%v</a></li>", attributes, menuItem.Name)
	}
	return menu + "</ul>"
}
