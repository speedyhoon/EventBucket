package main

import (
	"fmt"
	"net/http"
	"strings"
)

func scorecardsAll(w http.ResponseWriter, r *http.Request, parameters string) {
	scorecards(w, r, true, parameters)
}

func scorecardsIncomplete(w http.ResponseWriter, r *http.Request, parameters string) {
	scorecards(w, r, false, parameters)
}

func scorecards(w http.ResponseWriter, r *http.Request, showAll bool, parameters string) {
	// eventID/rangeID
	ids := strings.Split(parameters, "/")
	event, err := getEvent(ids[0])

	// If event not found in the database return error event not found (404).
	if err != nil {
		errorHandler(w, r, http.StatusNotFound, "event")
		return
	}

	var currentRange Range
	currentRange, err = eventRange(event.Ranges, ids[1], w, r)
	if err != nil {
		return
	}

	// TODO this is messy as hell - calculate it during class grade changes?
	var classesAdded []uint
	var myDosc []Discipline
	var temp bool
	for _, id := range event.Grades {
		temp = false
		for _, added := range classesAdded {
			if id == added {
				temp = true
				break
			}
		}
		if !temp {
			classesAdded = append(classesAdded, id) // TODO messy
			myDosc = append(myDosc, globalDisciplines[findGrade(id).ClassID])
		}
	}
	var longestShoot uint
	var longestClassID int
	for i, discipline := range myDosc {
		if discipline.QtySighters+discipline.QtyShots > longestShoot {
			longestShoot = discipline.QtySighters + discipline.QtyShots
			longestClassID = i
		}
	}
	for i, discipline := range myDosc {
		if longestShoot-discipline.QtySighters-discipline.QtyShots > 1 {
			myDosc[i].Colspan = longestShoot - discipline.QtySighters - discipline.QtyShots + 1
		}
	}

	templater(w, page{
		Title:   "Scorecards",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: event.Name,
		JS:      []string{"startShooting"},
		Data: map[string]interface{}{
			"Range":   currentRange,
			"Event":   event,
			"URL":     "scorecards",
			"ShowAll": showAll,

			// TODO this is messy as hell
			"Disciplines":    myDosc,
			"LongestClassID": longestClassID,
			"LongestShoot":   longestShoot,
			"Sighters":       make([]struct{}, globalDisciplines[longestClassID].QtySighters+1),
			"Shots":          make([]struct{}, globalDisciplines[longestClassID].QtyShots+1),
		},
	})
}

func updateShotScores(w http.ResponseWriter, r *http.Request, submittedForm form, redirect func()) {
	event, err := getEvent(submittedForm.Fields[1].Value)
	if err != nil {
		fmt.Fprintf(w, "Event with id %v doesn't exist", submittedForm.Fields[1].Value)
		http.NotFound(w, r)
		return
	}

	shooterID := submittedForm.Fields[3].valueUint
	if shooterID != event.Shooters[shooterID].ID {
		fmt.Fprintf(w, "Shooter with id %v doesn't exist in Event with id %v", shooterID, event.ID)
		http.NotFound(w, r)
		return
	}
	shooter := event.Shooters[shooterID]

	// Calculate the score with the shots given
	newScore := calcTotalCenters(submittedForm.Fields[0].Value, globalGrades[shooter.Grade].ClassID)

	err = updateDocument(tblEvent, event.ID, &shooterScore{
		rangeID: submittedForm.Fields[2].Value,
		id:      shooterID,
		score:   newScore,
	}, &Event{}, upsertScore)

	// Display any upsert errors onscreen.
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	// Return the score to the client
	if newScore.Centers == 0 {
		fmt.Fprint(w, newScore.Total)
		return
	}
	fmt.Fprintf(w, "%v<sup>%v</sup>", newScore.Total, newScore.Centers)
}

// This function assumes all validation on input "shots" has at least been done!
// AND input "shots" is verified to contain all characters in settings[class].validShots!
func calcTotalCenters(shots string, classID uint) Score {
	// TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot
	var total, centers uint
	var countBack string
	defaultClassSettings := defaultGlobalDisciplines()
	if classID >= 0 && classID < uint(len(defaultClassSettings)) {

		// Ignore the first sighter shots from being added to the total score. Unused sighters should be still be present in the data passed
		// for _, shot := range strings.Split(shots[defaultClassSettings[classID].QtySighters:], "") {
		for _, shot := range strings.Split(shots, "") {
			total += defaultClassSettings[classID].Marking.Shots[shot].Value
			centers += defaultClassSettings[classID].Marking.Shots[shot].Center

			// Append count back in reverse order so it can be ordered by the last few shots
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

	// Check the score can be saved
	if rangeErr != nil || shooterErr != nil || eventErr != nil || rangeID >= len(event.Ranges) || event.Ranges[rangeID].Locked || event.Ranges[rangeID].IsAgg {
		http.NotFound(w, r)
		return
	}

	// Calculate the score based on the shots given
	newScore := calcTotalCenters(validatedValues["shots"], grades()[event.Shooters[shooterID].Grade].classID)

	// Return the score to the client
	if newScore.Centers > 0 {
		fmt.Fprintf(w, "%v<sup>%v</sup>", newScore.Total, newScore.Centers)
	} else {
		fmt.Fprintf(w, "%v", newScore.Total)
	}

	go triggerScoreCalculation(newScore, rangeID, event.Shooters[shooterID], event)
}*/
