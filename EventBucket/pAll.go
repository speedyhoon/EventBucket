package main

import (
	"net/http"
	"os"
)

func all(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	templater(w, page{
		Title: "_All",
		Data: map[string]interface{}{
			"Event": map[string]interface{}{
				"Stuff":   "EVENT page!",
				"EventId": "eventId",
			},
			"About": map[string]interface{}{
				"Hostname":    hostname,
				"IpAddresses": localIPs(),
			},
			"Clubs": map[string]interface{}{
				"Stuff": "CLUBS page!",
			},
			"Shooters": map[string]interface{}{
				"Stuff": "SHOOTERS page!",
				"Fds": []field{
					{
						Error: "i caused an error!@",
						Options: []option{
							{Label: "label", Value: "2 3"},
							{Label: "text", Value: `"t`},
							{Label: "search", Value: ">s"},
						},
					},
					{
						Options: []option{
							{Label: "warrack", Value: "r23"},
							{Label: "horsham", Value: "t52"},
							{Label: "stawell", Value: "s82"},
						},
					},
				},
			},
		},
	})
}

func report(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Report",
	})
}
