package main

import (
"net/http"

"github.com/yosssi/ace"
"log"
"os"
"os/exec"

"github.com/tdewolff/minify"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := ace.Load("example", "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "Actions",
		"Msgs": []string{
			"Message1",
			"Message2",
			"Message3",
		},
	}
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}




	m := minify.NewMinifierDefault()
	m.AddCmd("text/javascript", exec.Command("java", "-jar", "build/compiler.jar"))

	if err := m.Minify("text/html", os.Stdout, os.Stdin); err != nil {
		log.Fatal("Minify:", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
