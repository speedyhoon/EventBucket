package main

import "net/http"

func club(w http.ResponseWriter, r *http.Request, clubID string) {
	//TODO disable checkbox when isDefault == true. Only non default clubs can steal the isDefault flag.

	club, err := getClub(clubID)
	//If club not found in the database return error club not found (404).
	if err != nil {
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
			"NewClub":   getSession(w, r),
			"ListClubs": listClubs,
		},
	})
}

func clubInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	name := submittedForm.Fields[0].Value
	isDefault := submittedForm.Fields[1].internalValue.(bool)

	//TODO these several db calls are not atomic.
	ID, err := getNextID(tblClub)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	if collectionQty(tblClub) == 0 {
		isDefault = true
	} else if isDefault {
		//Update all clubs isDefault to be false
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
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, "/club/"+ID, http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[6].Value
	err := updateDoc(tblClub, clubID, M{
		schemaName:      submittedForm.Fields[0].Value,
		schemaAddress:   submittedForm.Fields[1].Value,
		schemaTown:      submittedForm.Fields[2].Value,
		schemaPostcode:  submittedForm.Fields[3].Value,
		schemaLatitude:  submittedForm.Fields[4].Value,
		schemaLongitude: submittedForm.Fields[5].Value,
	})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[3].Value
	err := updateDoc(tblClub, clubID, M{"$push": M{
		schemaMound: Mound{
			Name:     submittedForm.Fields[0].Value,
			Distance: submittedForm.Fields[1].internalValue.(int),
			Unit:     submittedForm.Fields[2].Value,
		},
	}})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}
