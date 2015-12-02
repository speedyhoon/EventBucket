package main

import (
	"errors"
	"math/rand"
)

const (
	sessionIDLength  = 20 //Recommended to be at least 16 characters long.
	semicolon        = 59
	sessionCharStart = 33
	sessionCharEnd   = 126
	sessionCharRange = sessionCharEnd - sessionCharStart
)

var sessionForm map[string]form

//TODO will possibly need a chanel here to prevent locks occurring
func setSession(returns form) {
	var number int
	var sessionID string
	var ok bool
	//generate a new session id while it doesn't already exist
	for _, ok = sessionForm[sessionID]; sessionID != "" && !ok; {
		//clear any previous id that was generated
		sessionID = ""
		for len(sessionID) < sessionIDLength {
			number = rand.Intn(sessionCharRange) + sessionCharStart
			switch number {
			//ignore semicolons ";" as these characters terminate the end of the session id
			case semicolon:
				sessionID += " "
			default:
				//converts number to a alphanumeric or symbol character
				sessionID += string(number)
			}
		}
	}
	sessionForm[sessionID] = returns
}

//When a previous session id is used remove it.
func getSession(id string) (form, error) {
	contents, ok := sessionForm[id]
	if ok {
		//Clear the session contents as it has been returned to the user.
		delete(sessionForm, id)
		return contents, nil
	}
	return form{}, errors.New("Couldn't find a session with id " + id)
}
