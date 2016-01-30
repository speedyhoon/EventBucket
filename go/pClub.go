package main

import "net/http"

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	//	sessionForm := getSession(w, r)
	//	sessionFields := sessionForm.fields
	club, err := getClub(clubID)
	if err != nil {
		warn.Println(err)
	}
	//TODO redirect to clubs page if clubid doesn't exist in the database

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
	listClubs, err := getClubs()
	if err != nil {
		warn.Println(err)
	}
	templater(w, page{
		Title: "Clubs",
		Data: M{
			"NewClub":   getSession(w, r).fields,
			"ListClubs": listClubs,
		},
	})
}

func insertClub(w http.ResponseWriter, r *http.Request) {
	submittedFields, isValid := isValid(r, GlobalForms[clubNew].fields)
	name := submittedFields[0].Value
	isDefault := submittedFields[1].internalValue.(bool)

	goToClubsPage := func() { http.Redirect(w, r, "/clubs", http.StatusSeeOther) }
	if !isValid {
		setSession(w, form{
			action: clubNew,
			fields: submittedFields,
		})
		goToClubsPage()
		return
	}

	//TODO these several db calls are not atomic.
	ID, err := getNextID(tblClub)
	if err != nil {
		//TODO add error problems to form.
		goToClubsPage()
		return
	}
	if collectionQty(tblClub) == 0 {
		isDefault = true
	} else if isDefault {
		//update all clubs isDefault to be false
		updateAll(tblClub, M{schemaIsDefault: true}, M{"$unset": M{schemaIsDefault: ""}})
	}
	err = upsertDoc(tblClub, "", Club{
		ID:        ID,
		Name:      name,
		IsDefault: isDefault,
	})
	if err != nil {
		//TODO add error problems to form.
		goToClubsPage()
		return
	}
	http.Redirect(w, r, "/club/"+ID, http.StatusSeeOther)
}
