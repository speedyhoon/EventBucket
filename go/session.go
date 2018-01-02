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

var globalSessions = struct {
	sync.RWMutex
	m map[string]form
}{m: make(map[string]form)}

func init() {
	go sessionUpkeep()
}

//sessionUpkeep periodically deletes expired sessions
func sessionUpkeep() {
	for range time.NewTicker(sessionExpiryTime).C {
		//Can't directly change global variables in a go routine, so call an external function.
		sessionPurge()
	}
}

//Delete sessions where the expiry datetime has already lapsed.
func sessionPurge() {
	globalSessions.RLock()
	qty := len(globalSessions.m)
	globalSessions.RUnlock()
	if qty == 0 {
		return
	}

	t.Println("About to purge sessions. Qty:", qty)
	now := time.Now()
	globalSessions.Lock()
	for sessionID := range globalSessions.m {
		if globalSessions.m[sessionID].expiry.Before(now) {
			delete(globalSessions.m, sessionID)
		}
	}
	t.Println("Remaining sessions:", len(globalSessions.m))
	globalSessions.Unlock()
}

func sessionSet(w http.ResponseWriter, f form) {
	f.expiry = time.Now().Add(sessionExpiryTime)
	var ok bool
	var id string

	//Start mutex write lock.
	globalSessions.Lock()
	for {
		id = sessionID()
		_, ok = globalSessions.m[id]

		if !ok {
			//Assign the session ID if it isn't already assigned
			globalSessions.m[id] = f
			break
		}
		//else if sessionID is already assigned then regenerate a different session ID
	}
	globalSessions.Unlock()

	cookie := http.Cookie{
		Name:     sessionToken,
		Value:    id,
		HttpOnly: true,
		Expires:  f.expiry,
	}
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
	src := rand.NewSource(time.Now().UnixNano())
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

//sessionForms retrieves a slice of forms, including any errors (if any)
func sessionForms(w http.ResponseWriter, r *http.Request, formIDs ...uint8) (uint8, []form) {
	const noAction = 255

	//Get users session id from request cookie header
	cookie, err := r.Cookie(sessionToken)
	if err != nil || cookie == nil || cookie.Value == "" {
		//No session found. Return default forms.
		return noAction, getForms(formIDs...)
	}

	//Remove client session cookie
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
	if !ok {
		return noAction, getForms(formIDs...)
	}

	//Clear the session contents as it has been returned to the user.
	globalSessions.Lock()
	delete(globalSessions.m, cookie.Value)
	globalSessions.Unlock()

	var forms []form
	for _, id := range formIDs {
		if contents.action == id {
			forms = append(forms, contents)
		} else {
			forms = append(forms, getForm(id))
		}
	}
	return contents.action, forms
}

func getForms(actions ...uint8) (forms []form) {
	for _, id := range actions {
		forms = append(forms, getForm(id))
	}
	return forms
}
