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
/*func newSessionID() string {
	var sessionID string
	var i, randInt int
	for i < 24 {
		randInt = 33 + rand.Intn(93)
		if randInt != 59 { //ignore semicolons ;
			i++
			sessionID += string(randInt)
		}
	}
	return sessionToken + "=" + sessionID
}*/

//TODO create a ticker that checks the saved sessions every 90 seconds. If the session is older than 1 minute, delete it.

//When a session id is used remove it. Supply a list of expected forms to display error messages for. Don't show errors for different pages.
/*func getSession(w http.ResponseWriter, r *http.Request, formActions []uint8) form {
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
}*/

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

func sessionForms(w http.ResponseWriter, r *http.Request, formActions ...uint8) (*uint8, [][]field) {
	//Get users session id from request
	cookies, err := r.Cookie(sessionToken)
	if err != nil || cookies.Value == "" {
		//No session found. Return default forms.
		return nil, getForms(formActions...)
	}

	//Start a read lock to prevent concurrent reads while other parts are executing a write.
	globalSessions.RLock()
	contents, ok := globalSessions.m[cookies.Value]
	globalSessions.RUnlock()
	if ok {
		//Clear the session contents as it has been returned to the user.
		globalSessions.Lock()
		delete(globalSessions.m, cookies.Value)
		globalSessions.Unlock()

		//Remove cookie.
		w.Header().Set("Set-Cookie", fmt.Sprintf("%v=; expires=%v HttpOnly", sessionToken, time.Now().UTC().Add(-sessionExpiryTime).Format(formatGMT)))
		var forms [][]field
		for _, action := range formActions {
			if contents.action == action {
				//				forms = append(forms, contents)
				forms = append(forms, contents.Fields)
			} else {
				//				forms = append(forms, form{action: action, Fields: getForm(action)})
				forms = append(forms, getForm(action))
			}
		}
		//TODO is it possible that reference to contents could possibly blow up?
		return &contents.action, forms
	}
	return nil, getForms(formActions...)
}

func getForms(formActions ...uint8) [][]field {
	var forms [][]field //A group of forms, where each form has several fields
	for _, action := range formActions {
		forms = append(forms, getForm(action))
	}
	return forms
}

//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainSessions() {
	ticker := time.NewTicker(time.Second * 60)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		purgeOldSessions()
	}
}

//Delete sessions where the expiry datetime has already lapsed.
func purgeOldSessions() {
	globalSessions.Lock()
	for sessionID := range globalSessions.m {
		if globalSessions.m[sessionID].expiry.Before(time.Now()) {
			trace.Println("deleted sessionID:", sessionID)
			delete(globalSessions.m, sessionID)
		}
	}
	globalSessions.Unlock()
}

func sessionForms2(w http.ResponseWriter, r *http.Request, formActions ...uint8) (*uint8, []form) {
	//Get users session id from request
	cookies, err := r.Cookie(sessionToken)
	if err != nil || cookies.Value == "" {
		//No session found. Return default forms.
		//		return form{Fields: getForm(formActions[0])}
		return nil, getForms2(formActions...)
	}

	//Start a read lock to prevent concurrent reads while other parts are executing a write.
	globalSessions.RLock()
	contents, ok := globalSessions.m[cookies.Value]
	globalSessions.RUnlock()
	if ok {
		//Clear the session contents as it has been returned to the user.
		globalSessions.Lock()
		delete(globalSessions.m, cookies.Value)
		globalSessions.Unlock()

		//Remove cookie.
		w.Header().Set("Set-Cookie", fmt.Sprintf("%v=; expires=%v HttpOnly", sessionToken, time.Now().UTC().Add(-sessionExpiryTime).Format(formatGMT)))
		var forms []form
		for _, action := range formActions {
			if contents.action == action {
				//				forms = append(forms, contents)
				forms = append(forms, form{Fields: contents.Fields})
			} else {
				//				forms = append(forms, form{action: action, Fields: getForm(action)})
				forms = append(forms, form{Fields: getForm(action)})
			}
		}
		//TODO is it possible that reference to contents could possibly blow up?
		return &contents.action, forms
	}
	return nil, getForms2(formActions...)
}

func getForms2(formActions ...uint8) []form {
	var forms []form
	for _, action := range formActions {
		//		forms = append(forms, form{action: action, Fields: getForm(action)})
		forms = append(forms, form{Fields: getForm(action)})
	}
	return forms
}
