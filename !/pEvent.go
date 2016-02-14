package main

import "net/http"

func event(w http.ResponseWriter, r *http.Request, eventID string) {
	sessionForm := getSession(w, r, []uint8{eventShooterNew, eventShooterExisting})
	//	trace.Println("event fields len=", len(sessionForm.Fields))
	//	for i, input := range sessionForm.Fields {
	//		fmt.Println(i, input.name, input.Error)
	//	}
	var shooterEntry form
	switch sessionForm.action {
	case eventShooterNew:
		shooterEntry = sessionForm
		shooterEntry.Fields = append(shooterEntry.Fields[:3], append([]field{{}}, shooterEntry.Fields[3:]...)...)
	case eventShooterExisting:
		shooterEntry = sessionForm
		shooterEntry.Fields = append([]field{{}, {}, {}}, shooterEntry.Fields...)
	default:
		//	if sessionForm.action == eventShooterNew || sessionForm.action == eventShooterExisting {
		//		shooterEntry = sessionForm
		//	} else {
		listClubs, err := getClubs()
		shooterEntry = form{Fields: []field{
			{}, {},
			{Options: dataListClubs(listClubs)},
			{},
			{},
			{},
			{},
			//{},
		}}
		if err != nil {
			shooterEntry.Error = err.Error()
		}
	}
	shooterEntry.Fields = append(shooterEntry.Fields, field{Value: eventID})
	//	trace.Println("event fields len=", len(sessionForm.Fields))
	shooterEntry.Fields[6].Value = eventID
	//	shooterEntry.Fields[7].Value = eventID

	event, err := getEvent(eventID)
	//If club not found in the database return error club not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}
	templater(w, page{
		Title:  "Event",
		menu:   urlEvent,
		MenuID: eventID,
		Data: M{
			"Event":        event,
			"ShooterEntry": shooterEntry,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	sessionForm := getSession(w, r, []uint8{eventDetails})
	listEvents, err := getEvents()
	templater(w, page{
		Title: "Events",
		Error: err,
		Data: M{
			"NewEvent":   eventNewDefaultValues(sessionForm),
			"ListEvents": listEvents,
		},
	})
}

func eventInsert(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	ID, err := getNextID(tblEvent)
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}

	//Insert new event into database.
	err = upsertDoc(tblEvent, ID, Event{
		ID:   ID,
		Club: submittedForm.Fields[0].Value,
		Name: submittedForm.Fields[1].Value,
		Date: submittedForm.Fields[2].Value,
		Time: submittedForm.Fields[3].Value,
	})

	//Display any insert errors onscreen.
	if err != nil {
		formError(w, submittedForm, redirect, err)
		return
	}
	http.Redirect(w, r, urlEvent+ID, http.StatusSeeOther)
}

func eventNewDefaultValues(form form) form {
	if form.action != eventNew && len(form.Fields) == 0 {
		form.Fields = []field{
			{Required: hasDefaultClub()},
			{},
			{Value: defaultDate()[0]},
			{Value: defaultTime()[0]},
		}
	}
	return form
}
