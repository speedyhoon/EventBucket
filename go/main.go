//go:generate goversioninfo -icon=..\icon\app.ico

package main

import (
	"context"
	"flag"
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
	isPrivate bool

	//Logging
	//#ifdef DEBUG
	trc = log.New(os.Stdout, "\x1b[33;1mTRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	inf = log.New(os.Stdout, "\x1b[36;1m", 0)
	//#endif
)

func main() {
	//Command line flags
	flag.BoolVar(&mainTheme, "dark", false, "Switch EventBucket to use a dark theme for night shooting")
	gradesFilePath := flag.String("grades", "", "Load grade settings from a JSON file. If the file doesn't exist, EventBucket will try to create it & exit")
	httpListen := flag.String("http", "127.0.0.1:80", "host:port to listen on")
	dbPath := flag.String("dbpath", filepath.Join(os.Getenv("ProgramData"), "EventBucket", "EventBuc.ket"), "Directory for datafiles.")
	flag.Parse()

	//#ifdef DEBUG
	log.SetPrefix("\x1b[31;1mWARN: ") //Red
	//#endif
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

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
			log.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	//HTTP server host & port
	host, port, err := net.SplitHostPort(*httpListen)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}

		inf.Print("Started EventBucket HTTP server...")
		//#ifdef DEBUG
		inf.Println(httpAddr)
		//#else
		httpAddr = "http://" + httpAddr
		openBrowser(httpAddr)
		inf.Printf("A browser window should open. If not, please visit %s", httpAddr)
		//#endif
	}()

	<-stop
	inf.Println("Shutting down the server...")
	err = h.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
	inf.Println("EvenBucket server stopped.")
}
