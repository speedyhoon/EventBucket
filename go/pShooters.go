package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request, submittedForm form, isValid bool) {
	_, pageForms := sessionForms(w, r, shooterNew)

	if !debug && submittedForm.Fields[0].Value == "" && submittedForm.Fields[1].Value == "" && submittedForm.Fields[2].Value == "" {
		submittedForm.Fields[2].Value = defaultClubName()
	}
	shooters, err, shooterQty := getSearchShooters(submittedForm.Fields[0].Value, submittedForm.Fields[1].Value, submittedForm.Fields[2].Value)

	templater(w, page{
		Title: "Shooters",
		Error: err,
		JS:    "shooterDetails",
		Data: map[string]interface{}{
			"NewShooter":    pageForms[0],
			"ListShooters":  shooters,
			"ShooterSearch": submittedForm,
			"QtyShooters":   shooterQty,
			"Grades":        globalGradesDataList,
			"AgeGroups":     dataListAgeGroup(),
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	err := updateDocument(tblShooter, submittedForm.Fields[5].Value, &Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].valueUintSlice,
		AgeGroup:  submittedForm.Fields[4].valueUint,
	}, &Shooter{}, updateShooterDetails)
	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func eventSearchShooters(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	listShooters := []option{
		{Value: "sid", Label: "Firstname, Surname, Club"},
		{Value: "123", Label: "Tom, Dick, Harry"},
	}
	templater(w, page{
		Title:    "Shooter Search",
		template: templateNone,
		Data: map[string]interface{}{
			"ListShooters": listShooters,
		},
	})
}

func shooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	shooter := Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].valueUintSlice,
		AgeGroup:  submittedForm.Fields[4].valueUint,
	}
	_, err := insertShooter(shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}
