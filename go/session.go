package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	sessionToken      = "s"
	sessionExpiryTime = time.Minute * 2
)

var (
	src            = rand.NewSource(time.Now().UnixNano())
	globalSessions = struct {
		sync.RWMutex
		m map[string]form
	}{m: make(map[string]form)}
)

func setSession(w http.ResponseWriter, f form) {
	f.expiry = time.Now().Add(sessionExpiryTime)
	cookie := http.Cookie{
		Name:     sessionToken,
		Value:    sessionID(),
		HttpOnly: true,
		Expires:  f.expiry,
	}

	//Start mutex write lock.
	globalSessions.Lock()
	globalSessions.m[cookie.Value] = f
	globalSessions.Unlock()
	http.SetCookie(w, &cookie)
}

func sessionID() string {
	const (
		//string generated from validCookieValueByte golang source code net/http/cookie.go
		letterBytes   = "!#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
		n             = 24                   //Session ID length is recommended to be at least 16 characters long.
		letterIdxBits = 6                    //6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 //All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   //# of letter indices fitting in 63 bits
	)
	//author: icza, stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	b := make([]byte, n)
	//A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

//maintainSessions periodically deletes expired sessions
func maintainSessions() {
	ticker := time.NewTicker(sessionExpiryTime)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		purgeSessions()
	}
}

//Delete sessions where the expiry datetime has already lapsed.
func purgeSessions() {
	globalSessions.RLock()
	qty := len(globalSessions.m)
	globalSessions.RUnlock()
	if qty == 0 {
		return
	}

	t.Println("About to purge sessions, qty", qty)
	globalSessions.Lock()
	for sessionID := range globalSessions.m {
		if globalSessions.m[sessionID].expiry.Before(time.Now()) {
			delete(globalSessions.m, sessionID)
		}
	}
	t.Println("Remaining sessions:", len(globalSessions.m))
	globalSessions.Unlock()
}

func sessionForms(w http.ResponseWriter, r *http.Request, formActions ...uint8) (uint8, []form) {
	//Add generic unpopulated form for passing page errors from post requests to the next page served that isn't specific to a particular form.
	formActions = append(formActions, 255)

	//Get users session id from request cookie header
	cookie, err := r.Cookie(sessionToken)
	if err != nil || cookie == nil || cookie.Value == "" {
		//No session found. Return default forms.
		return 0, getForms(formActions...)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionToken,
		Value:    "",                                       //Remove cookie by setting it to nothing (empty string).
		HttpOnly: true,                                     //HttpOnly means the cookie can't be accessed by JavaScript
		Expires:  time.Now().UTC().Add(-sessionExpiryTime), //Using minus sessionExpiryTime so the session expires time is set to the past
	})

	//Start a read lock to prevent concurrent reads while other parts are executing a write.
	globalSessions.RLock()
	contents, ok := globalSessions.m[cookie.Value]
	globalSessions.RUnlock()
	if ok {
		//Clear the session contents as it has been returned to the user.
		globalSessions.Lock()
		delete(globalSessions.m, cookie.Value)
		globalSessions.Unlock()

		var forms []form
		for _, action := range formActions {
			if contents.action == action {
				forms = append(forms, contents)
			} else {
				forms = append(forms, getForm(action))
			}
		}
		return contents.action, forms
	}
	return 0, getForms(formActions...)
}

func getForms(formActions ...uint8) (forms []form) {
	for _, action := range formActions {
		forms = append(forms, getForm(action))
	}
	return forms
}
