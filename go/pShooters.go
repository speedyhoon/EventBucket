package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request) {
	_, pageForms := sessionForms2(w, r, shooterNew, shooterSearch, shooterDetails)
	shooters, _ := getShooters()

	templater(w, page{
		Title: "Shooters",
		Data: M{
			"NewShooter":     pageForms[0],
			"ListShooters":   shooters,
			"shooterSearch":  shooterSearch,
			"ShooterDetails": pageForms[1],
			"QtyShooters":    len(shooters),
		},
	})
}

func shooterUpdate(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	err := updateShooter(Shooter{
		ID:        submittedForm.Fields[0].Value,
		FirstName: submittedForm.Fields[1].Value,
		Surname:   submittedForm.Fields[2].Value,
		Club:      submittedForm.Fields[3].Value,
		Grade:     submittedForm.Fields[4].internalValue.(uint64),
		AgeGroup:  submittedForm.Fields[5].internalValue.(uint64),
	}, "")
	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}

func searchShooters(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	listShooters := []option{
		{Value: "sid", Label: "Firstname, Surname, Club"},
		{Value: "123", Label: "Tom, Dick, Harry"},
	}
	templater(w, page{
		Title: "Shooter Search",
		Ajax:  true,
		Data: M{
			"ListShooters": listShooters,
		},
	})
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
