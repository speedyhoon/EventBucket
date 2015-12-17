package main

import (
	"fmt"
	"net/http"
	"time"
)

func event(w http.ResponseWriter, r *http.Request, eventId string) {
	//	eventID := strings.TrimPrefix(r.URL.Path, urlEvent)
	//	if eventID == "" {
	//		http.Redirect(w, r, urlEvents, http.StatusNotFound)
	//	}
	cookie := r.Header.Get("Set-Cookie")
	if cookie != "" {

	}

	if r.URL.Path[len(urlEvent):] == "3C" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	//	eventId = "1A"
	templater(w, page{
		Title: "Event",
		Data: M{
			"EventId": eventId,
		},
	})
}

func events(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, "/events", http.StatusSeeOther)
	}
	templater(w, page{
		Title: "Events",
		Data: M{
			"Stuff": "EVENTS page!",
		},
	})
}

func insertEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		/*405 Method Not Allowed
		A request was made of a resource using a request method not supported by that resource; for example,
		using GET on a form which requires data to be presented via POST, or using POST on a read-only resource.
		//en.wikipedia.org/wiki/List_of_HTTP_status_codes*/

		//		globalSessions["fdsa"] = sessionInfo{inputs: []input{
		//			{Error: "0"},
		//			{Error: "1"},
		//			{Error: "2"},
		//			{Error: "3"},
		//		},
		//			expiry: time.Now(),
		//		}
		//		w.Header().Add("Set-Cookie", "z=fdsa; expires="+time.Now().Add(2*time.Minute).Format(gmtFormat))

		//		globalSessions["fdsa"] = sessionInfo{
		//			form: form{
		//				inputs: []input{
		//					{Error: "0"},
		//					{Required: true, Error: "1"},
		//					{Value: time.Now().Format("2006-01-02"), Error: "2"},
		//					{Value: time.Now().Format("15:04"), Error: "3"},
		//				},
		//			},
		//		}

		//		http.Redirect(w, r, "/", http.StatusSeeOther)
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)

		//redirect to home page
		//				errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	submittedFields, isValid := isValid(r, GlobalForms[0].fields)
	if !isValid {
		//		info.Println(err)

		sesId := newSessionID()
		expiry := time.Now().Add(2 * time.Minute)

		w.Header().Add("Set-Cookie", sessionError("invalid number of form items"))
		//
		//		http.Redirect(w, r, "/", http.StatusSeeOther)
		//		return

		globalSessions[sesId] = sessionInfo{
			expiry: expiry,
			form: form{
				action: "0",
				fields: submittedFields,
			},
		}
		w.Header().Add("Set-Cookie", fmt.Sprintf("z=%v; expires=%v", sesId, expiry.Format(gmtFormat)))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/event/1a", http.StatusSeeOther)

}
