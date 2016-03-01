package main

import "net/http"

func shooters(w http.ResponseWriter, r *http.Request) {
	shooters, _ := getShooters()
	templater(w, page{
		Title: "Shooters",
		Data: M{
			"ListShooters": shooters,
			"QtyShooters":  len(shooters),
		},
	})
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
		AgeGroup:  submittedForm.Fields[3].internalValue.(uint64),
		Grade:     submittedForm.Fields[4].internalValue.(uint64),
	}
	_, err := insertShooter(shooter)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlShooters, http.StatusSeeOther)
}
