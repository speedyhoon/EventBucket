package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/home.css", http.FileServer(http.Dir("./root/c/")))
	//	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//		// The "/" pattern matches everything, so we need to check
	//		// that we're at the root here.
	//		if req.URL.Path != "/" {
	//			http.NotFound(w, req)
	//			return
	//		}
	//		fmt.Fprintf(w, "Welcome to the home page!")
	//	})
	srv := &http.Server {
		Addr: ":80",
		Handler: mux,
	}
	srv.ListenAndServe()
}
