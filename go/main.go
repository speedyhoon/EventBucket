//go:generate goversioninfo -icon=..\icon\app.ico

package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

var (
	//Command line flags.
	portAddr  = ":"
	debug     bool
	isPrivate bool

	//Logging
	//TODO add t & info when debug == true during build time
	t    = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	info = log.New(os.Stdout, "", 0)
	warn = log.New(os.Stderr, "WARN: ", log.Lshortfile)
)

func main() {
	//Command line flags
	isDark := flag.Bool("dark", false, "Switch EventBucket to use a dark theme for night shooting")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging and turn off club maps.")
	gradesFilePath := flag.String("grades", "", "Load grade settings from a JSON file. If the file doesn't exist, EventBucket will try to create it & exit")
	httpListen := flag.String("http", "127.0.0.1:80", "host:port to listen on")
	dbPath := flag.String("dbpath", filepath.Join(os.Getenv("ProgramData"), "EventBucket", "EventBucket.db"), "Directory for datafiles.")
	flag.Parse()

	if *isDark {
		mainTheme = 1
	}

	//Turn on trace logging
	if debug {
		t.SetOutput(os.Stdout)
		t.SetPrefix("\x1b[33;1m" + t.Prefix())       //Yellow
		info.SetPrefix("\x1b[36;1m" + info.Prefix()) //Blue
		warn.SetPrefix("\x1b[31;1m" + warn.Prefix()) //Red
	}

	//Try to load the grades file if any are specified
	if loadGrades(*gradesFilePath) != nil {
		redoGlobals([]Discipline{})
		//If a file path was specified try to create one
		if *gradesFilePath != "" {
			buildGradesFile(*gradesFilePath)
			os.Exit(2)
		}
	}

	startDB(*dbPath)
	defer func() {
		if err := db.Close(); err != nil {
			warn.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	//HTTP server host & port
	host, port, err := net.SplitHostPort(*httpListen)
	if err != nil {
		warn.Fatal(err)
	}

	const localhost = "localhost"
	if host == "" {
		host = localhost
	}
	if port != "80" {
		portAddr += port
	}

	num, _ := stoU(port)
	isPrivate = host != "127.0.0.1" && host != localhost || num > 1023

	httpAddr := host + portAddr
	h := http.Server{Addr: httpAddr, Handler: nil}
	go func() {
		if err = h.ListenAndServe(); err != nil {
			warn.Fatal(err)
		} else {
			info.Print("Started EventBucket HTTP server...")
			url := "http://" + httpAddr
			if !debug && openBrowser(url) {
				info.Printf("A browser window should open. If not, please visit %s", url)
			} else {
				info.Printf("Please open your web browser and visit %s", url)
			}
		}
	}()

	<-stop
	info.Println("Shutting down the server...")
	err = h.Shutdown(context.Background())
	if err != nil {
		warn.Println(err)
	}
	info.Println("EvenBucket server stopped.")
}
