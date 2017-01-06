//go:generate goversioninfo -icon=..\icon\app.ico

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var (
	//Database connection.
	db *bolt.DB

	//Command line flags.
	portAddr, dbPath string
	debug            bool

	//Logging
	t    = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	info = log.New(os.Stdout, "", 0)
	warn = log.New(os.Stderr, "WARN: ", log.Lshortfile)
)

func init() {
	//go maintainExpiresTime()

	//Command line flags
	flag.StringVar(&dbPath, "dbpath", filepath.Join(os.Getenv("ProgramData"), `EventBucket`), "Directory for datafiles.")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging and turn off HTML file caching & club maps.")
	flag.BoolVar(&masterTemplate.Theme, "dark", false, "Switch EventBucket to use a dark theme for night shooting")
	gradesFilePath := flag.String("grades", "", "Load grade settings from a JSON file. If the file doesn't exist, EventBucket will try to create it & exit")
	port := flag.Uint("port", 80, "Assign a differnet port number for the HTTP server. Range: 1 through 65535. Some port numbers may already be in use on this system.")
	flag.Parse()

	//Create database directory if needed.
	err := mkDir(dbPath)
	if err != nil {
		warn.Println(err)
		os.Exit(1)
	}

	if debug {
		t.SetOutput(os.Stdout)
	}

	//Try to load the grades file if any is specified
	if loadGrades(*gradesFilePath) != nil {
		redoGlobals([]Discipline{})
		//If a file path was specified try to create one
		if *gradesFilePath != "" {
			buildGradesFile(*gradesFilePath)
			os.Exit(2)
		}
	}

	//Check port number
	if *port > math.MaxUint16 || *port < 1 {
		warn.Println("Port number must be between 1 and 65535. (default 80)")
		os.Exit(3)
	}
	if *port != 80 {
		portAddr = fmt.Sprintf(":%d", *port)
	}

	//setExpiresTime()
}

func main() {
	go maintainSessions()

	//Database save location
	dbPath = filepath.Join(dbPath, "EventBucket.db")
	info.Println("Opening database...", dbPath)
	var err error
	db, err = bolt.Open(dbPath, 0644, nil)
	if err != nil {
		warn.Println(err)
		os.Exit(4)
	}
	defer db.Close()
	//Prepare database by creating all buckets (tables) needed. Otherwise view (read only) transactions will fail.
	makeBuckets()

	pages()
	info.Print("Starting EventBucket HTTP server...")
	//Open the default browser
	if !debug {
		fullAddress := "http://localhost"+portAddr
		info.Print(fullAddress)
		if exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", fullAddress).Start() != nil {
			warn.Println("Unable to open a web browser for", fullAddress)
		}
	}
	warn.Printf("ListenAndServe: %v", http.ListenAndServe(portAddr, nil))
	info.Println("EvenBucket server stopped.")
}

//Attempt to create the path supplied if it doesn't exist.
func mkDir(path string) error {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		err = os.Mkdir(path, os.ModeDir)
		if err != nil {
			warn.Printf("Unable to create directory %v %v", path, err)
		}
	}
	return err
}
