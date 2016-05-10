package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	sessionToken      = "s"
	sessionIDLength   = 24                                                                                            //Recommended to be at least 16 characters long.
	letterBytes       = "!#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~" //string generated from validCookieValueByte golang source code net/http/cookie.go
	letterIdxBits     = 6                                                                                             // 6 bits to represent a letter index
	letterIdxMask     = 1<<letterIdxBits - 1                                                                          // All 1-bits, as many as letterIdxBits
	letterIdxMax      = 63 / letterIdxBits                                                                            // # of letter indices fitting in 63 bits
	sessionExpiryTime = 2 * time.Minute
)

var (
	src            = rand.NewSource(time.Now().UnixNano())
	globalSessions = struct {
		sync.RWMutex
		m map[string]form
	}{m: make(map[string]form)}
)

func setSession(w http.ResponseWriter, returns form) {
	cookie := http.Cookie{
		Name:     sessionToken,
		Value:    sessionID(sessionIDLength),
		HttpOnly: true,
		Expires:  time.Now().Add(sessionExpiryTime),
	}
	returns.expiry = cookie.Expires
	globalSessions.Lock()
	globalSessions.m[cookie.Value] = returns
	globalSessions.Unlock()
	http.SetCookie(w, &cookie)
}

func sessionID(n int) string {
	//author: icza, stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
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

//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainSessions() {
	ticker := time.NewTicker(sessionExpiryTime)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		purgeOldSessions()
	}
}

//Delete sessions where the expiry datetime has already lapsed.
func purgeOldSessions() {
	globalSessions.RLock()
	qty := len(globalSessions.m)
	globalSessions.RUnlock()
	if qty == 0 {
		return
	}
	globalSessions.Lock()
	t.Println("about to purge sessions, qty", len(globalSessions.m))
	for sessionID := range globalSessions.m {
		if globalSessions.m[sessionID].expiry.Before(time.Now()) {
			t.Println("deleted sessionID:", sessionID, len(globalSessions.m))
			delete(globalSessions.m, sessionID)
		}
	}
	t.Println("remaining sessions:", len(globalSessions.m))
	globalSessions.Unlock()
}

func sessionForms(w http.ResponseWriter, r *http.Request, formActions ...uint8) (uint8, []form) {
	//Add generic unpopulated form "pageError" for passing page errors from post requests to the next page served that doesn't isn't specific to a particular form.
	formActions = append(formActions, pageError)

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
				forms = append(forms, form{Fields: getForm(action)})
			}
		}
		return contents.action, forms
	}
	return 0, getForms(formActions...)
}

func getForms(formActions ...uint8) []form {
	var forms []form
	for _, action := range formActions {
		forms = append(forms, form{Fields: getForm(action)})
	}
	return forms
}
