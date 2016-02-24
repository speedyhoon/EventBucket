package main

import "net/http"

func entries(w http.ResponseWriter, r *http.Request, eventID string) {
	sessionForm := getSession(w, r, []uint8{eventShooterNew, eventShooterExisting})
	//	trace.Println("event fields len=", len(sessionForm.Fields))
	//	for i, input := range sessionForm.Fields {
	//		fmt.Println(i, input.name, input.Error)
	//	}

	listClubs, err := getClubs()

	var shooterEntry form
	switch sessionForm.action {
	case eventShooterNew:
		shooterEntry = sessionForm
		shooterEntry.Fields = append(shooterEntry.Fields[:3], append([]field{{}}, shooterEntry.Fields[3:]...)...)
	case eventShooterExisting:
		shooterEntry = sessionForm
		shooterEntry.Fields = append([]field{
			{},
			{},
			{Options: dataListClubs(listClubs)},
		}, shooterEntry.Fields...)
	default:
		//	if sessionForm.action == eventShooterNew || sessionForm.action == eventShooterExisting {
		//		shooterEntry = sessionForm
		//	} else {
		shooterEntry = form{Fields: []field{
			{}, {},
			{Options: dataListClubs(listClubs)},
			{},
			{},
			{},
			{},
			//{},
		}}
	}
	if err != nil {
		shooterEntry.Error = err.Error()
	}

	if len(shooterEntry.Fields[4].Options) == 0 {
		shooterEntry.Fields[4].Options = dataListGrades()
	}
	if len(shooterEntry.Fields[5].Options) == 0 {
		shooterEntry.Fields[5].Options = dataListAgeGroup()
	}
	shooterEntry.Fields = append(shooterEntry.Fields, field{Value: eventID})
	shooterEntry.Fields[6].Value = eventID

	event, err := getEvent(eventID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}
	templater(w, page{
		Title:   "Entries",
		Menu:    urlEvents,
		MenuID:  eventID,
		Heading: event.Name,
		Data: M{
			"Event":        event,
			"ShooterEntry": shooterEntry,
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	/*ID, err := getNextID(tblEvent)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}*/

	//Insert new event into database.
	ID, err := insertEvent(Event{
		//		ID:       ID,
		Club:     submittedForm.Fields[0].Value,
		Name:     submittedForm.Fields[1].Value,
		DateTime: submittedForm.Fields[2].Value,
	})

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEntries+ID, http.StatusSeeOther)
}
