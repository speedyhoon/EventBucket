package main

import "net/http"

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	//TODO disable checkbox when isDefault == true. Only non default clubs can steal the isDefault flag.

	club, err := getClub(clubID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		warn.Println(err)
		errorHandler(w, r, http.StatusNotFound, "club")
		return
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

func clubInsert(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {
	name := submittedFields[0].Value
	isDefault := submittedFields[1].internalValue.(bool)

	//TODO these several db calls are not atomic.
	ID, err := getNextID(tblClub)
	if err != nil {
		//TODO add error problems to form.
		redirect()
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
		redirect()
		return
	}
	http.Redirect(w, r, "/club/"+ID, http.StatusSeeOther)
}

func clubDetails(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {

}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {

}
