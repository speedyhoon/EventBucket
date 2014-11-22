package main

import (
	"net/http"
)

func events(w http.ResponseWriter, r *http.Request) {
	templator("admin", events_HTML(), events_Data(), w)
}
func events_Data() map[string]interface{} {
	return map[string]interface{}{
		"Title": "Event List",
		"EventList": eventList(),
	}
}

func events_HTML() string {
	return `<h1>{{ .Title}}</h1>
	{{if .EventList}}
		<table><tr><th>Name</th></tr>
			{{with .EventList}}
				{{range .}}
					<tr><td><a href=` + addQuotes(`/event/{{.Url}}`) + `>{{.Name}}</a></td></tr>
				{{end}}
			{{end}}
		</table>
	{{end}}`
}
