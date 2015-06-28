package main

import (
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
	if viewController.Theme == templateHome {
		viewController.Data["Menu"] = homeMenuItems
	}
	viewController.Data["Title"] = viewController.Title
	viewController.Data["CurrentYear"] = time.Now().Year()
	//Search in Theme html file & replace "^^BODY^^" with TemplateFile
	source := strings.Replace(loadHTM(viewController.Theme), "^^BODY^^", loadHTM(viewController.TemplateFile), -1)
	generator(w, source, viewController)
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
			return grades()[grade].className
		},
		"CLASSLONG": func(grade int) string {
			return grades()[grade].longName
		},
		"JSCLASS": func(grade int) string {
			return fmt.Sprintf("%v", grades()[grade].classID)
		},
		"GRADE": func(grade int) string {
			return grades()[grade].name
		},
		"Fieldset": func(title string) template.HTML {
			return template.HTML(fieldSet(title))
		},
		"EndFieldset": func() template.HTML {
			return template.HTML("</fieldset>")
		},
		"COLSPAN": func(longestShots []string, shortShots int) template.HTMLAttr {
			if len(longestShots) > shortShots {
				return template.HTMLAttr(fmt.Sprintf(" colspan=%v", len(longestShots)-shortShots+1))
			}
			return template.HTMLAttr("")
		},
		"ElementClass": func(className1, className2 interface{}) template.HTMLAttr {
			className3 := fmt.Sprintf("%v", className2)
			if className1 != "" && className3 != "" && className3 != "0" {
				return template.HTMLAttr(fmt.Sprintf(" class=%v%v", className1, className3))
			}
			return template.HTMLAttr("")
		},
		"DisplayShot": func(shotIndex int, score Score) string {
			if len(score.Shots) > shotIndex {
				return shotsToValue(string(score.Shots[shotIndex]))
			}
			return ""
		},
		"START_SHOOTING_SHOTS": func(score Score) template.HTML {
			var output string
			for _, shot := range strings.Split(score.Shots, "") {
				output += fmt.Sprintf("<td>%v", shotsToValue(shot))
			}
			return template.HTML(output)
		},
		"NOSHOOTERS": func() template.HTML {
			return template.HTML("<p>No Shooters entered in this event.</p>")
		},
		"VAR2STR": func(input interface{}) string {
			return fmt.Sprintf("%v", input)
		},
		"AgeGroupDisplay": func(value string) string {
			return ageGroupDisplay(value)
		},
		"POSITION": func(score Score) template.HTMLAttr {
			if score.Total == 0 && score.Centres == 0 {
				return template.HTMLAttr(fmt.Sprintf(" class=w%v", legendNoScore))
			}
			if score.Warning != 0 {
				return template.HTMLAttr(fmt.Sprintf(" class=w%v", score.Warning))
			}
			if score.Position > 0 {
				return template.HTMLAttr(fmt.Sprintf(" class=p%v", score.Position))
			}
			return template.HTMLAttr("")
		},
	})
	t := template.Must(html.Parse(fillin))
	err := t.Execute(w, viewController.Data)
	if err != nil {
		warning.Println(err)
	}
}

var eventMenuItems = []Menu{
	{
		Name: "Home",
		Link: "/",
	}, {
		Name: "Event",
		Link: urlEvent,
	}, {
		Name: "Event Settings",
		Link: urlEventSettings,
	}, {
		Name: "Scoreboard",
		Link: urlScoreboard,
	}, {
		Name:   "Start Shooting",
		Link:   urlStartShooting,
		Ranges: true,
	}, {
		Name:   "Total Scores",
		Link:   urlTotalScores,
		Ranges: true,
	}, {
		Name: "Close Menu",
		Link: "#",
	},
}

func eventMenu(eventID string, eventRanges []Range, pageURL string, isPrizeMeet bool) string {
	menu := "<ul>"
	var selected string
	var closeMenu bool
	for _, menuItem := range eventMenuItems {
		if menuItem.Link == pageURL {
			selected = " class=v"
		}
		if menuItem.Ranges {
			if (isPrizeMeet && menuItem.Link != urlTotalScores) || !isPrizeMeet {
				//The a tag is needed for my ipad
				if len(eventRanges) >= 1 {
					var menuRangeItems string
					for rangeID, rangeItem := range eventRanges {
						if !rangeItem.IsAgg && !rangeItem.Hidden {
							menuRangeItems += fmt.Sprintf("<li><a href=%v%v/%v>%v - %v</a>", menuItem.Link, eventID, rangeID, rangeID, rangeItem.Name)
						}
					}
					if menuRangeItems != "" {
						menu += fmt.Sprintf("<li%v><a href=#>%v</a><ul>%v</ul>", selected, menuItem.Name, menuRangeItems)
						closeMenu = true
					}
				}
			}
		} else if menuItem.Name == "Close Menu" {
			//Don't show the close menu item when there are no ranges available
			if len(eventRanges) >= 1 && closeMenu {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a>", selected, addQuotes(menuItem.Link), menuItem.Name)
			}
		} else {
			if menuItem.Link[len(menuItem.Link)-1:] == "/" && menuItem.Link != "/" {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a>", selected, addQuotes(menuItem.Link+eventID), menuItem.Name)
			} else {
				menu += fmt.Sprintf("<li%v><a href=%v>%v</a>", selected, addQuotes(menuItem.Link), menuItem.Name)
			}
		}
		selected = ""
	}
	if pageURL == urlScoreboard {
		menu += "<li><a id=scoreSettings href=#scoreboard_settings onclick=\"var d=document.getElementById('scoreboard_settings');d.style.display=(d.style.display?'':'block')\">&nbsp;</a>"
	}
	return menu + "</ul>"
}
