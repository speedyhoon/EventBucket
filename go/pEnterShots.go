package main

import (
	"fmt"
	"net/http"
	"strings"
)

func enterShotsAll(w http.ResponseWriter, r *http.Request, parameters string) {
	enterShots(w, r, true, parameters)
}

func enterShotsIncomplete(w http.ResponseWriter, r *http.Request, parameters string) {
	enterShots(w, r, false, parameters)
}

func enterShots(w http.ResponseWriter, r *http.Request, showAll bool, parameters string) {
	//eventID/rangeID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])

	//If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var currentRange Range
	currentRange, err = eventRange(event.Ranges, ids[1], w, r)
	if err != nil {
		return
	}

	var firstShooterID uint
	for _, shooter := range event.Shooters {
		if !shooter.Hidden && (showAll || !showAll && shooter.Scores[currentRange.StrID()].Total > 0) {
			firstShooterID = shooter.ID
			break
		}
	}

	templater(w, page{
		Title:   "Enter Shots",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		Data: map[string]interface{}{
			"Range":          currentRange,
			"Event":          event,
			"URL":            "enter-shots",
			"ShowAll":        showAll,
			"Disciplines":    globalDisciplines,
			"firstShooterID": firstShooterID,
		},
	})
}

func updateShotScores( /*w http.ResponseWriter, r *http.Request,*/ fields []field /*, redirect func()*/) string {
	event, err := getEvent(fields[1].Value)
	if err != nil {
		//		fmt.Fprintf(w, "Event with id %v doesn't exist", fields.Fields[1].Value)
		//		http.NotFound(w, r)
		//		return
		return fmt.Sprintf("Event with id %v doesn't exist", fields[1].Value)
	}

	shooterID := fields[3].valueUint
	if shooterID != event.Shooters[shooterID].ID {
		//		fmt.Fprintf(w, "Shooter with id %v doesn't exist in Event with id %v", shooterID, event.ID)
		//		http.NotFound(w, r)
		//		return
		return fmt.Sprintf("Shooter with id %v doesn't exist in Event with id %v", shooterID, event.ID)
	}
	shooter := event.Shooters[shooterID]

	//Calculate the score with the shots given
	newScore := calcTotalCenters(fields[0].Value, globalGrades[shooter.Grade].ClassID)

	err = updateDocument(tblEvent, event.ID, &shooterScore{
		rangeID: fields[2].Value,
		id:      shooterID,
		score:   newScore,
	}, &Event{}, upsertScore)

	//Display any upsert errors onscreen.
	if err != nil {
		//		fmt.Fprint(w, err.Error())
		//		return
		return err.Error()
	}

	//Return the score to the client
	if newScore.Centers == 0 {
		//		fmt.Fprint(w, newScore.Total)
		//		return
		return fmt.Sprintf("%v", newScore.Total)
	}
	//	fmt.Fprintf(w, "%v<sup>%v</sup>", newScore.Total, newScore.Centers)
	return fmt.Sprintf("%v<sup>%v</sup>", newScore.Total, newScore.Centers)
}

//This function assumes all validation on input "shots" has at least been done!
//AND input "shots" is verified to contain all characters in settings[class].validShots!
func calcTotalCenters(shots string, classID uint) Score {
	//TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot
	var total, centers uint
	var countBack string
	defaultClassSettings := defaultGlobalDisciplines()
	if classID >= 0 && classID < uint(len(defaultClassSettings)) {

		//Ignore the first sighter shots from being added to the total score. Unused sighters should be still be present in the data passed
		//for _, shot := range strings.Split(shots[defaultClassSettings[classID].QtySighters:], "") {
		for _, shot := range strings.Split(shots, "") {
			total += defaultClassSettings[classID].Marking.Shots[shot].Value
			centers += defaultClassSettings[classID].Marking.Shots[shot].Center

			//Append count back in reverse order so it can be ordered by the last few shots
			countBack = defaultClassSettings[classID].Marking.Shots[shot].CountBack + countBack
			/*if shot == "-" {
				warning = legendIncompleteScore
			}*/
		}
	}
	/*if warning == 0 && isScoreHighestPossibleScore(classID, 10, total, centers) {
		warning = legendHighestPossibleScore
	}*/
	return Score{Total: total, Centers: centers, Shots: shots, CountBack: countBack /*, warning: warning*/}
}

/*
func updateShotScores2(w http.ResponseWriter, r *http.Request) {
	validatedValues := checkForm(startShootingForm("", "", "", "").inputs, r)
	rangeID, rangeErr := strconv.Atoi(validatedValues["rangeid"])
	shooterID, shooterErr := strconv.Atoi(validatedValues["shooterid"])
	event, eventErr := getEvent(validatedValues["eventid"])

	//Check the score can be saved
	if rangeErr != nil || shooterErr != nil || eventErr != nil || rangeID >= len(event.Ranges) || event.Ranges[rangeID].Locked || event.Ranges[rangeID].IsAgg {
		http.NotFound(w, r)
		return
	}

	//Calculate the score based on the shots given
	newScore := calcTotalCenters(validatedValues["shots"], grades()[event.Shooters[shooterID].Grade].classID)

	//Return the score to the client
	if newScore.Centers > 0 {
		fmt.Fprintf(w, "%v<sup>%v</sup>", newScore.Total, newScore.Centers)
	} else {
		fmt.Fprintf(w, "%v", newScore.Total)
	}

	go triggerScoreCalculation(newScore, rangeID, event.Shooters[shooterID], event)
}*/
