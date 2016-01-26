package main

import "net/http"

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	//	sessionForm := getSession(w, r)
	//	sessionFields := sessionForm.fields
	club, err := getClub(clubID)
	if err != nil {
		warn.Println(err)
	}

	templater(w, page{
		Title:  "Club",
		MenuID: clubID,
		menu:   urlClub,
		Data: M{
			"Club": club,
		},
	})
}

func clubs(w http.ResponseWriter, r *http.Request) {
	//	var listClubs []Club
	//	var err error
	//	listClubs, err = getCollection(tblClub, listClubs)
	listClubs, err := getClubs()
	if err != nil {
		warn.Println(err)
	}
	temp := getSession(w, r).fields
	if len(temp) >= 2 {
		info.Println(temp[1].Value)
	} else {
		info.Println("checkbox has no value")
	}
	templater(w, page{
		Title: "Clubs",
		Data: M{
			"NewClub":   temp,
			"ListClubs": listClubs,
		},
	})
}

func insertClub(w http.ResponseWriter, r *http.Request) {
	trace.Println("\n\n\ninsertClub START")
	submittedFields, isValid := isValid(r, GlobalForms[2].fields)
	trace.Println(submittedFields)
	trace.Println("---")
	trace.Println(isValid)
	for _, n := range submittedFields {

		trace.Println("name=", n.name)
		trace.Println("err=", n.Error)
		trace.Println("val=", n.Value)
		trace.Println("int=", n.internalValue)
	}
	goToClubsPage := func() { http.Redirect(w, r, "/clubs", http.StatusSeeOther) }
	if !isValid {
		setSession(w, form{
			action: "2",
			fields: submittedFields,
		})
		goToClubsPage()
		return
	}

	ID, err := getNextID(tblClub)
	if err != nil {
		//TODO add error problems to form.
		goToClubsPage()
		return
	}
	err = upsertDoc(tblClub, "", Club{
		ID:        ID,
		Name:      submittedFields[0].Value,
		IsDefault: submittedFields[1].internalValue.(bool),
	})
	if err != nil {
		//TODO add error problems to form.
		goToClubsPage()
		return
	}
	http.Redirect(w, r, "/club/"+ID, http.StatusSeeOther)
}
