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
		Title:   "Club",
		MenuID:  clubID,
		Menu:    urlClubs,
		Heading: club.Name,
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
			"NewClub":   getFormSession(w, r, clubNew),
			"ListClubs": listClubs,
		},
	})
}

func clubInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	name := submittedForm.Fields[0].Value
	isDefault := submittedForm.Fields[1].internalValue.(bool)

	/*//TODO these several db calls are not atomic.
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
	})*/
	ID, err := insertClub(Club{Name: name, IsDefault: isDefault})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClub+ID, http.StatusSeeOther)
}

func clubDetailsUpsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[6].Value
	err := updateDoc(tblClub, clubID, M{
		schemaName:      submittedForm.Fields[0].Value,
		schemaAddress:   submittedForm.Fields[1].Value,
		schemaTown:      submittedForm.Fields[2].Value,
		schemaPostcode:  submittedForm.Fields[3].Value,
		schemaLatitude:  submittedForm.Fields[4].internalValue.(float64),
		schemaLongitude: submittedForm.Fields[5].internalValue.(float64),
	})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}

func clubMoundInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	clubID := submittedForm.Fields[2].Value
	err := updateDoc(tblClub, clubID, M{"$push": M{
		schemaMound: Mound{
			Distance: submittedForm.Fields[0].internalValue.(uint64),
			Unit:     submittedForm.Fields[1].Value,
		},
	}})
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlClubSettings+clubID, http.StatusSeeOther)
}

func dataListClubs(clubs []Club) []option {
	var options []option
	for _, club := range clubs {
		options = append(options, option{Label: club.Name})
	}
	return options
}

func getDataListClubs() []option {
	clubs, err := getClubs()
	if err != nil {
		warn.Println(err)
	}
	return dataListClubs(clubs)
}
