package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	sessionToken      = "s"
	sessionIDLength   = 24 //Recommended to be at least 16 characters long.
	sessionCharStart  = 33
	sessionCharRange  = 126 - sessionCharStart
	semicolon         = 59
	sessionExpiryTime = 2 * time.Minute
)

type sessionInfo struct {
	inputs []field
	form   form
	expiry time.Time
}

var globalSessions = make(map[string]sessionInfo, 1)

//TODO will possibly need a chanel here to prevent locks occurring
func setSession(w http.ResponseWriter, returns form) {
	var number int
	var sessionID string
	//var ok bool
	//generate a new session id while it doesn't already exist
	//for _, ok = globalSessions[sessionID]; sessionID != "" && !ok; {
	//clear any previous id that was generated
	sessionID = ""
	//TODO could possibly use cryptography random letter generator
	for len(sessionID) < sessionIDLength {
		number = rand.Intn(sessionCharRange) + sessionCharStart
		//ignore semicolons ";" as these characters terminate the end of the session id
		if number != semicolon {
			//converts number to a alphanumeric or symbol character
			sessionID += string(number)
		}
	}
	//	}
	expiry := time.Now().Add(sessionExpiryTime)
	globalSessions[sessionID] = sessionInfo{
		expiry: expiry,
		form:   returns,
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("%v=%v; expires=%v", sessionToken, sessionID, expiry.Format(formatGMT)))
}

//SessionID's should be at least 16 characters length can't have space or semicolon
//stackoverflow.com/questions/1969232/allowed-characters-in-cookies
func newSessionID() string {
	var sessionId string
	var i, randInt int
	for i < 24 {
		randInt = 33 + rand.Intn(93)
		if randInt != 59 { //ignore semicolons ;
			i++
			sessionId += string(randInt)
		}
	}
	return sessionToken + "=" + sessionId
}

/* TODO
create a ticker that checks the saved sessions every 90 seconds. If the session is older than 1 minute, delete it.
*/

//When a session id is used remove it.
func getSession(w http.ResponseWriter, r *http.Request) form {
	cookies, err := r.Cookie(sessionToken)
	if err != nil || cookies.Value == "" {
		return form{}
	}

	contents, ok := globalSessions[cookies.Value]
	if ok {
		//Clear the session contents as it has been returned to the user.
		delete(globalSessions, cookies.Value)
		w.Header().Set("Set-Cookie", fmt.Sprintf("%v=; expires=%v", sessionToken, time.Now().UTC().Add(-sessionExpiryTime).Format(formatGMT)))
		return contents.form
	}
	return form{}
}
