package main

import "net/http"

var shooterQty int

func shooters(w http.ResponseWriter, r *http.Request) {
	action, pageForms := sessionForms2(w, r, shooterNew, shooterSearch, shooterDetails)
	var shooters []Shooter
	var err error
	if action != nil && *action == shooterSearch {
		shooters, err = getSearchShooters(pageForms[1].Fields[0].Value, pageForms[1].Fields[1].Value, pageForms[1].Fields[2].Value)
	}
	if shooterQty < 1 {
		totalShooters, _ := getShooters()
		shooterQty = len(totalShooters)
	}

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
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	err := updateShooter(Shooter{
		ID:        submittedForm.Fields[5].Value,
		FirstName: submittedForm.Fields[0].Value,
		Surname:   submittedForm.Fields[1].Value,
		Club:      submittedForm.Fields[2].Value,
		Grade:     submittedForm.Fields[3].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[4].internalValue.(uint64),
	}, "")
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
		Grade:     submittedForm.Fields[3].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[4].internalValue.(uint64),
	}
	_, err := insertShooter(shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}
