// +build prod

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	Info    = log.New(ioutil.Discard, "INFO:    ", log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Lshortfile)
)

func loadHTM(pageName string) []byte {
	switch pageName {
	case "about.htm":
		return `^^about.htm^^`
	case "archive.htm":
		return `^^archive.htm^^`
	case "club.htm":
		return `^^club.htm^^`
	case "clubs.htm":
		return `^^clubs.htm^^`
	case "event.htm":
		return `^^event.htm^^`
	case "eventSettings.htm":
		return `^^eventSettings.htm^^`
	case "home.htm":
		return `^^home.htm^^`
	case "licence.htm":
		return `^^licence.htm^^`
	case "scoreboard.htm":
		return `^^scoreboard.htm^^`
	case "shooters.htm":
		return `^^shooters.htm^^`
	case "start-shooting.htm":
		return `^^start-shooting.htm^^`
	case "total-scores.htm":
		return `^^total-scores.htm^^`
	case "_template_admin.htm":
		return `^^_template_admin.htm^^`
	case "_template_empty.htm":
		return `^^_template_empty.htm^^`
	case "_template_home.htm":
		return `^^_template_home.htm^^`
	}
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
			gzipper(http.FileServer(http.Dir("^^dirRoot^^")), w, r)
		}))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		gzipper(h, w, r)
	}
}

func main() {
	start()
	url := "http://localhost"
	if exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start() != nil {
		Warning.Printf("Unable to open a web browser for " + url)
	}
	Warning.Println("ListenAndServe: %v", http.ListenAndServe(":80", nil))
}
