package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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

var globalSessions = struct {
	sync.RWMutex
	m map[string]form
}{m: make(map[string]form)}

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
	globalSessions.Lock()
	globalSessions.m[sessionID] = returns
	globalSessions.Unlock()
	returns.expiry = time.Now().Add(sessionExpiryTime)
	w.Header().Add("Set-Cookie", fmt.Sprintf("%v=%v; expires=%v HttpOnly", sessionToken, sessionID, returns.expiry.Format(formatGMT)))
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

//TODO create a ticker that checks the saved sessions every 90 seconds. If the session is older than 1 minute, delete it.

//When a session id is used remove it. Supply a list of expected forms to display error messages for. Don't show errors for different pages.
func getSession(w http.ResponseWriter, r *http.Request, formActions []uint8) form {
	cookies, err := r.Cookie(sessionToken)
	if err != nil || cookies.Value == "" {
		return form{}
	}

	globalSessions.RLock()
	contents, ok := globalSessions.m[cookies.Value]
	globalSessions.RUnlock()
	if ok {
		//Clear the session contents as it has been returned to the user.
		globalSessions.Lock()
		delete(globalSessions.m, cookies.Value)
		globalSessions.Unlock()
		w.Header().Set("Set-Cookie", fmt.Sprintf("%v=; expires=%v HttpOnly", sessionToken, time.Now().UTC().Add(-sessionExpiryTime).Format(formatGMT)))
		for _, action := range formActions {
			if action == contents.action {
				return contents
			}
		}
	}
	return form{}
}

func getFormSession(w http.ResponseWriter, r *http.Request, formActions ...uint8) form {
	cookies, err := r.Cookie(sessionToken)
	if err != nil || cookies.Value == "" {
		return form{Fields: getForm(formActions[0])}
	}

	globalSessions.RLock()
	contents, ok := globalSessions.m[cookies.Value]
	globalSessions.RUnlock()
	if ok {
		//Clear the session contents as it has been returned to the user.
		globalSessions.Lock()
		delete(globalSessions.m, cookies.Value)
		globalSessions.Unlock()
		w.Header().Set("Set-Cookie", fmt.Sprintf("%v=; expires=%v HttpOnly", sessionToken, time.Now().UTC().Add(-sessionExpiryTime).Format(formatGMT)))
		for _, action := range formActions {
			if contents.action == action {
				return contents
			}
		}
	}
	return form{Fields: getForm(formActions[0])}
}

/*


//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainExpiresTime() {
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
	//Can't directly change global variables in a go routine, so call an external function.
	setExpiresTime()
}
}

//Set expiry date 1 year, 0 months & 0 days in the future.
func setExpiresTime() {
	//Date format is the same as Go`s time.RFC1123 but uses "GMT" timezone instead of "UTC" time standard.
	expiresTime = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
	//w3.org says: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
	masterTemplate.CurrentYear = time.Now().Format("2006")
}
*/
