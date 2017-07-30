//go:generate goversioninfo -icon=..\icon\app.ico

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"context"
	"os/signal"
	"time"

	"github.com/boltdb/bolt"
)

var (
	//Command line flags.
	portAddr   = ":"
	dbPath     = flag.String("dbpath", filepath.Join(os.Getenv("ProgramData"), "EventBucket"), "Directory for datafiles.")
	debug      = flag.Bool("debug", false, "Turn on debugging and turn off club maps.")
	httpListen = flag.String("http", "127.0.0.1:80", "host:port to listen on")

	//Logging
	//TODO add t & info when debug == true during build time
	t    = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	info = log.New(os.Stdout, "", 0)
	warn = log.New(os.Stderr, "WARN: ", log.Lshortfile)
)

func init() {
	//Command line flags
	flag.BoolVar(&masterTemplate.Theme, "dark", false, "Switch EventBucket to use a dark theme for night shooting")
	gradesFilePath := flag.String("grades", "", "Load grade settings from a JSON file. If the file doesn't exist, EventBucket will try to create it & exit")
	flag.Parse()

	//Create database directory if needed.
	err := mkDir(*dbPath)
	if err != nil {
		warn.Fatal(err)
	}

	if *debug {
		t.SetOutput(os.Stdout)
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
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	//HTTP server host & port
	host, port, err := net.SplitHostPort(*httpListen)
	if err != nil {
		warn.Fatal(err)
	}

	//Database save location
	path := filepath.Join(*dbPath, "EventBucket.db")

	db, err = bolt.Open(path, 0644, &bolt.Options{
		Timeout:         time.Second * 8,
		InitialMmapSize: 1048576, //Initial database size = 1MB
	})
	if err != nil {
		warn.Fatal("Connection timeout. Unable to open", path)
	}
	defer db.Close()

	//Prepare database by creating all buckets (tables) needed. Otherwise view (read only) transactions will fail.
	makeBuckets()

	if host == "" {
		host = "localhost"
	}
	if port != "80" {
		portAddr = ":" + port
	}
	httpAddr := host + portAddr
	h := http.Server{Addr: httpAddr, Handler: nil}
	go func() {
		if err = h.ListenAndServe(); err != nil {
			warn.Fatal(err)
		} else {
			info.Print("Started EventBucket HTTP server...")
			url := "http://" + httpAddr
			if !*debug && openBrowser(url) {
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
