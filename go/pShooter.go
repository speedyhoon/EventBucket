package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Shooters",
		Data: M{
			"Stuff": "SHOOTERS page!",
			"Fds": []field{
				{
					Error: "I caused an error!@",
					Options: []option{
						{Label: "label", Value: "2 3"},
						{Label: "text", Value: `"T`},
						{Label: "search", Value: ">S"},
					},
				},
				{
					Options: []option{
						{Label: "Warrack", Value: "R23"},
						{Label: "Horsham", Value: "T52"},
						{Label: "Stawell", Value: "S82"},
					},
				},
				//			Date{},
				//			Time{},
				//			Check{},
				//			Hidden{},
			},
		},
	})
}
