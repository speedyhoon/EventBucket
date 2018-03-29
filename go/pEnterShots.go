package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/speedyhoon/forms"
)

func enterShotsAll(w http.ResponseWriter, r *http.Request, event Event, rangeID rID) {
	enterShots(w, r, true, event, rangeID)
}

func enterShotsIncomplete(w http.ResponseWriter, r *http.Request, event Event, rangeID rID) {
	enterShots(w, r, false, event, rangeID)
}

func enterShots(w http.ResponseWriter, r *http.Request, showAll bool, event Event, rangeID rID) {
	currentRange, err := eventRange(event.Ranges, rangeID, w, r)
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

	render(w, page{
		Title:   "Enter Shots",
		Menu:    urlEvents,
		MenuID:  event.ID,
		Heading: currentRange.Name,
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

func updateShotScores(fields []forms.Field) string {
	event, err := getEvent(fields[1].Str())
	if err != nil {
		return fmt.Sprintf("Event with id %v doesn't exist", fields[1].Value)
	}

	shooterID := fields[3].Uint()
	if shooterID != event.Shooters[shooterID].ID {
		return fmt.Sprintf("Shooter with id %v doesn't exist in Event with id %v", shooterID, event.ID)
	}
	shooter := event.Shooters[shooterID]

	//Calculate the score with the shots given
	newScore := calcTotalCenters(fields[0].Str(), globalGrades[shooter.Grade].ClassID)

	err = updateDocument(tblEvent, event.ID, &shooterScore{
		rangeID: fields[2].Str(),
		id:      shooterID,
		score:   newScore,
	}, &Event{}, upsertScore)

	//Display any upsert errors onscreen.
	if err != nil {
		return err.Error()
	}

	var score string

	//Return the score to the client
	if newScore.Centers == 0 {
		score = fmt.Sprintf("%v", newScore.Total)
	} else {
		score = fmt.Sprintf("%v<sup>%v</sup>", newScore.Total, newScore.Centers)
	}
	data := struct {
		R string
		S uint
		T string
	}{
		fields[2].Str(),
		shooterID,
		score,
	}

	var response []byte
	response, err = json.Marshal(data)
	if err != nil {
		warn.Println(err)
	}
	return fmt.Sprintf("%s", response)
}

//This function assumes all validation on input "shots" has at least been done!
//AND input "shots" is verified to contain all characters in settings[class].validShots!
func calcTotalCenters(shots string, classID uint) Score {
	//TODO need validation to check that the shots given match the required validation given posed by the event. e.g. sighters are not in the middle of the shoot or shot are not missing in the middle of a shoot
	var total, centers uint
	var countBack string
	defaultClassSettings := defaultGlobalDisciplines()
	if classID < uint(len(defaultClassSettings)) {

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
