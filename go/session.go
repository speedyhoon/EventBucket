package main

import (
	"errors"
	"math/rand"
	"time"
)

const (
	sessionIDLength  = 24 //Recommended to be at least 16 characters long.
	semicolon        = 59
	sessionCharStart = 33
	sessionCharRange = 126 - sessionCharStart
)

//TODO leave only one of these global variables
var sessionForm = make(map[string]form, 1)
var globalSessions map[string]sessionInfo

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
			//ignore semicolons ";" as these characters terminate the end of the session id
			if number != semicolon {
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

//SessionID's can't have space or semicolon
//should be 16 characters
//stackoverflow.com/questions/1969232/allowed-characters-in-cookies
func newSessionID() string {
	var newSessionId string
	var i, randInt int
	for i < 24 {
		randInt = 33 + rand.Intn(93)
		if randInt != 59 { //ignore semicolons ;
			i++
			newSessionId += string(randInt)
		}
	}
	return "z=" + newSessionId
}

func sessionError(error string) string {
	sessionID := newSessionID()
	//	globalSessions[sessionID] = error
	return sessionID + "; Expires=" + time.Now().Add(1*time.Minute).Format("Mon, Jan 2 2006 15:04:05 GMT")
}

type sessionInfo struct {
	inputs []input
	expiry time.Time
}

/* TODO
create a ticker that checks the saved sessions every 90 seconds. If the session is older than 1 minute, delete it.
*/
