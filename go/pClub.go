package main

import "net/http"

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	//	sessionForm := getSession(w, r)
	//	sessionFields := sessionForm.fields

	templater(w, page{
		Title:  "Club",
		MenuID: clubID,
		menu:   urlClub,
		//		Data: M{
		//			"NewClub": sessionFields,
		/*[]field{
				{
					Error:   "hey hey!",
					Value:   "true",
					Options: []option{},
				}, {
					Error: "I caused an error!@",
					Value: "fds",
					Options: []option{
						{Label: "label", Value: "2 3"},
						{Label: "text", Value: `"T`},
						{Label: "search", Value: ">S"},
					},
				}, {
					Value: "AbC",
				},
			},
		}*/
		//		},
		///
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r)
	sessionFields := sessionForm.fields

	templater(w, page{
		Title: "Clubs",
		Data: M{
			"NewClub": sessionFields,
			/*[]field{
					{
						Error:   "hey hey!",
						Value:   "true",
						Options: []option{},
					}, {
						Error: "I caused an error!@",
						Value: "fds",
						Options: []option{
							{Label: "label", Value: "2 3"},
							{Label: "text", Value: `"T`},
							{Label: "search", Value: ">S"},
						},
					}, {
						Value: "AbC",
					},
				},
			}*/
		},
	})
}
