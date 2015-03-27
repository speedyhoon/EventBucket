// +build prod

package main

import (
	"fmt"
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

func devModeCheckForm(check bool, message string) {
	if !check {
		Warning.Println(message)
	}
}

func loadHTM(pageName string) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(PATH_HTML_MINIFIED, pageName))
	if err != nil {
		//TODO inline html as bytes here
		Error.Println(err)
	}
	return bytes
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
			Gzip(http.FileServer(http.Dir("^^DIR_ROOT^^")), w, r)
		}))
}

func serveHtml(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHeaders(w, []string{"html", "noCache", "expireNow", "pragma"})
		Gzip(h, w, r)
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
