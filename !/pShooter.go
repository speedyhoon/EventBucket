package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request) {
	templater(w, page{
		Title: "Shooters",
		Data: M{
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

func searchShooters(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	listShooters := []option{
		{Value: "sid", Label: "Firstname, Surname, Club"},
		{Value: "123", Label: "Tom, Dick, Harry"},
	}
	templater(w, page{
		Title: "Shooter Search",
		Ajax:  true,
		Data: M{
			"ListShooters": listShooters,
		},
	})
}
