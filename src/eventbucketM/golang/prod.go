// +build prod

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	//TODO Move these to ONE location & make the build tags change these?
	trace   = log.New(ioutil.Discard, "TRACE:   ", log.Lshortfile)
	info    = log.New(os.Stdout, "INFO:    ", log.Lshortfile)
	warning = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
)

func loadHTM(pageName string) string {
	switch pageName {
	case "about":
		return `^^about.htm^^`
	case "archive":
		return `^^archive.htm^^`
	case "club":
		return `^^club.htm^^`
	case "clubs":
		return `^^clubs.htm^^`
	case "event":
		return `^^event.htm^^`
	case "eventSettings":
		return `^^eventSettings.htm^^`
	case "home":
		return `^^home.htm^^`
	case "licence":
		return `^^licence.htm^^`
	case "scoreboard":
		return `^^scoreboard.htm^^`
	case "shooters":
		return `^^shooters.htm^^`
	case "start-shooting":
		return `^^start-shooting.htm^^`
	case "total-scores":
		return `^^total-scores.htm^^`
	case templateAdmin:
		return `^^_template_admin.htm^^`
	case templateEmpty:
		return `^^_template_empty.htm^^`
	case templateHome:
		return `^^_template_home.htm^^`
	}
	return ""
}

func serveDir(contentType string) {
	http.Handle(contentType,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//If url is a directory return a 404 to prevent displaying a directory listing
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}
			httpHeaders(w, []string{"expire", "cache", contentType, "public"})
			gzipper(http.FileServer(http.Dir(dirRoot)), w, r)
		}))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		gzipper(h, w, r)
	}
}

func main() {
	go startDatabase(true)
	start()
	//TODO move opening browser code back here & add code to wait for the server to connect to the DB & finish loading the http server
	info.Println("EventBucket server starting...")
	warning.Println("ListenAndServe: %v", http.ListenAndServe(":80", nil))
}
