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
		AutoInc: AutoInc{
			Mound: 1,
		},
	})
	if err != nil {
		//TODO add error problems to form.
		redirect()
		return
	}
	http.Redirect(w, r, "/club/"+ID, http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {
	clubID := submittedFields[0].Value
	err := updateDoc(tblClub, clubID, M{
		schemaName:      submittedFields[1].Value,
		schemaAddress:   submittedFields[2].Value,
		schemaTown:      submittedFields[3].Value,
		schemaPostCode:  submittedFields[4].Value,
		schemaLatitude:  submittedFields[5].Value,
		schemaLongitude: submittedFields[6].Value,
	})
	if err != nil {
		//TODO add error problems to form.
		redirect()
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedFields []field, redirect func()) {
	clubID := submittedFields[0].Value

	info.Println(clubID)
	info.Println(submittedFields[1].Value)
	info.Println(submittedFields[2].internalValue.(int))
	info.Println(submittedFields[3].Value)

	err := updateDoc(tblClub, clubID, M{"$push": M{
		schemaMound: Mound{
			Name:     submittedFields[1].Value,
			Distance: submittedFields[2].internalValue.(int),
			Unit:     submittedFields[3].Value,
		},
	}})
	if err != nil {
		//TODO add error problems to form.
		redirect()
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}
