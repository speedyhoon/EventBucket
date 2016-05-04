//go:generate goversioninfo -icon=favicon.ico

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
	"time"

	"github.com/boltdb/bolt"
)

var (
	//Database connection.
	db *bolt.DB

	//Command line flags.
	portAddr, dbPath string
	debug            bool

	//Used for every HTTP request with cache headers set.
	cacheExpires string

	//Logging
	t    = log.New(ioutil.Discard, "TRACE: ", log.Lshortfile) //Flags can be log.Lshortfile|log.Ltime
	info = log.New(os.Stdout, "INFO:  ", log.Lshortfile)
	warn = log.New(os.Stderr, "WARN:  ", log.Lshortfile)
)

func init() {
	go maintainExpiresTime()
	go maintainSessions()

	//Command line flags
	flag.StringVar(&dbPath, "dbpath", filepath.Join(os.Getenv("ProgramData"), `EventBucket`), "Directory for datafiles.")
	flag.BoolVar(&debug, "debug", false, "Turn on debugging and turn off HTML file caching & club maps.")
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
	portAddr = fmt.Sprintf(":%d", *port)

	setExpiresTime()
}

func main() {
	//Database save location
	dbPath = filepath.Join(dbPath, "EventBucket.db")
	info.Println("Opening database...", dbPath)
	var err error
	db, err = bolt.Open(dbPath, 0644, nil)
	if err != nil {
		warn.Println(err)
		db.Close()
		os.Exit(4)
	}
	defer db.Close()
	//Prepare database by creating all buckets (tables) needed. Otherwise view (read only) transactions will fail.
	makeBuckets()

	pages()
	info.Print("Starting EventBucket HTTP server...")
	//Open the default browser
	if !debug {
		fullAddress := "http://localhost"
		if portAddr != ":80" {
			fullAddress += portAddr
		}
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

//Update the expires http header time, every 15 minutes rather than recalculating it on every http request.
func maintainExpiresTime() {
	ticker := time.NewTicker(time.Minute * 15)
	for range ticker.C {
		//Can't directly change global variables in a go routine, so call an external function.
		setExpiresTime()
	}
}

//Set expiry date 1 year, 0 months & 0 days in the future.
func setExpiresTime() {
	//Date format is the same as Go`s time.RFC1123 but uses "GMT" timezone instead of "UTC" time standard.
	cacheExpires = time.Now().UTC().AddDate(1, 0, 0).Format(formatGMT)
	//w3.org: "All HTTP date/time stamps MUST be represented in Greenwich Mean Time" under 3.3.1 Full Date //www.w3.org/Protocols/rfc2616/rfc2616-sec3.html
	masterTemplate.CurrentYear = time.Now().Format("2006")
}
