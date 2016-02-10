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
			"NewClub":   getSession(w, r, []int{clubNew}),
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

func dataListGrades() []option {
	return []option{
		{},
		{Value: "1", Label: "Target A"},
		{Value: "2", Label: "Target B"},
		{Value: "3", Label: "Target C"},
		{Value: "4", Label: "F Class A"},
		{Value: "5", Label: "F Class B"},
		{Value: "6", Label: "F Class Open"},
		{Value: "7", Label: "F/TR"},
		{Value: "8", Label: "Match Open"},
		{Value: "9", Label: "Match Reserve"},
		{Value: "10", Label: "303 Rifle"},
	}
}

func dataListAgeGroup() []option {
	return []option{
		{},
		{Value: "1", Label: "Junior U21"},
		{Value: "2", Label: "Junior U25"},
		{Value: "3", Label: "Veteran"},
		{Value: "4", Label: "Super Veteran"},
	}
}
