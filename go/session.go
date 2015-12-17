package main

import (
	"fmt"
	"math/rand"
	"net/http"
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
var globalSessions = make(map[string]sessionInfo, 1)

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
//func getSession(id string) (form, error) {
func getSession(id string) form {
	contents, ok := sessionForm[id]
	if ok {
		//Clear the session contents as it has been returned to the user.
		delete(sessionForm, id)
		//		return contents, nil
		return contents
	}
	//	return form{}, errors.New("Couldn't find a session with id " + id)
	return form{}
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
	inputs []field
	form   form
	expiry time.Time
}

/* TODO
create a ticker that checks the saved sessions every 90 seconds. If the session is older than 1 minute, delete it.
*/

const cookieToken = "s"

func getSessionForm2(w http.ResponseWriter, r *http.Request) form {
	cookies, err := r.Cookie("z")
	if err != nil || cookies.Value == "" {
		return form{}
	}

	contents, ok := globalSessions[cookies.Value]
	if ok {
		//Clear the session contents as it has been returned to the user.
		delete(globalSessions, cookies.Value)
		//		return contents, nil
		w.Header().Set("Set-Cookie", fmt.Sprintf("z=; expires=%v", time.Now().UTC().Add(-5*time.Minute).Format(gmtFormat)))
		return contents.form
	}
	//	return form{}, errors.New("Couldn't find a session with id " + id)
	return form{}
}

func sessionExpire(w http.ResponseWriter, r *http.Request) {
	//	var myForm []field
	//	info.Println("globalSessions", globalSessions)
	cookies, err := r.Cookie("z")
	info.Println("r.Cookie=", cookies, err)

	//	w.Header().Del("z=fdsa")

	//	w.Header().Set("Set-Cookie", fmt.Sprintf("z=; path=/; expires=%v", time.Now().UTC().Add(-5*time.Minute).Format(gmtFormat)))
	w.Header().Set("Set-Cookie", fmt.Sprintf("z=; expires=%v", time.Now().UTC().Add(-5*time.Minute).Format(gmtFormat)))
}

func newSessionForm(w http.ResponseWriter, r *http.Request) {
	globalSessions["fdsa"] = sessionInfo{
		expiry: time.Now(),
		form: form{fields: []field{
			{Error: "-----0"},
			{Error: "-----1"},
			{Error: "-----2"},
			{Error: "-----3"},
		}},
	}
	w.Header().Add("Set-Cookie", "z=fdsa; expires="+time.Now().Add(2*time.Minute).Format(gmtFormat))
}
