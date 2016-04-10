package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request) {
	action, pageForms := sessionForms(w, r, shooterNew, shooterSearch, shooterDetails)
	var shooters []Shooter
	var err error
	var shooterQty uint
	if action == shooterSearch {
		shooters, err, shooterQty = getSearchShooters(pageForms[1].Fields[0].Value, pageForms[1].Fields[1].Value, pageForms[1].Fields[2].Value)
	} else {
		shooterQty, err = collectionSize(tblShooter)
	}
	//TODO add query string so search is bookmarkable?

	templater(w, page{
		Title: "Shooters",
		Error: err,
		JS:    "shooterDetails",
		Data: map[string]interface{}{
			"NewShooter":     pageForms[0],
			"ListShooters":   shooters,
			"ShooterSearch":  pageForms[1],
			"ShooterDetails": pageForms[2],
			"QtyShooters":    shooterQty,
			"Grades":         globalGradesDataList,
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	err := updateDocument(tblShooter, submittedForm.Fields[5].Value, &Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].valueUint,
		AgeGroup:  submittedForm.Fields[4].valueUint,
	}, &Shooter{}, updateShooterDetails)
	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
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

func searchShooters(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	setSession(w, submittedForm)
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}

func shooterInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	shooter := Shooter{
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].valueUint,
		AgeGroup:  submittedForm.Fields[4].valueUint,
	}
	_, err := insertShooter(shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}
